package main

import (
	"fmt"
	"net"
)

var roota = "198.41.0.4:53"
var comtld = "192.41.162.30:53"
var googleas = "216.239.34.10:53"

func main() {
	msg := NewDNSMessage("google.com")

	remoteaddr, _ := net.ResolveUDPAddr("udp", comtld)

	conn, err := net.DialUDP("udp", nil, remoteaddr)
	if err != nil {
		fmt.Println("failed to connect,", err)
	}

	defer conn.Close()
	conn.Write(msg.bytes)
	fmt.Println("sent the packet")

	fmt.Println("reading from the conn")

	b := make([]byte, 1024)

	n, _, _ := conn.ReadFromUDP(b)
	bytes := b[:n]

	response := FromBytes(bytes)
	fmt.Println(response)
}
