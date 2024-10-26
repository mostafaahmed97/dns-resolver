package dns

import "net"

// returns host and number of bytes read
func parseHost(b []byte, offset int, dataLength int) (string, int) {
	cursor := offset
	host := ""

	for {
		if dataLength > 0 && (cursor-offset) >= dataLength {
			break
		}

		isPointer := b[cursor] == 0xc0
		part := ""

		// part can be a pointer referencing a name encountered before
		if isPointer {
			part, _ = parseHost(b, int(b[cursor+1]), 0)

			cursor += 2

			// assume name that starts with pointer
			// is only a pointer
			if host == "" {
				return part, 2
			}
		} else {
			len := int(b[cursor])
			cursor += 1

			part = string(b[cursor : cursor+len])
			cursor += len
		}

		if host == "" {
			host = part
		} else {
			host = host + "." + part
		}

		if b[cursor] == 0x00 {
			cursor += 1
			break
		}
	}

	return host, cursor - offset
}

func parseRR(b []byte, offset int) (ResourceRecord, int) {
	cursor := offset

	name, n := parseHost(b, cursor, 0)
	cursor += n

	rrtype := rrtypes[btoi16(b[cursor:cursor+2])]
	cursor += 2

	class := "IN"
	cursor += 2

	return ResourceRecord{
		Name:   name,
		RRType: rrtype,
		Class:  class,
	}, cursor - offset
}

type info struct {
	ttl        uint32
	nameserver string
	address    net.IP
}

func parseInfo(b []byte, offset int, hasAddress bool) (info, int) {
	cursor := offset

	ttl := btoi32(b[cursor : cursor+4])
	cursor += 4

	dataLen := int(btoi16(b[cursor : cursor+2]))
	cursor += 2

	address := net.IP{}
	nameserver := ""

	if hasAddress {
		address = net.IP(b[cursor : cursor+dataLen])
	} else {
		nameserver, _ = parseHost(b, cursor, dataLen)
	}

	cursor += dataLen

	return info{
		ttl:        ttl,
		nameserver: nameserver,
		address:    address,
	}, cursor - offset
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

	curr := 12
	for range response.Header.QueriesCount {
		rr, n := parseRR(b, curr)
		curr += n

		response.Queries = append(response.Queries,
			Query{Record: rr},
		)
	}

	for range response.Header.AnswersCount {
		rr, n := parseRR(b, curr)
		curr += n

		info, n := parseInfo(b, curr, true)
		curr += n

		response.Answers = append(
			response.Answers,
			Answer{
				Record:  rr,
				TTL:     info.ttl,
				Address: info.address,
			})
	}

	for range response.Header.AuthoritiesCount {
		rr, n := parseRR(b, curr)
		curr += n

		info, n := parseInfo(b, curr, false)
		curr += n

		response.Authorities = append(
			response.Authorities,
			Authority{
				Record:     rr,
				TTL:        info.ttl,
				Nameserver: info.nameserver,
			})
	}

	for range response.Header.AdditionalCount {
		rr, n := parseRR(b, curr)
		curr += n

		info, n := parseInfo(b, curr, true)
		curr += n

		response.Additional = append(
			response.Additional,
			Additional{
				Record:  rr,
				TTL:     info.ttl,
				Address: info.address,
			},
		)
	}

	return &response
}
