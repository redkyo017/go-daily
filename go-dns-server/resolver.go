package main

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
)

func resolver(domain string, qtype uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		log.Fatalf("[ERROR] : %v\n", err)
		return nil
	}

	if response == nil {
		log.Fatalf("[ERROR] : no response from server \n")
		return nil
	}

	for _, answer := range response.Answer {
		log.Println("con co")
		fmt.Printf("%s\n", answer.String())
	}

	return response.Answer
}
