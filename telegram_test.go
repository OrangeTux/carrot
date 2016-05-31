package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var message = `/KFM5KAIFA-METER

1-3:0.2.8(42)
0-0:1.0.0(160525205154S)
0-0:96.1.1(1234567890)
1-0:1.8.1(000001.117*kWh)
1-0:1.8.2(000004.491*kWh)
1-0:2.8.1(000000.000*kWh)
1-0:2.8.2(000000.000*kWh)
0-0:96.14.0(0002)
1-0:1.7.0(00.563*kW)
1-0:2.7.0(00.000*kW)
0-0:96.7.21(00001)
0-0:96.7.9(00001)
1-0:99.97.0(1)(0-0:96.7.19)(000101000001W)(2147483647*s)
1-0:32.32.0(00000)
1-0:32.36.0(00000)
0-0:96.13.1()
0-0:96.13.0()
1-0:31.7.0(004*A)
1-0:21.7.0(00.563*kW)
1-0:22.7.0(00.000*kW)
0-1:24.1.0(003)
0-1:96.1.0(4730303139333430323839323633363136)
0-1:24.2.1(160525200000S)(00000.866*m3)
!D57B`

func TestTelegram(t *testing.T) {
	telegram := &Telegram{}

	telegram.UnmarshalBinary([]byte(message))

	assert.Equal(t, telegram.EquipmentId, 1234567890, "")
	assert.Equal(t, telegram.PowerUsedLowTariff, 1.117, "")
	assert.Equal(t, telegram.PowerUsedNormalTariff, 4.491, "")
	assert.Equal(t, telegram.CurrentPowerUsage, 0.563, "")
	assert.Equal(t, telegram.GasUsed, 0.866, "")

}

func TestSplitTelegram(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("a!b!c"))
	scanner.Split(SplitTelegram)

	res := []string{}
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	assert.Equal(t, res, []string{"a", "b", "c"}, "")
}
