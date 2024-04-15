package main

import (
	"os"
	"time"
)

func main() {
	time.Sleep(time.Second * 2)
	cmd := os.Args[1]
	switch cmd {
	case "quic_c":
		Quic_client()
	case "quic_s":
		Quic_server()
	}
}
