package main

import (
	"fmt"
	"net"
)

var roota = "198.41.0.4:53"

func main() {
	dnspacket := []byte{}

	// DNS Packet Header
	// TransactionID
	dnspacket = append(dnspacket, 0xaa, 0xaa)

	// Flags
	dnspacket = append(dnspacket, 0x01, 0x00)

	// No. of questions
	dnspacket = append(dnspacket, 0x00, 0x01)

	// No. of answers
	dnspacket = append(dnspacket, 0x00, 0x00)

	// No. of authority RRs
	dnspacket = append(dnspacket, 0x00, 0x00)

	// No. of additional RRs
	dnspacket = append(dnspacket, 0x00, 0x00)

	// DNS message

	// Domain, added as <length of part><part>
	// 6 google
	dnspacket = append(dnspacket, byte(len("google")))
	dnspacket = append(dnspacket, "google"...)
	// dnspacket = append(dnspacket, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65)

	// 3 com
	dnspacket = append(dnspacket, byte(len("com")))
	dnspacket = append(dnspacket, "com"...)
	// dnspacket = append(dnspacket, 0x63, 0x6f, 0x6d)

	// Domain termination
	dnspacket = append(dnspacket, 0x00)

	// Type, A record
	dnspacket = append(dnspacket, 0x00, 0x01)

	// Class
	dnspacket = append(dnspacket, 0x00, 0x01)

	remoteaddr, _ := net.ResolveUDPAddr("udp", roota)

	conn, err := net.DialUDP("udp", nil, remoteaddr)
	if err != nil {
		fmt.Println("failed to connect,", err)
	}

	fmt.Println(dnspacket)

	defer conn.Close()
	conn.Write(dnspacket)
	fmt.Println("sent the packet")

	fmt.Println("reading from the conn")
	response := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(response)

	fmt.Println("read sth,", n, addr, string(response[:n]))
}
