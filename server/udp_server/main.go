package main

import (
    "fmt"
    "net"
)

func main() {
    addr := net.UDPAddr{
        Port: 8080,
        IP:   net.ParseIP("127.0.0.1"),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    buffer := make([]byte, 1024)
    for {
        _, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading from UDP connection:", err)
            continue
        }
        fmt.Printf("Received %s from %s\n", string(buffer), remoteAddr)
    }
}