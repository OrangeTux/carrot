package main

import (
	"bufio"
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Telegram struct {
	EquipmentId               int     `dsmr:"0-0:96.1.1"`
	PowerUsedLowTariff        float64 `dsmr:"1-0:1.8.1"`
	PowerUsedNormalTariff     float64 `dsmr:"1-0:1.8.2"`
	PowerProducedLowTariff    float64 `dsmr:"1-0:2.8.1"`
	PowerProducedNormalTariff float64 `dsmr:"1-0:2.8.2"`
	CurrentTariff             int     `dsmr:"0-0:96.14.0"`
	CurrentPowerUsage         float64 `dsmr:"1-0:1.7.0"`
	CurrentPowerProduced      float64 `dsmr:"1-0:2.7.0"`
	GasUsed                   float64 `dsmr:"0-1:24.2.1"`
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

		type_value := reflect.ValueOf(*t)

		// Iterate over Telegram's fields to find the field which tag
		// matches with the id from the telegram. Then the value
		// extracted from the telegraf is assigned to that field.
		for i := 0; i < type_value.NumField(); i++ {

			if type_value.Type().Field(i).Tag.Get("dsmr") == id {
				field := reflect.ValueOf(t).Elem().Field(i)

				if field.Kind() == reflect.Int {
					v, _ := strconv.ParseInt(value, 10, 64)
					reflect.ValueOf(t).Elem().Field(i).SetInt(v)
				}

				if field.Kind() == reflect.Float64 {
					v, _ := strconv.ParseFloat(value, 64)
					reflect.ValueOf(t).Elem().Field(i).SetFloat(v)
				}

			}
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
