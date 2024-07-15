# Speedtest and InfluxDB Integration

This project performs periodic internet speed tests using the `speedtest-go` library and stores the results in an InfluxDB instance. The results include latency, download speed, and upload speed, which are recorded at regular intervals.

## Prerequisites

- Go
- Docker
- InfluxDB running in Docker

## Setup

### 1. Install Go Packages

Run the following command to install the necessary Go packages:

```sh
go get github.com/showwin/speedtest-go
go get github.com/influxdata/influxdb-client-go/v2
go get github.com/joho/godotenv
```

### 2. Create `.env` File

```sh
INFLUXDB_TOKEN=your-influxdb-token
INFLUXDB_ORG=your-org
INFLUXDB_BUCKET=test
```

### Run InfluxDB with Docker

Pull first:

```sh
docker pull influxdb:latest
```

```sh
docker run -d -p 8086:8086 --name influxdb -v influxdb:/var/lib/influxdb influxdb:latest
```

Go to the `localhost:8086` to configure InfluxDB (your password, token etc.)

### Run the application

```
go run main.go
```
