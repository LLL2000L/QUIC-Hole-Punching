package main

import (
	"os"
)

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "tcp_c":
		Tcp_client()
	case "tcp_s":
		Tcp_server()
	}
}
