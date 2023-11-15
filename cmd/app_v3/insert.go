package main

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"time"
)

type ThermostatSetting struct {
	user string
	max  float64 //temperature
	avg  float64 //temperature
}

func writeEventWithFluentStyle(client *influxdb3.Client, t ThermostatSetting) error {
	point := influxdb3.NewPointWithMeasurement("thermostat").
		SetTag("unit", "temperature").
		SetTag("user", t.user).
		SetField("avg", t.avg).
		SetField("max", t.max).
		SetTimestamp(time.Now())

	// write point asynchronously
	err := client.WritePoints(context.Background(), point)
	if err != nil {
		return fmt.Errorf("error saving point to InfluxDB: %v", err)
	}
	return nil
}
