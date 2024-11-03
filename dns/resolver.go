package dns

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ResolveFromRoot(domain string, root string, depth int) string {
	target := root

	for {
		indent := strings.Repeat(" ", depth*2)
		fmt.Printf(
			"%s➜ Querying %s for %s\n",
			indent,
			target,
			domain,
		)

		conn, err := net.Dial("udp", target+":53")
		if err != nil {
			fmt.Println("connection to server failed")
			os.Exit(1)
		}
		defer conn.Close()

		message := NewDNSMessage(domain)
		conn.Write(message)

		b := make([]byte, 1024)
		n, _ := conn.Read(b)

		response := ParseDNSReponse(b[:n])
		fmt.Printf(
			"%s∞ Got: %d answers, %d authorities, %d additional\n",
			indent,
			response.Header.AnswersCount,
			response.Header.AuthoritiesCount,
			response.Header.AdditionalCount,
		)

		if response.Header.AuthoritiesCount > 0 {
			fmt.Printf(
				"%s\t Authorities: \n\t\t%s\n\t \n",
				indent,
				strings.Join(response.AuthorityNames(), ", \n"+indent+"\t\t"),
			)
		}

		if response.Header.AnswersCount > 0 {
			fmt.Printf(
				"%s✓ Found: %s -> %s\n",
				indent,
				domain,
				"["+strings.Join(response.AnswerAddresses(), ", ")+"]",
			)

			return response.Answers[0].Address.String()
		}

		// Resolve NS with no additional records
		if response.Header.AuthoritiesCount > 0 && response.Header.AdditionalCount == 0 {
			nameserver := response.Authorities[0].Nameserver
			fmt.Printf(
				"%s⤷ Need to resolve nameserver: %s\n",
				indent,
				nameserver,
			)

			target = ResolveFromRoot(nameserver, root, depth+1)
		}

		// Ask the next server
		for _, rr := range response.Additional {
			if rr.RRType == "A" {
				fmt.Printf(
					"%s⤷ Following referral to: %s (%s)\n",
					indent,
					rr.Address.String(),
					rr.Name,
				)

				target = rr.Address.String()
				break
			}
		}

		fmt.Printf("%s%s\n", indent, strings.Repeat("-", 20))
	}
}
