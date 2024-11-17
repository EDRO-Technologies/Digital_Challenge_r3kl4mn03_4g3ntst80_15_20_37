package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(proxy *Proxy) {
	if !proxy.config.SQLiteEnabled {
		return
	}
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS dns_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		client_ip TEXT NOT NULL,
		dest_ip TEXT NOT NULL,
		query_name TEXT NOT NULL,
		query_type TEXT NOT NULL,
		protocol TEXT NOT NULL,
		response_code TEXT NOT NULL,
		response_time_ms INTEGER NOT NULL,
		answers TEXT NOT NULL
	);`
	_, err := proxy.db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

}

func InsertLog(proxy *Proxy, logData map[string]interface{}) error {
	if !proxy.config.SQLiteEnabled {
		return nil
	}

	answersJSON, err := json.Marshal(logData["answers"])
	if err != nil {
		return fmt.Errorf("failed to serialize answers: %w", err)
	}

	insertSQL := `
	INSERT INTO dns_logs (
		timestamp, client_ip, dest_ip, query_name, query_type, protocol, response_code, response_time_ms, answers
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = proxy.db.Exec(insertSQL,
		logData["timestamp"],
		logData["client_ip"],
		logData["dest_ip"],
		logData["query_name"],
		logData["query_type"],
		logData["protocol"],
		logData["response_code"],
		logData["response_time_ms"],
		answersJSON,
	)
	return err
}

func GetFirstLog(proxy *Proxy) (map[string]interface{}, error) {
	if !proxy.config.SQLiteEnabled {
		return nil, nil
	}
	querySQL := `SELECT id, timestamp, client_ip, dest_ip, query_name, query_type, protocol, response_code, response_time_ms, answers FROM dns_logs ORDER BY id LIMIT 1`
	row := proxy.db.QueryRow(querySQL)

	var id int
	var timestamp, clientIP, destIP, queryName, queryType, protocol, responseCode, answersJSON string
	var responseTimeMs int64

	err := row.Scan(&id, &timestamp, &clientIP, &destIP, &queryName, &queryType, &protocol, &responseCode, &responseTimeMs, &answersJSON)
	if err != nil {
		return nil, err
	}

	// Deserialize the answers JSON back into an array of strings
	var answers []string
	err = json.Unmarshal([]byte(answersJSON), &answers)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize answers: %w", err)
	}

	logData := map[string]interface{}{
		"id":               id,
		"timestamp":        timestamp,
		"client_ip":        clientIP,
		"dest_ip":          destIP,
		"query_name":       queryName,
		"query_type":       queryType,
		"protocol":         protocol,
		"response_code":    responseCode,
		"response_time_ms": responseTimeMs,
		"answers":          answers,
	}
	return logData, nil
}

func DeleteFirstLog(proxy *Proxy) error {
	if !proxy.config.SQLiteEnabled {
		return nil
	}
	deleteSQL := `DELETE FROM dns_logs WHERE id = (SELECT id FROM dns_logs ORDER BY id LIMIT 1)`
	_, err := proxy.db.Exec(deleteSQL)
	return err
}
