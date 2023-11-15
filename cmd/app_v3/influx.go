package main

import (
	"errors"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"os"
)

// Connect to an Influx Database reading the credentials from
// environment variables INFLUXDB_TOKEN, INFLUXDB_URL
// return influxdb Client or errors
func connectToInfluxDB() (*influxdb3.Client, error) {
	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	database := "user_events"

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     dbURL,
		Token:    dbToken,
		Database: database,
	})

	return client, err
}
