package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

var rrtypes = map[uint16]string{
	1:  "A",
	2:  "NS",
	28: "AAA",
}

type Query struct {
	name  string
	qtype string
	class string
}

type Answer struct {
	name    string
	rrtype  string
	class   string
	ttl     int
	address string
}

type Authority struct {
	name   string
	rrtype string
	ns     string
	class  string
	ttl    int
}

type Additional struct {
	name    string
	rrtype  string
	class   string
	address string
	ttl     int
}

type DNSResponse struct {
	qcount    uint16
	anscount  uint16
	authcount uint16
	addcount  uint16

	queries    Query
	answers    []Answer
	autorities []Authority
	additional []Additional
}

func btoi(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// returns host and number of bytes read
func btohost(b []byte, offset int, datalength int) (string, int) {
	curr := offset
	host := ""

	for {
		if datalength > 0 && (curr-offset) >= datalength {
			break
		}

		ispointer := b[curr] == 0xc0
		part := ""

		// part can be a pointer referencing a name encountered before
		//
		if ispointer {
			part, _ = btohost(b, int(b[curr+1]), 0)

			// first byte of pointer
			curr += 2
		} else {
			partlen := int(b[curr])
			curr += 1

			part = string(b[curr : curr+partlen])
			curr += partlen
		}

		if host == "" {
			host = part
		} else {
			host = host + "." + part
		}

		if b[curr] == 0x00 {
			// skips terminator or second byte of pointer
			curr += 1
			break
		}
	}

	return host, curr - offset
}

func getipv4addr(b []byte) string {
	octets := []string{
		fmt.Sprintf("%d", int(b[0])),
		fmt.Sprintf("%d", int(b[1])),
		fmt.Sprintf("%d", int(b[2])),
		fmt.Sprintf("%d", int(b[3])),
	}

	return strings.Join(octets, ".")
}

func getipv6addr(b []byte) string {
	return "PLACEHOLDER IPV6"
}

func FromBytes(b []byte) *DNSResponse {
	response := DNSResponse{
		qcount:    btoi(b[4:6]),
		anscount:  btoi(b[6:8]),
		authcount: btoi(b[8:10]),
		addcount:  btoi(b[10:12]),
	}

	curr := 12
	for range response.qcount {

		name, read := btohost(b, curr, 0)
		curr += read

		qtype := rrtypes[btoi(b[curr:curr+2])]
		curr += 2

		class := ""
		if btoi(b[curr:curr+2]) == 1 {
			class = "IN"
		}
		curr += 2

		response.queries = Query{
			name:  name,
			qtype: qtype,
			class: class,
		}
	}

	for range response.anscount {
		name, _ := btohost(b, int(b[curr+1]), 0)
		curr += 2

		rrtype := rrtypes[btoi(b[curr:curr+2])]
		curr += 2

		// skip class
		curr += 2

		ttl := int(btoi(b[curr : curr+4]))
		curr += 4

		// skip data length
		curr += 2

		address := ""
		if rrtype == "A" {
			address = getipv4addr(b[curr : curr+4])
			curr += 4
		} else if rrtype == "AAA" {
			address = getipv6addr(b[curr : curr+16])
			curr += 16
		}

		response.answers = append(
			response.answers,
			Answer{
				ttl:     ttl,
				name:    name,
				rrtype:  rrtype,
				address: address,
			})
	}

	for range response.authcount {
		// name is compressed, get offset from start and
		// extract directly
		name, _ := btohost(b, int(b[curr+1]), 0)
		curr += 2

		rrtype := rrtypes[btoi(b[curr:curr+2])]
		curr += 2

		//skip class for now
		curr += 2

		ttl := int(btoi(b[curr : curr+4]))
		curr += 4

		// skip data length
		datalength := int(btoi(b[curr : curr+2]))
		curr += 2

		ns, n := btohost(b, curr, datalength)
		curr += n

		response.autorities = append(
			response.autorities,
			Authority{
				name:   name,
				ns:     ns,
				rrtype: rrtype,
				ttl:    ttl,
			})
	}

	for range response.addcount {
		name, _ := btohost(b, curr, 0)
		curr += 2

		rrtype := rrtypes[btoi(b[curr:curr+2])]
		curr += 2

		// skip class
		curr += 2

		ttl := int(btoi(b[curr : curr+4]))
		curr += 4

		// skip data length
		curr += 2

		address := ""
		if rrtype == "A" {
			address = getipv4addr(b[curr : curr+4])
			curr += 4
		} else if rrtype == "AAA" {
			address = getipv6addr(b[curr : curr+16])
			curr += 16
		}

		response.additional = append(
			response.additional,
			Additional{
				name:    name,
				address: address,
				rrtype:  rrtype,
				ttl:     ttl,
			},
		)
	}

	return &response
}
