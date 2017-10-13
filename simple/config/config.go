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
		TrackTurnSpeed:       1000000 / 2,
		SeekTurnSpeed:        1000000 / 2,
		TrackSpeed:           1000000,
		MaxIrValue:           50,
		StrategyR1Time:       500000,
		StrategyS1Time:       500000,
		StrategyR2Time:       500000,
		StrategyS2Time:       500000,
		StrategyStraightTime: 1000000,
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
