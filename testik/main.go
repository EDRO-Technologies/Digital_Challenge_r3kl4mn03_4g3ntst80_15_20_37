package main

import (
	"log"
	"time"

	"github.com/miekg/dns"
)

func main() {
	config, err := GetConfig()

	if err != nil {
		panic(err)
	}

	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"

	var dnsCache Cache
	if config.UseCache {
		dnsCache = InitCache(config.CacheExpiration)
	}

	dns.HandleFunc(".", func(writer dns.ResponseWriter, request *dns.Msg) {
		switch request.Opcode {
		case dns.OpcodeQuery:
			startTime := time.Now()

			response, err := ProcessRequest(dnsClient, &dnsCache, request, config)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", request, err.Error())
			}

			duration := time.Since(startTime)

			logData := GetRequestInfo(request, response, writer.RemoteAddr().String(), startTime, duration)

			go func() {
				err := SendToLogstash("localhost:50000", logData)
				if err != nil {
					log.Printf("Failed to send log to Logstash: %v\n", err)
				}
			}()

			response.SetReply(request)
			writer.WriteMsg(response)
		}
	})

	server := &dns.Server{Addr: config.Host, Net: "udp"}
	log.Printf("Starting at %s\n", config.Host)
	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("Failed to start server: %s\n ", err.Error())
	}
}
