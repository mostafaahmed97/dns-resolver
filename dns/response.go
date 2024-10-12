package dns

var rrtypes = map[uint16]string{
	1:  "A",
	2:  "NS",
	28: "AAAA",
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

	queries     []Query
	answers     []Answer
	authorities []Authority
	additional  []Additional
}

func (r DNSResponse) hasAnswer() bool {
	return r.anscount > 0
}

// returns host and number of bytes read
func btohost(b []byte, offset int, dataLength int) (string, int) {
	cursor := offset
	host := ""

	for {
		if dataLength > 0 && (cursor-offset) >= dataLength {
			break
		}

		isPointer := b[cursor] == 0xc0
		label := ""

		// part can be a pointer referencing a name encountered before
		if isPointer {
			label, _ = btohost(b, int(b[cursor+1]), 0)

			// first byte of pointer
			cursor += 2
		} else {
			len := int(b[cursor])
			cursor += 1

			label = string(b[cursor : cursor+len])
			cursor += len
		}

		if host == "" {
			host = label
		} else {
			host = host + "." + label
		}

		if b[cursor] == 0x00 {
			cursor += 1
			break
		}
	}

	return host, cursor - offset
}

func ParseDNSReponse(b []byte) *DNSResponse {
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

		response.queries = append(response.queries,
			Query{
				name:  name,
				qtype: qtype,
				class: class,
			},
		)
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
			address = parseIPv4Addr(b[curr : curr+4])
			curr += 4
		} else if rrtype == "AAA" {
			address = parseIPv6Addr(b[curr : curr+16])
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

		response.authorities = append(
			response.authorities,
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
			address = parseIPv4Addr(b[curr : curr+4])
			curr += 4
		} else if rrtype == "AAAA" {
			address = parseIPv6Addr(b[curr : curr+16])
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
