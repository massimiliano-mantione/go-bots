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
	MaxIrValue           int
	StrategyR1Time       int
	StrategyS1Time       int
	StrategyR2Time       int
	StrategyS2Time       int
	StrategyStraightTime int
}

// Default Config data
func Default() Config {
	return Config{
		MaxSpeed:             1000000,
		TrackTurnSpeed:       400000,
		SeekTurnSpeed:        400000,
		TrackSpeed:           1000000,
		MaxIrValue:           30,
		StrategyR1Time:       220000,
		StrategyS1Time:       120000,
		StrategyR2Time:       470000,
		StrategyS2Time:       120000,
		StrategyStraightTime: 400000,
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
