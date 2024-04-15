package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
	"net"
	"os"
	"time"
)

func Quic_client() {
	serverAddress := os.Args[2]
	localAddress := ":8002" // default port
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}
	server, _ := net.ResolveUDPAddr("udp", serverAddress)
	local, _ := net.ResolveUDPAddr("udp", localAddress)
	udpConn, err := net.ListenUDP("udp", local)
	if err != nil {
		fmt.Printf("UDP connection error: %s\n", err)
		panic(err)
	}

	tr := quic.Transport{
		Conn: udpConn,
	}

	tlsConfDial := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic"},
	}

	startTime := time.Now()
	conn, err := tr.Dial(context.Background(), server, tlsConfDial, nil)
	if err != nil {
		fmt.Printf("failed: %s\n", err)
		panic(err)
	} else {
		endTime := time.Now()
		log.Println("startTime--->", startTime)
		log.Println("endTime--->", endTime)
		executionTime := endTime.Sub(startTime)
		log.Println("TotalTime", executionTime)
	}

	conn.CloseWithError(0x42, "error 0x42 occurred")

}
