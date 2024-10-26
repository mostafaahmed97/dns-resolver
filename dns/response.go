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
}

type Query struct {
	Record ResourceRecord
}

type Answer struct {
	Record  ResourceRecord
	Address net.IP
	TTL     uint32
}

type Authority struct {
	Record     ResourceRecord
	Nameserver string
	TTL        uint32
}

type Additional struct {
	Record  ResourceRecord
	Address net.IP
	TTL     uint32
}

type DNSMessage struct {
	Header Header

	Queries     []Query
	Answers     []Answer
	Authorities []Authority
	Additional  []Additional
}
