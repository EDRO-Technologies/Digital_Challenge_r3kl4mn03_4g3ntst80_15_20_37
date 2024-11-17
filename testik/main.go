package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/miekg/dns"
)

type Proxy struct {
	config Config
	cache  *Cache
	client dns.Client
	db     *sql.DB
}

func main() {
	proxy := Proxy{}

	config, err := GetConfig()
	proxy.config = config

	if err != nil {
		panic(err)
	}

	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"

	proxy.client = *dnsClient

	var dnsCache Cache
	if config.UseCache {
		dnsCache = InitCache(config.CacheExpiration)
	}
	proxy.cache = &dnsCache

	var db *sql.DB
	if config.SQLiteEnabled {
		db, err = sql.Open("sqlite3", "./logdata.db")
		if err != nil {
			panic(err)
		}
		proxy.db = db
		InitDB(&proxy)
	}
	defer db.Close()

	dns.HandleFunc(".", func(writer dns.ResponseWriter, request *dns.Msg) {
		switch request.Opcode {
		case dns.OpcodeQuery:
			startTime := time.Now()

			response, err := ProcessRequest(&proxy, request, config)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", request, err.Error())
			}

			duration := time.Since(startTime)

			logData := GetRequestInfo(request, response, writer.RemoteAddr().String(), startTime, duration)

			go func() {
				err := SendToLogstash(&proxy, logData)
				if err != nil {
					log.Printf("Failed to send log to Logstash: %s\n", err)
					err = InsertLog(&proxy, logData)
					if err != nil {
						log.Printf("Failed to add log to SQLite: %s\n", err)
					}
				} else {
					LateSend(&proxy)
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
