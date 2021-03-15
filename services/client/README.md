# MQTT CLIENT MONITORING

This application is used to monitor an MQTT server running with [VerneMQ](https://vernemq.com/). 

The usage is to get metrics about clients on the server, and send it on an [Influxdb](https://docs.influxdata.com/) database to be shown on a [Grafana](https://grafana.com/) dashboard. 

## HOW TO USE IT 

Some informations are required to run this application. 

### STACK SIDE

You need to have : 

* an Influxdb V2 instance with a `token`, an `org` and a `bucket` created on it
* an VerneMQ MQTT server with a `token` generated with the vmq-admin CLI

### APPLCIATION SIDE

You need to give informations to this application like environments variables :  

* `INFLUX_TOKEN` : the token you've created on influxdb - (mandatory)
* `INFLUX_ENDPOINT` : the endpoint address of the influxdb instance - (optional | default to `http://localhost:8086`)
* `INFLUX_CLIENT_BUCKET` : the bucket you want to send data to - (mandatory)
* `INFLUX_CLIENT_ORG` : the organization in relation with the bucket - (mandatory)
* `MQTT_ENDPOINT` : the MQTT server endpoint - (optional | default to `http://localhost:8888`)
* `MQTT_API_KEY` : the MQTT api key to query the MQTT API - (mandatory)

### EXAMPLE OF USE WITH `docker-compose`

```yaml
version: "2"
    services:
    mqtt-client-monit:
        container_name: mqtt-client-monit
        image: fredericpinaud/mqtt-client-monit:latest
        environment: 
        - INFLUX_TOKEN=<YOU_TOKEN>
        - INFLUX_ENDPOINT=<INFLUXDB ENDPOINT>
        - INFLUX_CLIENT_BUCKET=<INFLUXDB CLIENT BUCKET>
        - INFLUX_CLIENT_ORG=<INFLUXDB CLIENT ORGANIZATION>
        - MQTT_ENDPOINT=<MQTT ENDPOINT ADDRESS>
        - MQTT_API_KEY=<MQTT API KEY>
```