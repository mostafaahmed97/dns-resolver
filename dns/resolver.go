package dns

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ResolveFromRoot(domain string, root string) string {
	target := root

	for {

		conn, err := net.Dial("udp", target+":53")
		if err != nil {
			fmt.Println("connection to server failed")
			os.Exit(1)
		}
		defer conn.Close()

		fmt.Printf("Asking: %s\n", target)
		message := NewDNSMessage(domain)
		conn.Write(message)

		b := make([]byte, 1024)
		n, _ := conn.Read(b)

		response := ParseDNSReponse(b[:n])

		if response.Header.AnswersCount > 0 {
			fmt.Printf("Found: \n")
			fmt.Printf("\tAnswer: %s\n", response.Answers[0].Address.String())
			return response.Answers[0].Address.String()
		}

		// Resolve NS with no additional records
		if response.Header.AuthoritiesCount > 0 && response.Header.AdditionalCount == 0 {
			nameserver := response.Authorities[0].Nameserver
			target = ResolveFromRoot(nameserver, root)
		}

		// Ask the next server
		for _, rr := range response.Additional {
			if rr.RRType == "A" {
				fmt.Printf("Found: \n")
				fmt.Printf("\tAuthority: %s\n", response.Authorities[0].Nameserver)
				fmt.Printf("\tFor: %s\n", response.Authorities[0].Nameserver)
				fmt.Printf("\tAt: %s\n", rr.Address.String())

				target = rr.Address.String()
				break
			}
		}

		fmt.Println(strings.Repeat("-", 30))
	}

}
