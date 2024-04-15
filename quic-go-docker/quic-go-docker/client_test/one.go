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
	// 增加 UDP 缓冲区大小
	//err = udpConn.SetReadBuffer(2048 * 1024) // 设置为期望的大小
	//if err != nil {
	//	fmt.Printf("Error setting UDP buffer size: %s\n", err)
	//	panic(err)
	//}

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

	go func() {
		// 记录holepunch开始时间
		executionTime := holepunch(&tr, tlsConfListen, tlsConfDial, peerAddr) // 在线程2中使用 peerAddr 值执行 chatter 函数
		fmt.Println("holepunch took", executionTime)
	}()

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
		fmt.Println("自己的公网IP:port--->", localAddr)
		fmt.Println("对端的公网IP:port--->", peerAddr)
		//stream.Close()
		break
	}
	//接受从服务器发送过来的对端地址和端口号
	return peerAddr
}

// 客户端之间的连接（打洞）
func holepunch(tr *quic.Transport, tlsConfListen *tls.Config, tlsConfDial *tls.Config, remote string) time.Duration {
	remoteAddr, _ := net.ResolveUDPAddr("udp", remote)

	////创建监听对象
	listener, err := tr.Listen(tlsConfListen, nil)
	if err != nil {
		fmt.Printf("Create listener object error: %s\n", err)
		panic(err)
	}

	////定义一个接受流的管道
	//receiveStreamChan := make(chan quic.ReceiveStream)
	////定义一个发送流的管道
	//sendStreamChan := make(chan quic.SendStream)
	//var executionTime time.Duration
	var endTime time.Time
	var startTime time.Time
	done := make(chan struct{}) // 创建一个用于通知函数结束的channel

	go func() {
		var conntoPeer quic.Connection
		go func() {
			defer close(done)
			log.Printf("Listening for connection from %s\n", remote)
			connListen, err := listener.Accept(context.Background())
			if err != nil {
				fmt.Printf("Failed to accept when listen: %s\n", err)
				return
			}

			receiveStream, err := connListen.AcceptUniStream(context.Background())
			if err != nil {
				fmt.Println("Failed to accept stream", err)
			}
			buffer := make([]byte, 1024)
			bytesRead, err := receiveStream.Read(buffer)
			if err != nil {
				fmt.Printf("Failed to read the data from peer: %s\n", err)
				return
			}
			log.Println("[INCOMING-FROM-peer]", string(buffer[0:bytesRead]))
			if string(buffer[0:bytesRead]) == "Hello!" {

				log.Println("Successfully received [Hello!] message from peer ! ! !")
				log.Println("NOW ! Now you can write your message at [Enter message] ")
				log.Println()
			}
			//receiveStreamChan <- receiveStream
		}()
		log.Printf("Attempting to dial %s\n", remote)
		startTime = time.Now()
		log.Println("startTime", startTime)
		for i := 0; i < 10; i++ {
			conntoPeer, err = tr.Dial(context.Background(), remoteAddr, tlsConfDial, nil)
			if err == nil {
				//计算打洞结束时间
				endTime = time.Now()
				log.Println("endTime", endTime)
				break
			}
			log.Println("Dial failed, reattempting", err)
			time.Sleep(time.Second)
		}
		if conntoPeer == nil {
			log.Println("Unable to establish connection")
			return
		}

		log.Println("Successfully established connection with peer")
		sendStream, err := conntoPeer.OpenUniStreamSync(context.Background())
		if err != nil {
			log.Printf("Error with OpenUniStreamSync: %s\n", err)
			return
		}
		_, err = sendStream.Write([]byte("Hello!"))
		if err != nil {
			fmt.Printf("Failed to send message: %s\n", err)
		} else {
			log.Println("Successfully sent [Hello!] to ---> ", remote)
		}
		//sendStreamChan <- sendStream

	}()
	<-done // 等待通知，确保goroutine执行完毕
	// 计算时间差
	executionTime := endTime.Sub(startTime)
	return executionTime

	//executionTime := endTime.Sub(startTime)
	//放在<-done的结尾最后一行，第一次返回的值是的一个复数，不可用

	//定义互斥锁
	//var streamMutex sync.Mutex
	////从管道中提取流
	//receiveStream := <-receiveStreamChan
	//sendStream := <-sendStreamChan
	////双方都接收到hello之后，开始进行消息传输，在控制台输入消息
	//go func() {
	//	go func() {
	//		for {
	//			buffer := make([]byte, 1024)
	//			newRead, _ := receiveStream.Read(buffer)
	//			log.Println("[INCOMING-FROM-PEER]", string(buffer[:newRead]))
	//		}
	//	}()
	//	for {
	//		//在控制台输入消息
	//		reader := bufio.NewReader(os.Stdin)
	//		log.Print("[Enter message]")
	//		text, err := reader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("Failed to read input: %s\n", err)
	//			return
	//		}
	//
	//		streamMutex.Lock()
	//		_, err = sendStream.Write([]byte(text))
	//		streamMutex.Unlock()
	//		if err != nil {
	//			fmt.Printf("Failed to send message: %s\n", err)
	//			return
	//		}
	//	}
	//}()
}
