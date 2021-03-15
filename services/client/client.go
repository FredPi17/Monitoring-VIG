package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

var mqtt_api_key = GetEnv("MQTT_API_KEY", "")
var mqtt_endpoint = GetEnv("MQTT_ENDPOINT", "http://localhost:8888")
var influx_token = GetEnv("INFLUX_TOKEN", "")
var influx_endpoint = GetEnv("INFLUX_ENDPOINT", "http://localhost:8086")
var mqtt_client_bucket = GetEnv("INFLUX_CLIENT_BUCKET", "")
var org = GetEnv("INFLUX_CLIENT_ORG", "")

func main() {
	wait := make(chan struct{})
	fmt.Println("Starting server loop")
	for {
		ClientPlay()
		go loop(wait)
		<-wait
	}
}

func ClientPlay() {
	httpRequest := &http.Client{}
	influx := influxdb2.NewClient(influx_endpoint, influx_token)
	writeAPIClient := influx.WriteAPI(org, mqtt_client_bucket)

	clientSession, err := http.NewRequest("GET", "http://"+mqtt_api_key+"@"+mqtt_endpoint+"/api/v1/session/show", nil)
	if err != nil {
		fmt.Println("error get api call")
	}
	resp, err := httpRequest.Do(clientSession)
	var response ClientResponse
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(body)), &response)

	for i := range response.Table {
		p := influxdb2.NewPointWithMeasurement("client").
			AddField("client_id", response.Table[i].ClientId).AddField("connected", response.Table[i].IsOnline).AddField("ip", response.Table[i].PeerHost).SetTime(time.Now())
		writeAPIClient.WritePoint(p)
		writeAPIClient.Flush()
	}
	// always close client at the end
	defer influx.Close()
}

///// DOCUMENTATION /////
// this struct object is the unit client object returned by the MQTT API
///// DOCUMENTATION /////

type Client struct {
	ClientId   string `json:"client_id"`
	IsOnline   bool   `json:"is_online"`
	Mountpoint string `json:"mountpoint"`
	PeerHost   string `json:"peer_host"`
	PeerPort   int    `json:"peer_port"`
	User       string `json:"user"`
}

///// DOCUMENTATION /////
// this struct object is the array client object returned by the MQTT API
///// DOCUMENTATION /////

type ClientResponse struct {
	Table []Client
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

///// DOCUMENTATION /////
// this function make waiting time like a sleeping mode
///// DOCUMENTATION /////

func loop(ch chan struct{}) {
	fmt.Println("Server API requested... waiting for a minute")
	time.Sleep(1 * time.Minute)
	ch <- struct{}{}
}
