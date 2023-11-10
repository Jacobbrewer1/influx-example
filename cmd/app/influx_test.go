package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"testing"
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
			got, err := connectToInfluxDB()
			require.NoError(t, err, "connectToInfluxDB() error = %v, wantErr %v", err, tt.wantErr)

			health, err := got.Health(context.Background())
			require.Equal(t, tt.wantErr, err, "Error getting health of InfluxDB")
			require.Equal(t, domain.HealthCheckStatusPass, health.Status, "InfluxDB is not healthy")

			// List influx databases
			databases, err := got.OrganizationsAPI().GetOrganizations(context.Background())
			require.NoError(t, err, "Error getting list of databases")
			for _, db := range *databases {
				fmt.Printf("Database: %s\n", db.Name)
			}

			got.Close()
		})
	}
}

func Test_write_event_with_line_protocol(t *testing.T) {
	tests := []struct {
		name  string
		f     func(influxdb2.Client, []ThermostatSetting)
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
			f: func(c influxdb2.Client, datas []ThermostatSetting) {
				// Send all the data to the DB
				for {
					for _, data := range datas {
						for i := 0; i < 10; i++ {
							d := data
							d.avg += float64(i)
							d.max += float64(i)
							writeEventWithLineProtocol(c, d)
							writeEventWithFluentStyle(c, d)
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// helper to initialise and clean the database
			client := initTestdb(t)
			// call function under test
			tt.f(client, tt.datas)
		})
	}
}

func initTestdb(t *testing.T) influxdb2.Client {
	t.Helper() // Tells `go test` that this is a helper

	err := godotenv.Load("../../influxdb_example.env") // load environement variable
	require.NoError(t, err, "Error loading .env file")

	bucket = `user_events`
	org = `iot`

	client, err := connectToInfluxDB() // create the client
	require.NoError(t, err, "Error connecting to InfluxDB")

	// Clean the database by deleting the bucket
	ctx := context.Background()
	bucketsAPI := client.BucketsAPI()
	dBucket, err := bucketsAPI.FindBucketByName(ctx, bucket)
	if err == nil {
		err := client.BucketsAPI().DeleteBucketWithID(context.Background(), *dBucket.Id)
		require.NoError(t, err, "Error deleting bucket")
	}

	// create new empty bucket
	dOrg, _ := client.OrganizationsAPI().FindOrganizationByName(ctx, org)
	_, err = client.BucketsAPI().CreateBucketWithNameWithID(ctx, *dOrg.Id, bucket)
	require.NoError(t, err, "Error creating bucket")

	return client
}
