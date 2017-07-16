package logic

import (
	"go-bots/ui"
	"go-bots/xl4/config"
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
var keys <-chan ui.KeyEvent
var quit chan<- bool

// Init initializes the logic module
func Init(d <-chan Data, c func(*Commands), k <-chan ui.KeyEvent, q chan<- bool) {
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
	c.LedLeftGreen = 255 * d.VisionIntensity * ((config.VisionAngleMax - d.VisionAngle) / 2) / (config.VisionIntensityMax * config.VisionAngleMax)
	c.LedRightGreen = 255 * d.VisionIntensity * ((config.VisionAngleMax + d.VisionAngle) / 2) / (config.VisionIntensityMax * config.VisionAngleMax)
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

func waitBeginOrQuit(start int) {
	for {
		select {
		case d := <-data:
			now := d.Millis
			c.Millis = now
			speed(0, 0)
			ledsFromData(d)

			cmd()
		case k := <-keys:
			if k.Key == ui.Quit || k.Key == ui.Back {
				quit <- true
				return
			}
			if k.Key == ui.Enter {
				go pauseBeforeBegin(k.Millis)
				return
			}
		}
	}
}

func pauseBeforeBegin(start int) {
	for {
		select {
		case d := <-data:
			now := d.Millis
			c.Millis = now
			elapsed := now - start
			if elapsed >= config.StartTime {
				go begin(now)
				return
			}
			speed(0, 0)
			intensity := ((elapsed % 1000) * 255) / (config.StartTime / 5)
			if elapsed > (config.StartTime * 4 / 5) {
				leds(intensity, intensity, intensity, intensity)
			} else {
				leds(0, 0, intensity, intensity)
			}
			cmd()
		case k := <-keys:
			if k.Key == ui.Quit || k.Key == ui.Back {
				quit <- true
				return
			}
		}
	}
}

func begin(start int) {
	for {
		select {
		case d := <-data:
			now := d.Millis
			c.Millis = now
			elapsed := now - start

			adjustInner := d.CornerLeft * config.AdjustInnerMax / 100
			inner := config.InnerSpeed - adjustInner

			if elapsed >= 5000 {
				quit <- true
				return
			}
			speed(config.OuterSpeed, inner)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if k.Key == ui.Quit || k.Key == ui.Back {
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
