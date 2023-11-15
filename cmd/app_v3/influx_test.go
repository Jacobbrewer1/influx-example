package main

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func Test_connectToInfluxDB(t *testing.T) {
	// load environment variable from a file for test purposes
	err := godotenv.Load("../../influxdb_example.env")
	require.NoError(t, err, "Error loading .env file")

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Successful connection to InfluxDB",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := connectToInfluxDB()
			require.NoError(t, err, "connectToInfluxDB() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func Test_write_event_with_line_protocol(t *testing.T) {
	// load environment variable from a file for test purposes
	err := godotenv.Load("../../influxdb_example.env")
	require.NoError(t, err, "Error loading .env file")

	tests := []struct {
		name  string
		f     func(*influxdb3.Client, []ThermostatSetting)
		datas []ThermostatSetting
	}{
		{
			name: "Write new record with line protocol",
			// Your data Points
			datas: []ThermostatSetting{
				{
					user: "foo",
					avg:  35.5,
					max:  42,
				},
			},
			f: func(c *influxdb3.Client, datas []ThermostatSetting) {
				// Send all the data to the DB
				for {
					for _, data := range datas {
						for i := 0; i < 10; i++ {
							d := data
							d.avg += float64(i)
							d.max += float64(i)
							err := writeEventWithFluentStyle(c, d)
							require.NoError(t, err, "Error writing event with line protocol")
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// helper to initialise and clean the database
			client, err := connectToInfluxDB()
			require.NoError(t, err, "Error connecting to InfluxDB")

			// call function under test
			tt.f(client, tt.datas)
		})
	}
}
