package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/tarm/serial"
)

func main() {
	var device = flag.String("device", "/dev/ttyUSB0", "specify the device")
	var baudrate = flag.Int("baudrate", 115200, "specify the baudrate")

	var influxdb_host = flag.String("influx-host", "http://mon.mijnbaopt.nl:8086", "specify host of InfluxDB")
	var influxdb_user = flag.String("influx-user", "", "specify InfluxDB user")
	var influxdb_pass = flag.String("influx-password", "", "specify InfluxDB password")
	var influxdb_db = flag.String("influx-db", "rommel", "specify InfluxDB database")

	flag.Parse()

	serial_config := &serial.Config{
		Name: *device,
		Baud: *baudrate,
	}

	s, err := serial.OpenPort(serial_config)
	if err != nil {
		panic(fmt.Sprintf("Could not open serial port '%s': %s", device, err))
	}

	influx_client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     *influxdb_host,
		Username: *influxdb_user,
		Password: *influxdb_pass,
	})
	if err != nil {
		panic(fmt.Sprintf("Could not create InfluxDB client: %s", err))
	}
	influx_client.Close()

	_, _, err = influx_client.Ping(2)

	if err != nil {
		panic(fmt.Sprintf("Cant ping InfluxDB: %s", err))
	}

	batch_point_config := client.BatchPointsConfig{
		Database:  *influxdb_db,
		Precision: "s",
	}

	scanner := bufio.NewScanner(bufio.NewReader(s))

	r, err := regexp.Compile("(.*)\\(([^\\*]*)(?:\\*(.*))?\\)")

	if err != nil {
		panic(fmt.Sprintf("Regex didn't compile: %s", err))
	}

	var bp client.BatchPoints
	var fields map[string]interface{}

	for scanner.Scan() {
		l := scanner.Text()

		if strings.Contains(l, "KAIFA") {
			if bp != nil {
				// Write old one
				pt, _ := client.NewPoint("energy_usage", make(map[string]string), fields)
				bp.AddPoint(pt)
				err = influx_client.Write(bp)

				if err != nil {
					panic(fmt.Sprintf("Could not write points to influx: %s", err))
				}
			}

			bp, err = client.NewBatchPoints(batch_point_config)
			fields = make(map[string]interface{})

			if err != nil {
				panic(fmt.Sprintf("Could not create batch poitns: %s", err))
			}

		}
		matches := r.FindStringSubmatch(l)

		if len(matches) >= 2 {

			switch matches[1] {
			case "1-0:1.7.0":
				v, _ := strconv.ParseFloat(matches[2], 64)
				fields["electricity_actual_usage"] = v
			case "1-0:1.8.1":
				v, _ := strconv.ParseFloat(matches[2], 64)
				fields["electricity_total_usage_normal_tariff"] = v
			case "1-0:1.8.2":
				v, _ := strconv.ParseFloat(matches[2], 64)
				fields["electricity_total_usage_low_tariff"] = v
			case "0-1:24.2.1":
				v, _ := strconv.ParseFloat(matches[2], 64)
				fields["gas_total_usage"] = v
			}

		}
	}

	if err = scanner.Err(); err != nil {
		panic(fmt.Sprintf("Could not read from serial port: %s", err))
	}
}
