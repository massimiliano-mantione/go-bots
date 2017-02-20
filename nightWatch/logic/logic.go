package logic

import (
	"log"
	"time"
)

type Data struct {
	Start            time.Time
	Millis           int
	CornerRightIsOut bool
	CornerLeftIsOut  bool
	VisionIntensity  int
	VisionAngle      int
}

type Commands struct {
	SpeedFront int
	SpeedRight int
	SpeedLeft  int
}

var data <-chan Data
var commands chan<- Commands

func Init(d <-chan Data, c chan<- Commands) {
	data = d
	commands = c
}

func Run() {
	select {
	case d := <-data:
		log.Println("data", d)
		if d.Millis > 2000 {
			return
		}
	}
}
