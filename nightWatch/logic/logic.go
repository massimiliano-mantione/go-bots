package logic

import (
	"go-bots/ui"
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
var keys <-chan ui.Key

func Init(d <-chan Data, c chan<- Commands, k <-chan ui.Key) {
	data = d
	commands = c
	keys = k
}

func Run() {
	var lastMillis int
	var targetMillis int
	const deltaMillis = 50
	const slowSpeed = 60
	const fastSpeed = 100

	for {
		select {
		case d := <-data:
			// log.Println("data", d)
			commands <- Commands{
				SpeedFront: 0,
				SpeedRight: 0,
				SpeedLeft:  0,
			}
			lastMillis = d.Millis
			if targetMillis > 0 && lastMillis <= targetMillis {
				targetMillis = 0
				commands <- Commands{
					SpeedFront: 0,
				}
			}
		case k := <-keys:
			switch k {
			case ui.Up:
				commands <- Commands{
					SpeedFront: fastSpeed,
				}
				targetMillis = lastMillis + deltaMillis
			case ui.Down:
				commands <- Commands{
					SpeedFront: -fastSpeed,
				}
				targetMillis = lastMillis + deltaMillis
			case ui.Right:
				commands <- Commands{
					SpeedFront: slowSpeed,
				}
				targetMillis = lastMillis + deltaMillis
			case ui.Left:
				commands <- Commands{
					SpeedFront: -slowSpeed,
				}
				targetMillis = lastMillis + deltaMillis
			case ui.Quit:
				log.Println("Logic got clean quit")
				return
			}
		}
	}
}
