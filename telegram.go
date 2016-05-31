package main

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

type Telegram struct {
	EquipmentId               int     // 0-0:96.1.1
	PowerUsedLowTariff        float64 // 1-0:1.8.1
	PowerUsedNormalTariff     float64 // 1-0:1.8.2
	PowerProducedLowTariff    float64 // 1-0:2.8.1
	PowerProducedNormalTariff float64 // 1-0:2.8.2
	CurrentTariff             int     // 0-0:96.14.0
	CurrentPowerUsage         float64 // 1-0:1.7.0
	CurrentPowerProduced      float64 // 1-0:2.7.0
	GasUsed                   float64 // 0-1:24.2.1
}

func (t *Telegram) UnmarshalBinary(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	// Regex which matches the string till the first '('.
	//
	// 1-0:2.8.2(000000.000*kWh)
	// ^^^^^^^^^
	//regex_id, _ := regexp.Compile("^.*?(?=[\\(|&])")

	// Regex matching the values, excluding possible units.
	//
	// 0-0:96.7.9(00001)
	//            ^^^^^
	//
	// 0-1:24.2.1(160525200000S)(00000.866*kWh)
	//            ^^^^^^^^^^^^   ^^^^^^^^^
	regex_values, _ := regexp.Compile("(?:\\()([0-9\\.]*)")

	for scanner.Scan() {
		l := scanner.Text()

		i := strings.Index(l, "(")
		if i == -1 {
			continue
		}

		id := l[0:i]

		values := regex_values.FindAllStringSubmatch(l, -1)
		if len(values) == 0 {
			continue
		}

		value := values[len(values)-1][1]

		switch id {
		case "0-0:96.1.1":
			v, _ := strconv.Atoi(value)
			t.EquipmentId = v
		case "1-0:1.7.0":
			v, _ := strconv.ParseFloat(value, 64)
			t.CurrentPowerUsage = v
		case "1-0:1.8.1":
			v, _ := strconv.ParseFloat(value, 64)
			t.PowerUsedLowTariff = v
		case "1-0:1.8.2":
			v, _ := strconv.ParseFloat(value, 64)
			t.PowerUsedNormalTariff = v
		case "0-1:24.2.1":
			v, _ := strconv.ParseFloat(value, 64)
			t.GasUsed = v
		}
	}

	return nil
}

func SplitTelegram(data []byte, atEOF bool) (advanced int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte{'!'}); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
