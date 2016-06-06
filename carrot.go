package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/orangetux/carrot/outputs"
	"github.com/tarm/serial"
)

func main() {
	var device = flag.String("device", "/dev/ttyUSB0", "specify the device")
	var baudrate = flag.Int("baudrate", 115200, "specify the baudrate")

	var config_file = flag.String("conf", "carrot.conf", "path to config")

	flag.Parse()

	f, err := os.Open(*config_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	i, err := output.NewOutput(b)

	s, err := serial.OpenPort(&serial.Config{
		Name: *device,
		Baud: *baudrate,
	})

	if err != nil {
		panic(fmt.Sprintf("Could not open serial port '%s': %s", device, err))
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.Config.C.Database,
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
		i.Write(bp)

		if err != nil {
			fmt.Sprintf("Could not write data: %s", err)
		}
	}

	if err = scanner.Err(); err != nil {
		panic(fmt.Sprintf("Could not read from serial port: %s", err))
	}
}
