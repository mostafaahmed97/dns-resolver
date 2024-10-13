package dns

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ResolveURLFromRoot(url string, root string) string {
	target := root

	for {
		fmt.Printf("Asking: %s\n", target)
		msg := NewDNSMessage(url)

		address, err := net.ResolveUDPAddr("udp", target+":53")
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

		if response.AnsCount > 0 {
			fmt.Printf("Found: \n")
			fmt.Printf("\tAnswer: %s\n", response.Answers[0].Address.String())
			return response.Answers[0].Address.String()
		}

		if response.AuthCount == 0 {
			fmt.Println("no auth found for", url)
			os.Exit(1)
		}

		// Ask the next server
		for _, rr := range response.Additional {
			if rr.Record.RRType == "A" {
				fmt.Printf("Found: \n")
				fmt.Printf("\tAuthority: %s\n", response.Authorities[0].Nameserver)
				fmt.Printf("\tFor: %s\n", response.Authorities[0].Record.Name)
				fmt.Printf("\tAt: %s\n", rr.Address.String())

				target = rr.Address.String()
				break
			}
		}

		fmt.Println(strings.Repeat("-", 30))
	}

}
