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
	IrValueRight     int
	IrValueLeft      int
	VisionIntensity  int
	VisionAngle      int
}

// Commands contains commands for motors and leds
type Commands struct {
	Millis        int
	SpeedLeft     int
	SpeedRight    int
	FrontActive   bool
	EyesActive    bool
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

// Run starts the logic module
func Run() {
	go chooseStrategy(0)
}
