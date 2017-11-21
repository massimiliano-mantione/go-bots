package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config data
type Config struct {
	MaxSpeed int

	OutTimeMs   int
	OutPowerMax int
	OutPowerMin int
	OutPowerRn  int
	OutPowerRd  int

	MaxSteeringPC int
	MaxSlowPC     int
	SensorHole    int
	SensorRadius  int
	SensorSpan    int
	SensorMin     int
	MinDTicks     int
	MaxDTicks     int
	MinOutMillis  int
	DAvgMillis    int
	SpeedEstRn    int
	SpeedEstRd    int
	KPn           int
	KPd           int
	KDn           int
	KDd           int
	KIn           int
	KId           int
	KIrn          int
	KIrd          int
	KEReduction   int
	KELimit       int
	KEn           int
	KEd           int
	KErn          int
	KErd          int
	MaxPos        int
	MaxPosD       int
	SlowStart1    int
	SlowEnd1      int
	SlowSpeed1    int
	SlowStart2    int
	SlowEnd2      int
	SlowSpeed2    int
	SlowStart3    int
	SlowEnd3      int
	SlowSpeed3    int
	SlowStart4    int
	SlowEnd4      int
	SlowSpeed4    int
	Timeout       int
}

// CompleteConfig fills in computed configutation fields
func CompleteConfig(c *Config) {
	c.MaxPos = c.SensorRadius * 3
	c.MaxPosD = c.MaxPos / c.MinOutMillis

	if c.SlowSpeed1 > 0 {
		c.SlowStart1 *= 1000
		c.SlowEnd1 *= 1000
	} else {
		c.SlowSpeed1 = 0
		c.SlowStart1 = 0
		c.SlowEnd1 = 0
	}

	if c.SlowSpeed2 > 0 {
		c.SlowStart2 *= 1000
		c.SlowEnd2 *= 1000
	} else {
		c.SlowSpeed2 = 0
		c.SlowStart2 = 0
		c.SlowEnd2 = 0
	}

	if c.SlowSpeed3 > 0 {
		c.SlowStart3 *= 1000
		c.SlowEnd3 *= 1000
	} else {
		c.SlowSpeed3 = 0
		c.SlowStart3 = 0
		c.SlowEnd3 = 0
	}

	if c.SlowSpeed4 > 0 {
		c.SlowStart4 *= 1000
		c.SlowEnd4 *= 1000
	} else {
		c.SlowSpeed4 = 0
		c.SlowStart4 = 0
		c.SlowEnd4 = 0
	}

	if c.Timeout > 0 {
		c.Timeout *= 1000
	} else {
		c.Timeout = 0
	}
}

// Default Config data
func Default() Config {
	result := Config{
		// MaxSpeed:  100,
		MaxSpeed:      30,
		MaxSteeringPC: 120,
		MaxSlowPC:     60,
		SensorHole:    10,
		SensorRadius:  100,
		SensorSpan:    700,
		SensorMin:     80,
		MinDTicks:     10,
		MaxDTicks:     30000,
		MinOutMillis:  10,
		DAvgMillis:    42,
		KPn:           40,
		KPd:           100,
		KDn:           20,
		KDd:           100,
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
