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

	var answers []string
	var destIP string
	if len(response.Answer) > 0 {
		for _, ans := range response.Answer {
			answers = append(answers, ans.String())
			switch v := ans.(type) {
			case *dns.A:
				destIP = v.A.String()
			case *dns.AAAA:
				destIP = v.AAAA.String()
			}
		}
	}

	logData := map[string]interface{}{
		"timestamp":        startTime.Format(time.RFC3339),
		"client_ip":        clientIP,
		"dest_ip":          destIP,
		"query_name":       request.Question[0].Name,
		"query_type":       dns.TypeToString[request.Question[0].Qtype],
		"protocol":         "udp",
		"response_code":    dns.RcodeToString[response.Rcode],
		"response_time_ms": duration.Milliseconds(),
		"answers":          answers,
	}

	return logData
}

func SendToLogstash(proxy *Proxy, data map[string]interface{}) error {
	conn, err := net.Dial("tcp", proxy.config.LogstashAddr)
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

func LateSend(proxy *Proxy) {
	if !proxy.config.SQLiteEnabled {
		return
	}
	data, err := GetFirstLog(proxy)
	if err != nil {
		fmt.Printf("Error when getting FirstLog: %s\n", err)
		return
	}
	fmt.Printf("SendLATA")
	for len(data) != 0 {
		SendToLogstash(proxy, data)
		DeleteFirstLog(proxy)
		data, _ = GetFirstLog(proxy)
	}
}
