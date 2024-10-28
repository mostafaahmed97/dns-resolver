package dns

import (
	"net"
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
