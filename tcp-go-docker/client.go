package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Client --
func Client() {
	signalAddress := os.Args[2]
	localAddress := ":9001" // default port
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}

	//-------------------------------------上述代码是公用部分（client-relay || client-client）-------------------------------------
	type AddrTime struct {
		PeerAddr  string
		StartTime time.Time
	}
	peerAddrChan := make(chan AddrTime) // 创建一个通道用于传递 peerAddr 和 startTime

	go func() {
		// 函数执行完成后减少一个等待的 goroutine
		peerAddr, startTime := register(localAddress, signalAddress)       // 在线程1中执行 register() 函数
		peerAddrChan <- AddrTime{PeerAddr: peerAddr, StartTime: startTime} // 将 peerAddr 和 startTime 封装成结构体并发送到通道中
	}()

	addrTime := <-peerAddrChan // 从通道中接收 peerAddr 和 startTime 值
	peerAddr := addrTime.PeerAddr
	startTime := addrTime.StartTime

	//先使得双方都收到对端公网地址再打洞
	//time.Sleep(time.Second * 1)
	go func() {
		// 记录结束时间
		endTime := holepunch(localAddress, peerAddr) // 在线程2中使用 peerAddr 值执行 chatter 函数
		log.Println("startTime", startTime)
		log.Println("endTime", endTime)
	}()

	time.Sleep(time.Second * 20) //等待100s主线程退出
	log.Println("Main goroutine exit")
}

// 客户端和中继连接，并返回公网地址和端口号
func register(localAddress string, signalAddress string) (string, time.Time) {

	// 创建TLS连接
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true, // 不验证证书
	//	NextProtos:         []string{"tcp-holepunch"},
	//	MinVersion:         tls.VersionTLS13,
	//}
	// 创建TCP拨号连接
	conn, err := reuseport.Dial("tcp", localAddress, signalAddress)

	//dconn, err := net.Dial("tcp", signalAddress)
	//conn, err := tls.Dial("tcp", signalAddress, tlsConfig)
	//dconn, err := reuseport.Dial("tcp", localAddress, signalAddress)

	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close() // 延迟关闭连接

	var registerTime time.Time

	//向服务器发送【register】
	go func() {
		registerTime = time.Now()
		//fmt.Printf("开始给服务器发送的时间：", time.Now())
		bytesWritten, err := conn.Write([]byte("register"))
		if err != nil {
			panic(err)
			fmt.Printf("发送数据出现错误: %s\n", err)
		}
		fmt.Println(bytesWritten, " bytes written")
	}()

	var peerAddr string // 声明 peerAddr 变量
	for {
		fmt.Println("listening")

		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}

		data := string(buffer[0:bytesRead])
		addresses := strings.Split(data, ",")
		localAddr := addresses[0]
		peerAddr = addresses[1]
		fmt.Println("自己的公网IP:port--->", localAddr, time.Now())
		fmt.Println("对端的公网IP:port--->", peerAddr)
		break
	}
	//接受从服务器发送过来的对端地址和端口号
	return peerAddr, registerTime
}

// 客户端之间的连接（打洞）
// TCP打洞双方的客户端不需要监听，只需要dial就可以了
func holepunch(localAddress string, remote string) time.Time {
	var endTime time.Time
	// 记录holepunch开始时间
	//var startTime time.Time

	done := make(chan struct{}) // 创建一个用于通知函数结束的channel
	go func() {
		var conntoPeer net.Conn
		log.Printf("Attempting to dial %s\n", remote)
		//计算打洞开始时间
		//startTime = time.Now() //打洞时间从这里开始的目的是不能够忽略掉dial失败的时间
		for i := 0; i < 20; i++ {
			var err error
			conntoPeer, err = reuseport.Dial("tcp", localAddress, remote)
			//dconn, err := net.Dial("tcp", remote)
			if err == nil {
				endTime = time.Now()
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

		_, err := conntoPeer.Write([]byte("Hello!"))
		if err != nil {
			fmt.Printf("Failed to send message: %s\n", err)
		} else {
			log.Println("Successfully sent [Hello!] to ---> ", remote, time.Now())
		}

		go func() {
			defer close(done)
			log.Printf("Receiving for connection from %s\n", remote)
			buffer := make([]byte, 1024)
			bytesRead, err := conntoPeer.Read(buffer)
			if err != nil {
				fmt.Printf("Failed to read the data from peer: %s\n", err)
				return
			} else {
				log.Println("[INCOMING-FROM-peer]-->[Hello!]", time.Now())
			}
			log.Println("[INCOMING-FROM-peer]", string(buffer[0:bytesRead]))
			if string(buffer[0:bytesRead]) == "Hello!" {
				log.Println("Successfully received [Hello!] message from peer ! ! !")
				log.Println("NOW ! Now you can write your message at [Enter message] ")
				log.Println()
			}
		}()
	}()

	<-done // 等待通知，确保goroutine执行完毕
	//log.Println("startTime", startTime)
	//log.Println("endTime", endTime)
	//executionTime := endTime.Sub(startTime)
	//return executionTime
	return endTime
}
