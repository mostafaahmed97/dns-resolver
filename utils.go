package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func btoi(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
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
