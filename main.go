package main

import (
	"log"

	"github.com/miekg/dns"
)

func main() {
	host := "10.0.2.15:53"

	dns.HandleFunc(".", func(writer dns.ResponseWriter, request *dns.Msg) {
		switch request.Opcode {
		case dns.OpcodeQuery:
			response, err := processRequest(request)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", request, err.Error())
				response.SetReply(request)
				writer.WriteMsg(response)
				return
			}
			response.SetReply(request)
			writer.WriteMsg(response)
		}
	})

	server := &dns.Server{Addr: host, Net: "udp"}
	log.Printf("Starting at %s\n", host)
	err := server.ListenAndServe()
	if err != nil {
		log.Panicf("Failed to start server: %s\n ", err.Error())
	}
}
