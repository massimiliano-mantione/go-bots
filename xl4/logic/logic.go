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

const rampMillis = 500

func computeSpeed(rightTargetSpeed int, leftTargetSpeed int, currentMillis int, startMillis int) (int, int) {
	deltaTime := currentMillis - startMillis
	if deltaTime < rampMillis {
		return (rightTargetSpeed * deltaTime) / rampMillis, (leftTargetSpeed * deltaTime) / rampMillis
	}
	return rightTargetSpeed, leftTargetSpeed
}

func Run() {
	var currentMillis int
	var startMillis int
	var rightTargetSpeed int
	var leftTargetSpeed int
	const slowSpeed = 60
	const fastSpeed = 100

	for {
		select {
		case d := <-data:
			currentMillis = d.Millis
			r, l := computeSpeed(rightTargetSpeed, leftTargetSpeed, currentMillis, startMillis)
			commands <- Commands{
				SpeedRight: r,
				SpeedLeft:  l,
			}
		case k := <-keys:
			switch k {
			case ui.Up:
				startMillis = currentMillis
				rightTargetSpeed = fastSpeed
				leftTargetSpeed = fastSpeed
			case ui.Down:
				startMillis = currentMillis
				rightTargetSpeed = -fastSpeed
				leftTargetSpeed = -fastSpeed
			case ui.Right:
				startMillis = currentMillis
				rightTargetSpeed = -fastSpeed
				leftTargetSpeed = fastSpeed
			case ui.Left:
				startMillis = currentMillis
				rightTargetSpeed = fastSpeed
				leftTargetSpeed = -fastSpeed
			case ui.Enter:
				startMillis = currentMillis
				rightTargetSpeed = 0
				leftTargetSpeed = 0
			case ui.Quit:
				log.Println("Logic got clean quit (kill)")
				return
			case ui.Back:
				log.Println("Logic got clean quit (back)")
				return
			}
		}
	}
}
