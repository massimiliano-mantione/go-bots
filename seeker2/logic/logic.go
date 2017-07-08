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

// SpeedMax is the maximum speed
const SpeedMax = 10000
const SpeedReverse = 5000
const SpeedTurning = 5000

// AngleMax is the maximum turning angle (positive on the right)
const AngleMax = 100

// Commands contains commands for motors and leds
type Commands struct {
	Millis        int
	SpeedL        int
	SpeedR        int
	FrontActive   bool
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

func movement(speedL int, speedR int) {
	c.SpeedL = speedL
	c.SpeedR = speedR
}

func ledsOff() {
	c.LedLeftGreen = 0
	c.LedRightGreen = 0
	c.LedLeftRed = 0
	c.LedRightRed = 0
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

const startTime = 50

func waitBeginOrQuit(start int) {
	now := start
	for {
		select {
		case d := <-data:
			now = d.Millis
			c.Millis = now
			movement(0, 0)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if k.Key == ui.Quit || k.Key == ui.Back {
				quit <- true
				return
			}
			if k.Key == ui.Enter {
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
			c.Millis = now
			elapsed := now - start
			if elapsed >= startTime {
				go begin(now)
				return
			}
			movement(0, 0)
			intensity := ((elapsed % 1000) * 255) / (startTime / 5)
			if elapsed > (startTime * 4 / 5) {
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
	now := start
	speedL := 0
	speedR := 0
	for {
		select {
		case d := <-data:
			now = d.Millis
			c.Millis = now
			// elapsed := now - start
			movement(speedL, speedR)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			switch k.Key {
			case ui.Enter:
				speedL = 0
				speedR = 0
				c.FrontActive = false
			case ui.Up:
				speedL = SpeedMax
				speedR = SpeedMax
				c.FrontActive = true
			case ui.Down:
				speedL = -SpeedReverse
				speedR = -SpeedReverse
				c.FrontActive = true
			case ui.Left:
				speedL = -SpeedTurning
				speedR = SpeedTurning
				c.FrontActive = true
			case ui.Right:
				speedL = SpeedTurning
				speedR = -SpeedTurning
				c.FrontActive = true
			case ui.Quit:
				quit <- true
				return
			case ui.Back:
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
