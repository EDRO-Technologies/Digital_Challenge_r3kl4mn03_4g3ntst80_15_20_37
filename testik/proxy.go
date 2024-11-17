package main

import (
	"errors"

	"github.com/miekg/dns"
)

func GetCache(dnsCache *Cache, request *dns.Msg, config Config) *dns.Msg {
	if !config.UseCache {
		return nil
	}
	question := request.Question[0]
	if question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA {
		cached, found := dnsCache.Get(question.Name)
		if found {
			rr := cached.(*dns.RR)
			response := new(dns.Msg)
			response.Answer = append(response.Answer, *rr)
			return response
		}
	}
	return nil
}

func SetCache(dnsCache *Cache, response *dns.Msg, question dns.Question, config Config) {
	if !config.UseCache {
		return
	}
	if len(response.Answer) > 0 {
		switch response.Answer[len(response.Answer)-1].(type) {
		case *dns.A, *dns.AAAA:
			dnsCache.Set(
				question.Name,
				&response.Answer[len(response.Answer)-1],
			)
		default:
			return
		}
	}
}

func ProcessRequest(proxy *Proxy, request *dns.Msg, config Config) (*dns.Msg, error) {
	switch request.Opcode {
	case dns.OpcodeQuery:
		if len(request.Question) > 0 {

			response := GetCache(&proxy.cache, request, config)
			if response != nil {
				return response, nil
			}

			response, _, err := proxy.client.Exchange(request, config.DNSServer)
			if err != nil {
				return nil, err
			}

			SetCache(&proxy.cache, response, request.Question[0], config)

			return response, nil
		}
		return nil, nil
	default:
		return nil, errors.New("wrong opcode")
	}
}
