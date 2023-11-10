package main

import (
	"fmt"
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
