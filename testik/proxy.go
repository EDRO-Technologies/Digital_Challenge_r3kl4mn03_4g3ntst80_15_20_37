package main

import (
	"errors"

	"github.com/miekg/dns"
)

func processRequest(dnsClient *dns.Client, request *dns.Msg, config Config) (*dns.Msg, error) {
	switch request.Opcode {
	case dns.OpcodeQuery:
		if len(request.Question) > 0 {
			dnsServer := config.DNSServer
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
