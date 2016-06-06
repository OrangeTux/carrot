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

	var influxdb_host = flag.String("influx-host", "http://localhost:8086", "specify host of InfluxDB")
	var influxdb_user = flag.String("influx-user", "", "specify InfluxDB user")
	var influxdb_pass = flag.String("influx-password", "", "specify InfluxDB password")
	var influxdb_db = flag.String("influx-db", "carrot", "specify InfluxDB database")

	flag.Parse()

	s, err := serial.OpenPort(&serial.Config{
		Name: *device,
		Baud: *baudrate,
	})

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
	defer influx_client.Close()

	_, _, err = influx_client.Ping(2)

	if err != nil {
		panic(fmt.Sprintf("Can't ping InfluxDB: %s", err))
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  *influxdb_db,
		Precision: "s",
	})

	scanner := bufio.NewScanner(bufio.NewReader(s))
	scanner.Split(SplitTelegram)

	if err != nil {
		panic(fmt.Sprintf("Could not create batch points: %s", err))
	}

	for scanner.Scan() {
		t := Telegram{}
		t.UnmarshalBinary(scanner.Bytes())

		pt, _ := client.NewPoint(
			"energy_usage",
			map[string]string{"equipment_id": strconv.Itoa(t.EquipmentId)},
			map[string]interface{}{
				"electricity_actual_usage":              t.CurrentPowerUsage,
				"electricity_total_usage_normal_tariff": t.PowerUsedTariff2,
				"electricity_total_usage_low_tariff":    t.PowerUsedTariff1,
				"gas_total_usage":                       t.GasUsed,
			})

		bp.AddPoint(pt)
		err = influx_client.Write(bp)

		if err != nil {
			fmt.Sprintf("Could not write data: %s", err)
		}
	}

	if err = scanner.Err(); err != nil {
		panic(fmt.Sprintf("Could not read from serial port: %s", err))
	}
}
