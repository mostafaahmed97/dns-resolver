package dns

import (
	"fmt"
	"net"
	"os"
)

func ResolveURLFromRoot(url string, root string) string {
	target := root + ":53"

	for {
		msg := NewDNSMessage(url)

		address, err := net.ResolveUDPAddr("udp", target)
		if err != nil {
			fmt.Println("failed to resolve server address", target, err)
			os.Exit(1)
		}

		conn, err := net.DialUDP("udp", nil, address)
		if err != nil {
			fmt.Println("connection to server failed")
			os.Exit(1)
		}

		defer conn.Close()

		conn.Write(msg.bytes)

		b := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(b)

		response := ParseDNSReponse(b[:n])

		if response.anscount > 0 {
			return response.answers[0].address
		}

		if response.authcount == 0 {
			fmt.Println("no auth found for", url)
			os.Exit(1)
		}

		// Ask the next server
		for _, rr := range response.additional {
			if rr.rrtype == "A" {
				target = rr.address + ":53"
			}
		}
	}

}
