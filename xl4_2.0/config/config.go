package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config data
type Config struct {
	AccelPerTicks        int
	MaxSpeed             int
	TrackTurnSpeed       int
	SeekTurnSpeed        int
	TrackSpeed           int
	MaxIrValue           int
	StrategyR1Time       int
	StrategyS1Time       int
	StrategyR2Time       int
	StrategyS2Time       int
	StrategyStraightTime int
}

// Default Config data
func Default() Config {
	result := Config{
		AccelPerTicks:        40,
		MaxSpeed:             100,
		TrackTurnSpeed:       40,
		SeekTurnSpeed:        40,
		TrackSpeed:           100,
		MaxIrValue:           30,
		StrategyR1Time:       220,
		StrategyS1Time:       120,
		StrategyR2Time:       470,
		StrategyS2Time:       120,
		StrategyStraightTime: 400,
	}
	fixConfig(&result)
	return result
}

// FromString reads Config data from a TOML string
func FromString(data string) (Config, error) {
	result := Config{}
	_, err := toml.Decode(data, &result)
	if err != nil {
		return result, err
	}
	fixConfig(&result)
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

func fixConfig(c *Config) {
	c.StrategyR1Time *= 1000
	c.StrategyS1Time *= 1000
	c.StrategyR2Time *= 1000
	c.StrategyS2Time *= 1000
	c.StrategyStraightTime *= 1000
}
