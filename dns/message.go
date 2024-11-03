package dns

import (
	"net"
	"strings"
)

var rrtypes = map[uint16]string{
	1:  "A",
	2:  "NS",
	28: "AAAA",
}

type Header struct {
	TransactionId    int
	QueriesCount     uint16
	AnswersCount     uint16
	AuthoritiesCount uint16
	AdditionalCount  uint16
}

type ResourceRecord struct {
	Name   string
	RRType string
	Class  string
	TTL    uint32

	// Available on RRs with type `NS`
	Nameserver string

	// Available on RRs with type `A` & `AAAA`
	Address net.IP
}

type DNSMessage struct {
	Header Header

	Queries     []ResourceRecord
	Answers     []ResourceRecord
	Authorities []ResourceRecord
	Additional  []ResourceRecord
}

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
