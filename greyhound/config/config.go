package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config data
type Config struct {
	MaxSpeed         int
	ParamP1          int
	ParamP2          int
	ParamD1          int
	ParamD2          int
	ParamR           int
	InnerBrakeFactor int
}

// Default Config data
func Default() Config {
	return Config{
		MaxSpeed:         1000000,
		ParamP1:          0,
		ParamP2:          0,
		ParamD1:          0,
		ParamD2:          0,
		ParamR:           0,
		InnerBrakeFactor: 0,
	}
}

// FromString reads Config data from a TOML string
func FromString(data string) (Config, error) {
	result := Config{}
	_, err := toml.Decode(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// FromFile reads Config data from a TOML file
func FromFile(fileName string) (Config, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	return FromString(string(b))
}
