package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

func GetRequestInfo(request *dns.Msg, response *dns.Msg, clientAddr string, startTime time.Time, duration time.Duration) map[string]interface{} {
	clientIP, _, err := net.SplitHostPort(clientAddr)
	if err != nil {
		log.Printf("Failed to parse client address: %s\n", clientAddr)
		clientIP = clientAddr
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

	return logData
}

func SendToLogstash(address string, data map[string]interface{}) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("error connecting to Logstash: %w", err)
	}
	defer conn.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error sending data to Logstash: %w", err)
	}

	_, err = conn.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("error writing newline: %w", err)
	}

	return nil
}
