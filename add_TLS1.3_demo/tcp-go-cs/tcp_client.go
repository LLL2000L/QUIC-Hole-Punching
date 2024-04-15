package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func Tcp_client() {
	serverAddress := os.Args[2]
	localAddress := ":9002" // default port
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}
	fmt.Println("localAddress", localAddress)

	// 创建TLS连接
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true, // 不验证证书
	//	NextProtos:         []string{"tcp+tls1.3"},
	//	MinVersion:         tls.VersionTLS13,
	//	//MinVersion: tls.VersionTLS12,
	//	//MaxVersion: tls.VersionTLS12,
	//}

	//conn, err := net.Dial("tcp", serverAddress)
	//conn, err := tls.Dial("tcp", serverAddress, tlsConfig)

	startTime := time.Now()
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	} else {
		endTime := time.Now()
		log.Println("startTime--->", startTime)
		log.Println("endTime--->", endTime)
		executionTime := endTime.Sub(startTime)
		log.Println("TotalTime", executionTime)
	}
	conn.Close()

}
