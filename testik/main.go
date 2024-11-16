package main

import (
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

func main() {
	host := "192.168.1.114:53"

	dns.HandleFunc(".", func(writer dns.ResponseWriter, request *dns.Msg) {
		switch request.Opcode {
		case dns.OpcodeQuery:
			startTime := time.Now()

			response, err := processRequest(request)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", request, err.Error())
			}

			duration := time.Since(startTime)

			clientAddr := writer.RemoteAddr().String()
			clientIP, _, err := net.SplitHostPort(clientAddr)
			if err != nil {
				log.Printf("Failed to parse client address: %s\n", clientAddr)
				clientIP = clientAddr // Fallback to the original address
			}

			logData := map[string]interface{}{
				"timestamp":        startTime.Format(time.RFC3339),
				"client_ip":        clientIP,
				"query_name":       request.Question[0].Name,
				"query_type":       dns.TypeToString[request.Question[0].Qtype],
				"protocol":         "udp",
				"response_code":    dns.RcodeToString[response.Rcode],
				"response_time_ms": duration.Milliseconds(),
			}

			if len(response.Answer) > 0 {
				var answers []string
				for _, ans := range response.Answer {
					answers = append(answers, ans.String())
				}
				logData["answers"] = answers
			}

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

	server := &dns.Server{Addr: host, Net: "udp"}
	log.Printf("Starting at %s\n", host)
	err := server.ListenAndServe()
	if err != nil {
		log.Panicf("Failed to start server: %s\n ", err.Error())
	}
}
