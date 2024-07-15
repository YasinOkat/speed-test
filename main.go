package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
	"github.com/showwin/speedtest-go/speedtest"
)

const bytesToMegabytes = 1000000.0

func runSpeedTest(speedtestClient *speedtest.Speedtest, writeAPI api.WriteAPIBlocking) {
	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		log.Printf("Failed to fetch servers: %v\n", err)
		return
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		log.Printf("Failed to find server: %v\n", err)
		return
	}

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()

		downloadSpeed := s.DLSpeed.Mbps()
		uploadSpeed := s.ULSpeed.Mbps()

		fmt.Printf("Latency: %s, Download: %.2f Mbps, Upload: %.2f Mbps\n", s.Latency, downloadSpeed, uploadSpeed)
		writeResultToInfluxDB(writeAPI, s.Latency.Seconds(), downloadSpeed, uploadSpeed)
		s.Context.Reset()
	}
}

func writeResultToInfluxDB(writeAPI api.WriteAPIBlocking, latency float64, download float64, upload float64) {
	p := influxdb2.NewPointWithMeasurement("speedtest").
		AddField("latency", latency).
		AddField("download", download).
		AddField("upload", upload).
		SetTime(time.Now())
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Printf("Failed to write point: %v\n", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	influxDBToken := os.Getenv("INFLUX_TOKEN")
	influxDBOrg := os.Getenv("INFLUX_ORG")
	influxDBBucket := os.Getenv("INFLUX_BUCKET")

	influxDB := influxdb2.NewClient("http://localhost:8086", influxDBToken)
	defer influxDB.Close()
	writeAPI := influxDB.WriteAPIBlocking(influxDBOrg, influxDBBucket)

	speedtestClient := speedtest.New()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runSpeedTest(speedtestClient, writeAPI)
		}
	}
}
