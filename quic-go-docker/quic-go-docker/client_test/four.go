package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Client --
func Client() {
	signalAddress := os.Args[2]
	localAddress := ":8888" // default port
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}

	remoteRelay, _ := net.ResolveUDPAddr("udp", signalAddress)
	local, _ := net.ResolveUDPAddr("udp", localAddress)
	udpConn, err := net.ListenUDP("udp", local) //创建UDP连接
	if err != nil {
		fmt.Printf("UDP connection error: %s\n", err)
		panic(err)
	}
	//在单个 UDP 套接字上运行的 QUIC 连接
	tr := quic.Transport{
		Conn: udpConn,
	}

	//监听的tls
	cert, err := generateSelfSignedCertificate()
	if err != nil {
		fmt.Printf("Error generating certificate: %s\n", err)
		return
	}
	// QUIC 监听配置
	tlsConfListen := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-holepunch"},
	}
	// QUIC 拨号配置
	tlsConfDial := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-holepunch"},
	}
	//-------------------------------------上述代码是公用部分（client-relay || client-client）-------------------------------------
	peerAddrChan := make(chan string) // 创建一个用于传递 peerAddr 的通道
	go func() {
		// 函数执行完成后减少一个等待的 goroutine
		peerAddr := register(&tr, remoteRelay, tlsConfDial) // 在线程1中执行 register() 函数
		peerAddrChan <- peerAddr                            // 将 peerAddr 发送到通道中
	}()
	peerAddr := <-peerAddrChan // 从通道中接收 peerAddr 值
	//先使得双方都收到对端公网地址再打洞
	time.Sleep(time.Second)
	go func() {
		// 记录holepunch开始时间
		executionTime := holepunch(&tr, tlsConfListen, tlsConfDial, peerAddr) // 在线程2中使用 peerAddr 值执行 chatter 函数
		fmt.Println("holepunch took", executionTime, time.Now())
	}()

	time.Sleep(time.Second * 150) //等待100s主线程退出
	log.Println("Main goroutine exit")
}

// 客户端和中继连接，并返回公网地址和端口号
func register(tr *quic.Transport, remote *net.UDPAddr, tlsConf *tls.Config) string {
	// 使用 QUIC 和中继连接
	conn, err := tr.Dial(context.Background(), remote, tlsConf, nil)
	if err != nil {
		fmt.Printf("客户端连接出现错误: %s\n", err)
		panic(err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		_, err := stream.Write([]byte("register"))
		if err != nil {
			panic(err)
			fmt.Printf("发送数据出现错误: %s\n", err)
		}
	}()

	var peerAddr string // 声明 peerAddr 变量
	for {
		fmt.Println("listening...")
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		buffer := make([]byte, 1024)
		bytesRead, err := stream.Read(buffer)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}

		data := string(buffer[0:bytesRead])
		addresses := strings.Split(data, ",")
		localAddr := addresses[0]
		peerAddr = addresses[1]
		fmt.Println("自己的公网IP:port--->", localAddr, time.Now())
		fmt.Println("对端的公网IP:port--->", peerAddr, time.Now())
		break
	}
	//接受从服务器发送过来的对端地址和端口号
	return peerAddr
}

// 客户端之间的连接（打洞）
func holepunch(tr *quic.Transport, tlsConfListen *tls.Config, tlsConfDial *tls.Config, remote string) time.Duration {
	remoteAddr, _ := net.ResolveUDPAddr("udp", remote)

	////创建监听对象
	listener, err := tr.Listen(tlsConfListen, &quic.Config{KeepAlivePeriod: 1})
	if err != nil {
		fmt.Printf("Create listener object error: %s\n", err)
		panic(err)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 2s handshake timeout
	//defer cancel()

	var endTime time.Time
	var startTime time.Time
	done := make(chan struct{}) // 创建一个用于通知函数结束的channel

	go func() {
		defer close(done)
		var connToPeer quic.Connection //dial的连接
		var connListen quic.Connection //accept的连接
		dialDone := make(chan struct{})
		acceptDone := make(chan struct{})

		// 开始 Dial
		go func() {
			defer close(dialDone)
			//计算打洞开始时间
			startTime = time.Now() //打洞时间从这里开始的目的是不能够忽略掉dial失败的时间
			for i := 0; i < 20; i++ {
				log.Printf("Attempting to dial %s\n", remote)
				var dialErr error
				connToPeer, dialErr = tr.Dial(context.Background(), remoteAddr, tlsConfDial, &quic.Config{KeepAlivePeriod: 1})
				if dialErr == nil {
					endTime = time.Now()
					//计算打洞结束时间
					break
				} else {
					log.Println("Dial failed, reattempting", dialErr)
					time.Sleep(time.Second) //间隔1s之后再打洞
				}
			}
		}()

		// 开始 Accept
		go func() {
			defer close(acceptDone)
			for {
				var acceptErr error
				log.Printf("Listening for connection from %s\n", remote, time.Now())
				connListen, acceptErr = listener.Accept(context.Background())
				if acceptErr != nil {
					fmt.Printf("Failed to accept when listen: %s\n", acceptErr, time.Now())
				} else {
					log.Println("hhhhhhhhhhhhhhhh", time.Now())
					break
				}
			}

		}()

		if connToPeer == nil {
			log.Println("Unable to establish connection")
			return
		}
		log.Println("Successfully established connection with peer")
		sendStream, errSendStream := connToPeer.OpenUniStreamSync(context.Background())
		if errSendStream != nil {
			log.Printf("Error with OpenUniStreamSync: %s\n", errSendStream)
			return
		}
		_, err = sendStream.Write([]byte("Hello!"))
		if err != nil {
			fmt.Printf("Failed to send message: %s\n", err)
		} else {
			log.Println("Successfully sent [Hello!] to ---> ", remote, time.Now())
		}

		//接收流
		receiveStream, errReceiveStream := connListen.AcceptUniStream(context.Background())
		if errReceiveStream != nil {
			fmt.Println("Failed to accept stream", errReceiveStream)
		}
		buffer := make([]byte, 1024)
		bytesRead, errBytesRead := receiveStream.Read(buffer)
		if errBytesRead != nil {
			fmt.Printf("Failed to read the data from peer: %s\n", errBytesRead)
			return
		}
		log.Println("[INCOMING-FROM-peer]", string(buffer[0:bytesRead]), time.Now())
		if string(buffer[0:bytesRead]) == "Hello!" {
			log.Println("Successfully received [Hello!] message from peer ! ! !")
			log.Println("NOW ! Now you can write your message at [Enter message] ")
			log.Println()
		}
	}()
	<-done // 等待通知，确保goroutine执行完毕
	log.Println("startTime", startTime)
	log.Println("endTime", endTime)
	// 计算时间差
	executionTime := endTime.Sub(startTime)
	return executionTime
}
