package output

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/naoina/toml"
)

type Config struct {
	C struct {
		Server   string
		Username string
		Password string
		Database string
	} `toml:"influxdb"`
}

type InfluxDB struct {
	Config Config
	Client client.Client
}

func (i *InfluxDB) Write(bp client.BatchPoints) error {
	return i.Client.Write(bp)
}

func NewOutput(b []byte) (*InfluxDB, error) {
	var conf Config

	if err := toml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.C.Server,
		Username: conf.C.Username,
		Password: conf.C.Password,
	})

	i := InfluxDB{Config: conf, Client: c}

	if err != nil {
		return nil, err
	}
	//_, _, err = i.Client.Ping(2)

	if err != nil {
		return nil, err
	}

	return &i, nil
}
