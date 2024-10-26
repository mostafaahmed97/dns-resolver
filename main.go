package main

import (
	"fmt"
	"os"

	"github.com/mostafaahmed97/dns-resolver/dns"
)

var roots = map[string]string{
	"a": "198.41.0.4",
	"b": "170.247.170.2",
	"c": "192.33.4.12",
	"d": "199.7.91.13",
	"e": "192.203.230.10",
	"f": "192.5.5.241",
	"g": "192.112.36.4",
	"h": "198.97.190.53",
	"i": "192.36.148.17",
	"j": "192.58.128.30",
	"k": "193.0.14.129",
	"l": "199.7.83.42",
	"m": "202.12.27.33",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing url")
		os.Exit(1)
	}

	r := "a"
	url := os.Args[1]

	if len(os.Args) >= 3 {
		r = os.Args[2]
	}

	fmt.Printf("Resolving: %s from Root Server: %s(%s)\n", url, r, roots[r])
	ip := dns.ResolveFromRoot(url, roots[r])

	fmt.Printf("%s is at %s\n", url, ip)
}
