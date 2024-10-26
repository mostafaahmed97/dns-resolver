package dns

import (
	"strings"
)

func NewDNSMessage(query string) []byte {
	message := []byte{}

	header := []byte{
		// Transaction ID
		0xaa, 0xaa,

		// Flags
		0x01, 0x00,

		// Queries, Answers, Authority Nameservers & Additional Count
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}

	// Domain name, encoded as <PART_LENGTH><PART>
	qname := []byte{}
	parts := strings.Split(query, ".")

	for _, p := range parts {
		l := byte(len(p))

		qname = append(qname, l)
		qname = append(qname, p...)
	}

	// Termination byte
	qname = append(qname, 0x00)

	// QTYPE,  1 = A record (IPv4 address)
	qtype := []byte{0x00, 0x01}

	// QCLASS, 1 = IN (Internet)
	qclass := []byte{0x00, 0x01}

	message = append(message, header...)
	message = append(message, qname...)
	message = append(message, qtype...)
	message = append(message, qclass...)

	return message
}
