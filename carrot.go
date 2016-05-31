package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/tarm/serial"
)

func main() {
	var device = flag.String("device", "/dev/ttyUSB0", "specify the device")
	var baudrate = flag.Int("baudrate", 115200, "specify the baudrate")

	var influxdb_host = flag.String("influx-host", "http://localhost.nl:8086", "specify host of InfluxDB")
	var influxdb_user = flag.String("influx-user", "", "specify InfluxDB user")
	var influxdb_pass = flag.String("influx-password", "", "specify InfluxDB password")
	var influxdb_db = flag.String("influx-db", "carrot", "specify InfluxDB database")

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
	scanner.Split(SplitTelegram)

	bp, err := client.NewBatchPoints(batch_point_config)
	if err != nil {
		panic(fmt.Sprintf("Could not create batch poitns: %s", err))
	}

	for scanner.Scan() {
		t := Telegram{}
		t.UnmarshalBinary(scanner.Bytes())

		pt, _ := client.NewPoint(
			"energy_usage",
			map[string]string{"equiment_id": strconv.Itoa(t.EquipmentId)},
			map[string]interface{}{
				"electricity_actual_usage":              t.CurrentPowerUsage,
				"electricity_total_usage_normal_tariff": t.PowerUsedNormalTariff,
				"electricity_total_usage_low_tariff":    t.PowerUsedLowTariff,
				"gas_total_usage":                       t.GasUsed,
			})

		bp.AddPoint(pt)
		err = influx_client.Write(bp)
	}

	if err = scanner.Err(); err != nil {
		panic(fmt.Sprintf("Could not read from serial port: %s", err))
	}
}
