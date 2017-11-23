package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config data
type Config struct {
	MaxSpeed             int
	TrackTurnSpeed       int
	SeekTurnSpeed        int
	TrackSpeed           int
	MaxIrFront           int
	MaxIrSide            int
	StrategyR1Time       int
	StrategyS1Time       int
	StrategyR2Time       int
	StrategyS2Time       int
	StrategyStraightTime int
	WaitTime             int
}

// Default Config data
func Default() Config {
	return Config{
		MaxSpeed:             100,
		TrackTurnSpeed:       100,
		SeekTurnSpeed:        100,
		TrackSpeed:           100,
		MaxIrFront:           40,
		MaxIrSide:            30,
		StrategyR1Time:       500,
		StrategyS1Time:       900,
		StrategyR2Time:       1000,
		StrategyS2Time:       900,
		StrategyStraightTime: 1200,
		WaitTime:             4800,
	}
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
	c.WaitTime *= 1000
}
