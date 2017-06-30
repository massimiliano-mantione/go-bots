package logic

import (
	"go-bots/ui"
	"time"
)

// Data contains readings from sensors (adjusted for logic)
type Data struct {
	Start            time.Time
	Millis           int
	CornerRightIsOut bool
	CornerLeftIsOut  bool
	CornerRight      int
	CornerLeft       int
	VisionIntensity  int
	VisionAngle      int
}

// VisionIntensityMax is the maximum vision intensity
const VisionIntensityMax = 100

// VisionAngleMax is the maximum vision angle (positive on the right)
const VisionAngleMax = 100

// Commands contains commands for motors and leds
type Commands struct {
	Millis        int
	SpeedRight    int
	SpeedLeft     int
	LedRightRed   int
	LedRightGreen int
	LedLeftRed    int
	LedLeftGreen  int
}

var data <-chan Data
var commandProcessor func(*Commands)
var keys <-chan ui.Key
var quit chan<- bool

// Init initializes the logic module
func Init(d <-chan Data, c func(*Commands), k <-chan ui.Key, q chan<- bool) {
	data = d
	commandProcessor = c
	keys = k
	quit = q
}

var c = Commands{}

func cmd() {
	commandProcessor(&c)
}

func speed(left int, right int) {
	c.SpeedLeft = left
	c.SpeedRight = right
}

func leds(leftGreen int, rightGreen int, leftRed int, rightRed int) {
	c.LedLeftGreen = leftGreen
	c.LedRightGreen = rightGreen
	c.LedLeftRed = leftRed
	c.LedRightRed = rightRed
}

func ledsFromData(d Data) {
	c.LedLeftGreen = 255 * d.VisionIntensity * ((VisionAngleMax - d.VisionAngle) / 2) / (VisionIntensityMax * VisionAngleMax)
	c.LedRightGreen = 255 * d.VisionIntensity * ((VisionAngleMax + d.VisionAngle) / 2) / (VisionIntensityMax * VisionAngleMax)
	if d.CornerLeftIsOut {
		c.LedLeftRed = 255
	} else {
		c.LedLeftRed = 0
	}
	if d.CornerRightIsOut {
		c.LedRightRed = 255
	} else {
		c.LedRightRed = 0
	}
}

const maxSpeed = 10000

const startTime = 5000

func waitBeginOrQuit(start int) {
	now := start
	for {
		select {
		case d := <-data:
			now = d.Millis
			speed(0, 0)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if k == ui.Quit || k == ui.Back {
				quit <- true
				return
			}
			if k == ui.Enter {
				go pauseBeforeBegin(now)
				return
			}
		}
	}
}

func pauseBeforeBegin(start int) {
	now := start
	for {
		select {
		case d := <-data:
			now = d.Millis
			elapsed := now - start
			if elapsed >= startTime {
				go begin(now)
				return
			}
			speed(0, 0)
			intensity := ((elapsed % 1000) * 255) / (startTime / 5)
			if elapsed > (startTime * 4 / 5) {
				leds(intensity, intensity, intensity, intensity)
			} else {
				leds(0, 0, intensity, intensity)
			}
			cmd()
		case k := <-keys:
			if k == ui.Quit || k == ui.Back {
				quit <- true
				return
			}
		}
	}
}

func begin(start int) {
	now := start
	for {
		select {
		case d := <-data:
			now = d.Millis
			elapsed := now - start
			if elapsed >= 3000 {
				quit <- true
				return
			}
			speed(maxSpeed, maxSpeed)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if k == ui.Quit || k == ui.Back {
				quit <- true
				return
			}
		}
	}
}

// Run starts the logic module
func Run() {
	go waitBeginOrQuit(0)
}
