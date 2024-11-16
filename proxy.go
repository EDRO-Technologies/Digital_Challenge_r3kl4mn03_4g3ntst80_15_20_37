package main

import (
	"errors"

	"github.com/miekg/dns"
)

func processRequest(request *dns.Msg) (*dns.Msg, error) {
	switch request.Opcode {
	case dns.OpcodeQuery:
		if len(request.Question) > 0 {
			dnsServer := "1.1.1.1:53"
			dnsClient := new(dns.Client)
			dnsClient.Net = "udp"
			response, _, err := dnsClient.Exchange(request, dnsServer)
			if err != nil {
				return nil, err
			}
			return response, nil
		}
		return nil, nil
	default:
		return nil, errors.New("wrong opcode")
	}
}
