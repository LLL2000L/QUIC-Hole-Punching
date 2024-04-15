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
	"time"
)

func Tcp_server() {
	localAddress := ":9001"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	//cert, err := tgenerateSelfSignedCertificate()
	//if err != nil {
	//	fmt.Printf("Error generating certificate: %s\n", err)
	//	return
	//}
	//tlsConfListen := &tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	NextProtos:   []string{"tcp+tls1.3"},
	//	MinVersion:   tls.VersionTLS13,
	//	//MinVersion: tls.VersionTLS12,
	//	//MaxVersion: tls.VersionTLS12,
	//}
	//创建TLS监听器
	//listener, err := tls.Listen("tcp4", localAddress, tlsConfListen)

	// 中继服务器地监听址
	listener, err := net.Listen("tcp4", localAddress)
	if err != nil {
		fmt.Println("Failed to ListenTCP:", err)
		os.Exit(1)
	}

	//defer listener.Close()

	// 接受客户端连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept client connection:", err)
			continue
		} else {
			endTime := time.Now()
			log.Println("服务器endTime--->", endTime)
		}

		defer clientConn.Close()
		log.Println("Client connected:", clientConn.RemoteAddr())

		//处理进来的TLS
		//buffer := make([]byte, 1024)
		//bytesRead, err := clientConn.Read(buffer)
		//if err != nil {
		//	fmt.Println("Failed to read from client:", err)
		//	return
		//}
		//
		//incoming := string(buffer[0:bytesRead])
		//fmt.Println("[INCOMING]", incoming)
		//if incoming != "" {
		//	return
		//}
	}
}

func tgenerateSelfSignedCertificate() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"tcp-go"},
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
