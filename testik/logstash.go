package main

import (
	"encoding/json"
	"fmt"
	"net"
)

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
