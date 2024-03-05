package main

import (
	"fmt"

	"github.com/miekg/dns"
)

func main() {
	fmt.Println("very first line of Go DN server practice")
	// domain := "ajiyba.com"
	// resolver(domain, dns.TypeA)

	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":53",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port 53")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
