package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type clientType struct {
	mux     sync.Mutex
	clients map[string][]quic.Connection
}

var clients = clientType{
	clients: make(map[string][]quic.Connection),
}

func (c *clientType) keys(filter string) string {
	output := []string{}
	for key := range c.clients {
		if key != filter {
			output = append(output, key)
		}
	}

	return strings.Join(output, ",")
}

// Server --
func Server() {
	localAddress := ":8888"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	addr, _ := net.ResolveUDPAddr("udp4", localAddress)
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("连接出现错误", err)
	}

	// 增加 UDP 缓冲区大小
	err = udpConn.SetReadBuffer(2048 * 1024) // 设置为期望的大小
	if err != nil {
		fmt.Printf("设置 UDP 缓冲区大小出现错误: %s\n", err)
		return
	}

	// QUIC 配置
	// 生成自签名证书
	cert, err := generateSelfSignedCertificate()
	if err != nil {
		fmt.Printf("生成证书出现错误: %s\n", err)
		return
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-holepunch"},
	}

	//创建quic连接（单个 UDP 套接字上运行的 QUIC 连接）
	tr := quic.Transport{
		Conn: udpConn,
	}
	listener, err := tr.ListenEarly(tlsConf, &quic.Config{Allow0RTT: true})
	//listener, err := quic.Listen(udpConn, tlsConf, nil)
	if err != nil {
		fmt.Printf("连接出现错误: %s\n", err)
		return
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Printf("接受连接时出现错误: %s\n", err)
			continue
		}
		go handleSession(conn)
	}
}

func handleSession(sess quic.Connection) {
	log.Println("Client connected:", sess.RemoteAddr())
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		fmt.Printf("接受 Stream 时出现错误: %s\n", err)
		return
	}
	defer stream.Close()

	remoteAddr := sess.RemoteAddr().String()
	// Send remote address to the client

	buffer := make([]byte, 1024)
	bytesRead, err := stream.Read(buffer)
	if err != nil {
		//fmt.Printf("读取数据时出现错误: %s\n", err)
		//return
		panic(err)
	}

	incoming := string(buffer[0:bytesRead])
	log.Println("[INCOMING]", incoming)
	if incoming != "register" {
		return
	}

	//将客户端和连接映射，遍历每个连接以便于将对端IP和port发给客户端
	//remoteAddr := conn.RemoteAddr().String()
	clients.clients[remoteAddr] = append(clients.clients[remoteAddr], sess)

	for client, conns := range clients.clients {
		resp := clients.keys(client)
		if len(resp) > 0 {
			//遍历连接，获取与对端不同的连接
			for _, conn := range conns {
				clients.mux.Lock()
				s, err := conn.OpenStreamSync(context.Background())
				if err != nil {
					fmt.Printf("[INFO] Error opening stream to %s: %s\n", client, err)
					continue
				}
				//将两端的地址和port都告诉双方
				data := client + "," + resp
				_, err = s.Write([]byte(data))

				if err != nil {
					fmt.Printf("[INFO] Error responding to %s: %s\n", client, err)
					return
				} else {
					log.Printf("[INFO] Responded to %s: %s\n", client, resp)
				}
				clients.mux.Unlock()
			}
		}
	}

}

func generateSelfSignedCertificate() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"quic-go"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	cert, err := tls.X509KeyPair(certPEM, privPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}
