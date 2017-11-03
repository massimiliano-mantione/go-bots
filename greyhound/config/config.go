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
	MinOutMillis  int
	KP            int
	KP2           int
	KD            int
	KD2           int
	MaxSteering   int
	MaxPos        int
	MaxPos2       int
	MaxPosD       int
	MaxPosD2      int
}

// CompleteConfig fills in computed configutation fields
func CompleteConfig(c *Config) {
	c.MaxSteering = (c.MaxSpeed * c.MaxSteeringPC) / 100
	c.MaxPos = c.SensorRadius * 3
	c.MaxPos2 = c.MaxPos * c.MaxPos
	c.MaxPosD = c.MaxPos / c.MinOutMillis
	c.MaxPosD2 = c.MaxPos * c.MaxPos
}

// Default Config data
func Default() Config {
	result := Config{
		// MaxSpeed:  100,
		MaxSpeed:      30,
		MaxSteeringPC: 120,
		SensorRadius:  100,
		SensorSpan:    700,
		SensorMin:     80,
		MinDTicks:     10,
		MaxDTicks:     30000,
		MinOutMillis:  10,
		KP:            70,
		KP2:           10,
		KD:            0,
		KD2:           0,
	}
	CompleteConfig(&result)
	return result
}

// FromString reads Config data from a TOML string
func FromString(data string) (Config, error) {
	result := Config{}
	_, err := toml.Decode(data, &result)
	if err != nil {
		return result, err
	}
	CompleteConfig(&result)
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
