package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config data
type Config struct {
	MaxSpeed      int
	MaxSteeringPC int
	SensorRadius  int
	SensorSpan    int
	SensorMin     int
	MinDTicks     int
	MaxDTicks     int
	DTicksBoost   int
	KP            int
	KP2           int
	KPR           int
	KD            int
	KD2           int
	KDR           int
	KR            int
}

// Default Config data
func Default() Config {
	return Config{
		// MaxSpeed:  1000000,
		MaxSpeed:      300000,
		MaxSteeringPC: 120,
		SensorRadius:  500,
		SensorSpan:    700,
		SensorMin:     80,
		MinDTicks:     10,
		MaxDTicks:     30000,
		DTicksBoost:   100000,
		KP:            1800,
		KP2:           1,
		KPR:           4,
		KD:            0,
		KD2:           0,
		KDR:           1,
		KR:            1,
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
