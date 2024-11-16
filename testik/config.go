package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host            string
	UseCache        bool
	CacheExpiration int
	DNSServer       string
	LogstashAddr    string
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetConfig() (Config, error) {
	host := os.Getenv("Host")
	if host == "" {
		host = GetOutboundIP().String() + ":53"
	}

	cacheExpirationString, present := os.LookupEnv("CacheExpiration")
	var cacheExpiration int
	var err error
	if present {
		cacheExpiration, err = strconv.Atoi(cacheExpirationString)
	}
	if err != nil {
		return Config{}, err
	}

	dnsServer := os.Getenv("DNSServer")
	if dnsServer == "" {
		dnsServer = "1.1.1.1:53"
	} else if !strings.Contains(dnsServer, ":") {
		dnsServer += ":53"
	}

	logstashAddr := os.Getenv("LogstashAddr")
	if logstashAddr == "" {
		logstashAddr = "localhost:50000"
	}

	return Config{
		Host:            host,
		UseCache:        present,
		CacheExpiration: cacheExpiration,
		DNSServer:       dnsServer,
		LogstashAddr:    logstashAddr,
	}, nil
}
