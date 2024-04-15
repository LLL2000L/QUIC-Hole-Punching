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
	"math/big"
	"net"
	"os"
	"time"
)

func Quic_server() {
	localAddress := ":8001"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	addr, _ := net.ResolveUDPAddr("udp4", localAddress)
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("failed connect", err)
	}

	cert, err := generateSelfSignedCertificate()
	if err != nil {
		fmt.Printf("failed: %s\n", err)
		return
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic"},
	}

	tr := quic.Transport{
		Conn: udpConn,
	}
	listener, err := tr.Listen(tlsConf, nil)
	if err != nil {
		fmt.Printf("failed: %s\n", err)
		return
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Printf("failed connect: %s\n", err)
			continue
		}
		fmt.Println(conn)
		//go handleSession(conn)
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
