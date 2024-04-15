package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type clientType struct {
	mux sync.Mutex
	//clients map[string][]net.TCPConn
	clients map[string][]net.Conn
}

var clients = clientType{
	//clients: make(map[string][]net.TCPConn),
	clients: make(map[string][]net.Conn),
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
	localAddress := ":9999"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	//cert, err := generateSelfSignedCertificate()
	//if err != nil {
	//	fmt.Printf("Error generating certificate: %s\n", err)
	//	return
	//}
	//tlsConfListen := &tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	NextProtos:   []string{"tcp-holepunch"},
	//	MinVersion:   tls.VersionTLS13,
	//}
	// 中继服务器地监听址

	// 创建TLS监听器
	//listener, err := tls.Listen("tcp4", localAddress, tlsConfListen)
	//listener, err := reuseport.Listen("tcp4", localAddress)
	listener, err := net.Listen("tcp4", localAddress)
	if err != nil {
		fmt.Println("Failed to ListenTCP:", err)
		os.Exit(1)
	}
	defer listener.Close()

	// 接受客户端连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept client connection:", err)
			continue
		}
		//defer clientConn.Close()
		log.Println("Client connected:", clientConn.RemoteAddr())
		go handleSession(clientConn)
	}
}

func handleSession(sess net.Conn) {
	//获取客户端IP地址和端口
	remoteAddr := sess.RemoteAddr().String()
	// Send remote address to the client

	buffer := make([]byte, 1024)
	bytesRead, err := sess.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read from client:", err)
		return
	}

	incoming := string(buffer[0:bytesRead])
	fmt.Println("[INCOMING]", incoming)
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
				//将两端的地址和port都告诉双方
				data := client + "," + resp
				_, err = conn.Write([]byte(data))

				if err != nil {
					fmt.Printf("[INFO] Error responding to %s: %s\n", client, err)
					return
				} else {
					fmt.Printf("[INFO] Responded to %s: %s\n", client, resp)
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
