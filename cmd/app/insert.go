package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type ThermostatSetting struct {
	user string
	max  float64 //temperature
	avg  float64 //temperature
}

func writeEventWithLineProtocol(client influxdb2.Client, t ThermostatSetting) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	// write line protocol
	writeAPI.WriteRecord(fmt.Sprintf("thermostat,unit=temperature,user=%s avg=%f,max=%f", t.user, t.avg, t.max))
	// Flush writes
	writeAPI.Flush()
}

func writeEventWithFluentStyle(client influxdb2.Client, t ThermostatSetting) {
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPI(org, bucket)

	// create point using fluent style
	p := influxdb2.NewPointWithMeasurement("thermostat").
		AddTag("unit", "temperature").
		AddTag("user", t.user).
		AddField("avg", t.avg).
		AddField("max", t.max).
		SetTime(time.Now())

	writeAPI.WritePoint(p)

	// Flush writes
	writeAPI.Flush()
}
