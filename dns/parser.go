package dns

import (
	"net"
	"strings"
)

func parseName(b []byte, offset int) (string, int) {
	labels := []string{}
	cursor := offset

	for {
		isPtr := b[cursor]&0xC0 == 0xC0

		// Pointer always signifies end, RFC 1035 4.14
		if isPtr {
			pointerOffset := btoi16([]byte{
				b[cursor] & 0x3F,
				b[cursor+1],
			})

			l, _ := parseName(b, int(pointerOffset))
			labels = append(labels, l)
			cursor += 2

			break
		}

		labelLen := int(b[cursor])
		cursor += 1

		// Octet labels terminator
		if labelLen == 0 {
			break
		}

		l := string(b[cursor : cursor+labelLen])
		cursor += labelLen
		labels = append(labels, l)
	}

	return strings.Join(labels, "."), cursor - offset
}

func parseRR(b []byte, offset int, isQuery bool) (ResourceRecord, int) {
	var n int
	var rr ResourceRecord

	cursor := offset

	rr.Name, n = parseName(b, cursor)
	cursor += n

	rr.RRType = rrtypes[btoi16(b[cursor:cursor+2])]
	cursor += 2

	rr.Class = "IN"
	cursor += 2

	// For queries, we're done after the common fields
	if isQuery {
		return rr, cursor - offset
	}

	rr.TTL = btoi32(b[cursor : cursor+4])
	cursor += 4

	dataLen := int(btoi16(b[cursor : cursor+2]))
	cursor += 2

	switch rr.RRType {
	case "NS":
		rr.Nameserver, _ = parseName(b, cursor)
	case "A", "AAAA":
		rr.Address = net.IP(b[cursor : cursor+dataLen])
	}

	cursor += dataLen

	return rr, cursor - offset
}

func ParseDNSReponse(b []byte) *DNSMessage {
	response := DNSMessage{
		Header: Header{
			QueriesCount:     btoi16(b[4:6]),
			AnswersCount:     btoi16(b[6:8]),
			AuthoritiesCount: btoi16(b[8:10]),
			AdditionalCount:  btoi16(b[10:12]),
		},
	}

	cursor := 12
	for range response.Header.QueriesCount {
		rr, n := parseRR(b, cursor, true)
		cursor += n

		response.Queries = append(response.Queries, rr)
	}

	for range response.Header.AnswersCount {
		rr, n := parseRR(b, cursor, false)
		cursor += n

		response.Answers = append(response.Answers, rr)
	}

	for range response.Header.AuthoritiesCount {
		rr, n := parseRR(b, cursor, false)
		cursor += n

		response.Authorities = append(response.Authorities, rr)
	}

	for range response.Header.AdditionalCount {
		rr, n := parseRR(b, cursor, false)
		cursor += n

		response.Additional = append(response.Additional, rr)
	}

	return &response
}
