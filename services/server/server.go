package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

///// DOCUMENTATION /////
// the main() function connects to MQTT broker to get health status
// then send the value to the influxdb instance
///// DOCUMENTATION /////
var influx_token = GetEnv("INFLUX_TOKEN", "")
var influx_server = GetEnv("INFLUX_ENDPOINT", "http://localhost:8086")
var influx_bucket = GetEnv("INFLUX_SERVER_BUCKET", "")
var influx_org = GetEnv("INFLUX_SERVER_ORG", "")
var mqtt_server = GetEnv("MQTT_ENDPOINT", "http://localhost:8888")

func main() {
	wait := make(chan struct{})
	fmt.Println("Starting server loop")
	for {
		ServerPlay()
		go loop(wait)
		<-wait
	}
}

func ServerPlay() {
	httpRequest := &http.Client{}
	influx := influxdb2.NewClient(influx_server, influx_token)
	writeAPIServer := influx.WriteAPI(influx_org, influx_bucket)
	serverHealth, err := http.NewRequest("GET", "http://"+mqtt_server+"/health", nil)
	if err != nil {
		fmt.Println("error get api call")
	}
	if influx_token == "" || influx_bucket == "" || influx_org == "" {
		fmt.Errorf("INFLUX_TOKEN, INFLUX_BUCKET and INFLUX_ORG are mandatory ! error : %w", err)
	} else {
		resp, err := httpRequest.Do(serverHealth)
		if err != nil {
			fmt.Errorf("Error while calling MQTT Server at : "+mqtt_server+" with error: %v", err)
		}
		if resp.StatusCode == 200 {
			p := influxdb2.NewPointWithMeasurement("server").AddField("status", true).SetTime(time.Now())
			writeAPIServer.WritePoint(p)
			writeAPIServer.Flush()
		} else {
			p := influxdb2.NewPointWithMeasurement("server").AddField("status", false).SetTime(time.Now())
			writeAPIServer.WritePoint(p)
			writeAPIServer.Flush()
		}
		defer influx.Close()
	}
}

///// DOCUMENTATION /////
// this function gets environment variable with a key and a default value
// the key is the environment variable to get
// the defaultValue is the value returned if the environment variable length is equal to 0
///// DOCUMENTATION /////

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func loop(ch chan struct{}) {
	fmt.Println("Server API requested... waiting for a minute")
	time.Sleep(1 * time.Minute)
	ch <- struct{}{}
}
