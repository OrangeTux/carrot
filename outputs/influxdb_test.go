package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tomlConfig = `
[influxdb]
server = "http://localhost:8086"
username = "orangetux"
password = "something_secret"
database = "carrot"
`

func TestNewOutput(t *testing.T) {
	i, err := NewOutput([]byte(tomlConfig))

	if err != nil {
		panic(err)
	}

	assert.Equal(t, i.Config.C.Server, "http://localhost:8086", "should be equal")
}
