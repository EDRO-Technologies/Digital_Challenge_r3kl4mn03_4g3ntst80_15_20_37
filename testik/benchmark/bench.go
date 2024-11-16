package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func main() {
	N := 1000
	start := time.Now()
	dnsClient := new(dns.Client)
	msg := new(dns.Msg)
	dnsClient.Net = "udp"
	msg.Question = append(msg.Question, dns.Question{
		Name:  "www.google.com.",
		Qtype: 1,
	})
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := dnsClient.Exchange(msg, "192.168.1.114:53")
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Made %d requests, it took %f seconds\n", N, time.Since(start).Seconds())
}
