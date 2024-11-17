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
	CacheExpiration int64
	DNSServer       string
	LogstashAddr    string
	SQLiteEnabled   bool
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

	cacheExpirationString, cacheEnabled := os.LookupEnv("CacheExpiration")
	var cacheExpiration int64
	var err error
	if cacheEnabled && cacheExpirationString != "" {
		cacheExpiration, err = strconv.ParseInt(cacheExpirationString, 10, 64)
	} else {
		cacheEnabled = false
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

	SQLiteEnabledString, SQLiteEnabled := os.LookupEnv("SQLiteEnabled")
	if SQLiteEnabledString == "" && SQLiteEnabled {
		SQLiteEnabled = false
	}

	return Config{
		Host:            host,
		UseCache:        cacheEnabled,
		CacheExpiration: cacheExpiration * 1000000000,
		DNSServer:       dnsServer,
		LogstashAddr:    logstashAddr,
		SQLiteEnabled:   SQLiteEnabled,
	}, nil
}
