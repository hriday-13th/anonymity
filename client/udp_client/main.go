package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	message := "Hello there UDP server"
	_, err = conn.Write([]byte(message))

	if err != nil {
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("Message sent to UDP server")
}