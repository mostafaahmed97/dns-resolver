package main

import (
	"fmt"
	"net"
)

var roota = "198.41.0.4:53"

func main() {
	msg := NewDNSMessage("spotify.net")

	remoteaddr, _ := net.ResolveUDPAddr("udp", roota)

	conn, err := net.DialUDP("udp", nil, remoteaddr)
	if err != nil {
		fmt.Println("failed to connect,", err)
	}

	defer conn.Close()
	conn.Write(msg.bytes)
	fmt.Println("sent the packet")

	fmt.Println("reading from the conn")
	response := make([]byte, 1024)

	n, addr, _ := conn.ReadFromUDP(response)
	fmt.Println("read sth,", n, addr, string(response[:n]))
}
