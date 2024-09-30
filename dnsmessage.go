package main

import "strings"

type DNSMessage struct {
	query string
	bytes []byte
}

func NewDNSMessage(query string) *DNSMessage {
	msg := DNSMessage{query: query}

	txnId := []byte{0xaa, 0xaa}
	flags := []byte{0x01, 0x00}
	qryCount := []byte{0x00, 0x01}
	ansCount := []byte{0x00, 0x00}
	authRRCount := []byte{0x00, 0x00}
	addnRRCount := []byte{0x00, 0x00}

	headerSections := [][]byte{
		txnId, flags, qryCount, ansCount, authRRCount, addnRRCount,
	}

	bytes := []byte{}

	// DNS Packet Header
	for _, section := range headerSections {
		bytes = append(bytes, section...)
	}

	// DNS message
	// Domain, added as <length of part><part>
	urlparts := strings.Split(msg.query, ".")
	for _, part := range urlparts {
		bytes = append(bytes, byte(len(part)))
		bytes = append(bytes, part...)
	}

	// Domain termination
	bytes = append(bytes, 0x00)

	// Type A, class IN
	recordType := []byte{0x00, 0x01}
	recordClass := []byte{0x00, 0x01}

	bytes = append(bytes, recordType...)
	bytes = append(bytes, recordClass...)

	msg.bytes = bytes

	return &msg
}
