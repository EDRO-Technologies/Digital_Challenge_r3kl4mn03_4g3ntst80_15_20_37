package main

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	UseCache bool
	CacheExpiration int
	DNSServer string
}

func GetConfig() (Config, error) {
	cacheExpirationString, present := os.LookupEnv("CacheExpiration")
	cacheExpiration, err := strconv.Atoi(cacheExpirationString)
	if err != nil {
		return Config{}, err
	}

	dnsServer := os.Getenv("DNSSERVER")
	if dnsServer == "" {
		dnsServer = "1.1.1.1:53"
	} else if strings.Contains(dnsServer, ":") {
		dnsServer += ":53"
	}
	

	return Config{
		UseCache: present,
		CacheExpiration: cacheExpiration,
		DNSServer: dnsServer,
	}, nil
}
