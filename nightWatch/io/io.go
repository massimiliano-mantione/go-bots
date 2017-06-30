package io

import (
	"go-bots/ev3"
	"go-bots/nightWatch/logic"
	"time"
)

func colorIsOut(v int) bool {
	return v > 20
}

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var pf, ml, mc, mr *ev3.Attribute
var df string
var colR, colL, irR, irL *ev3.Attribute

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

var start time.Time

// StartTime gets the time when the bot started
func StartTime() time.Time {
	return start
}

// Init initializes the io module
func Init(d chan<- logic.Data) {
	devs = ev3.Scan(nil)
	data = d

	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	ev3.CheckDriver(devs.In2, ev3.DriverIr, ev3.In2)
	ev3.CheckDriver(devs.In3, ev3.DriverColor, ev3.In3)
	ev3.CheckDriver(devs.In4, ev3.DriverColor, ev3.In4)

	// A center
	ev3.CheckDriver(devs.OutA, ev3.DriverTachoMotorLarge, ev3.OutA)
	// B left
	ev3.CheckDriver(devs.OutB, ev3.DriverTachoMotorLarge, ev3.OutB)
	// C front
	ev3.CheckDriver(devs.OutC, ev3.DriverTachoMotorMedium, ev3.OutC)
	// D right
	ev3.CheckDriver(devs.OutD, ev3.DriverTachoMotorLarge, ev3.OutD)

	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.IrModeProx)
	ev3.SetMode(devs.In3, ev3.ColorModeReflect)
	ev3.SetMode(devs.In4, ev3.ColorModeReflect)

	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)

	irR = ev3.OpenByteR(devs.In1, ev3.BinData)
	irL = ev3.OpenByteR(devs.In2, ev3.BinData)
	colR = ev3.OpenByteR(devs.In3, ev3.BinData)
	colL = ev3.OpenByteR(devs.In4, ev3.BinData)
	// D right
	mr = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	// A center
	mc = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	// B left
	ml = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	// C front
	df = devs.OutC
	pf = ev3.OpenTextR(devs.OutC, ev3.Position)

	ledLG = ev3.OpenTextW(devs.LedLeftGreen, ev3.Brightness)
	ledLR = ev3.OpenTextW(devs.LedLeftRed, ev3.Brightness)
	ledRG = ev3.OpenTextW(devs.LedRightGreen, ev3.Brightness)
	ledRR = ev3.OpenTextW(devs.LedRightRed, ev3.Brightness)
	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	mr.Value = 0
	mc.Value = 0
	ml.Value = 0

	mr.Sync()
	mc.Sync()
	ml.Sync()

	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)
	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)

	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdRunForever)
}

const forwardAcceleration = 10000 / 400
const reverseAcceleration = 10000 / 1
const angularAcceleration = 100 / 1

var speed, direction int
var lastMillis, currentMillis int

var frontActive bool

func setFrontActive(active bool) {
	if active != frontActive {
		speed := "0"
		if active {
			// max is 1560
			speed = "560"
		}
		ev3.WriteStringAttribute(df, ev3.SpeedSp, speed)
		ev3.RunCommand(df, ev3.CmdRunForever)
	}
	frontActive = active
}

func computeSpeed(currentSpeed int, targetSpeed int, millis int) int {
	if currentSpeed < targetSpeed {
		currentSpeed += (forwardAcceleration * millis)
		if currentSpeed > targetSpeed {
			currentSpeed = targetSpeed
		}
	}
	if currentSpeed > targetSpeed {
		currentSpeed -= (reverseAcceleration * millis)
		if currentSpeed < targetSpeed {
			currentSpeed = targetSpeed
		}
	}
	return currentSpeed
}
func computeAngle(currentAngle int, targetAngle int, millis int) int {
	if currentAngle < targetAngle {
		currentAngle += (angularAcceleration * millis)
		if currentAngle > targetAngle {
			currentAngle = targetAngle
		}
	}
	if currentAngle > targetAngle {
		currentAngle -= (angularAcceleration * millis)
		if currentAngle < targetAngle {
			currentAngle = targetAngle
		}
	}
	return currentAngle
}

func computeSpeeds(currentSpeed int, currentAngle int) (left int, center int, right int) {
	left = currentSpeed
	center = currentSpeed
	right = currentSpeed
	if currentAngle == 0 {
		return
	}

	if currentAngle > 0 {
		right *= (logic.AngleMax - currentAngle)
		right /= logic.AngleMax
	} else if currentAngle < 0 {
		left *= (logic.AngleMax + currentAngle)
		left /= logic.AngleMax
	}
	center = (left + right) / 2
	return
}

func ProcessCommand(c *logic.Commands) {
	currentMillis = c.Millis
	millis := currentMillis - lastMillis
	speed = computeSpeed(speed, c.Speed, millis)
	direction = computeAngle(direction, c.Direction, millis)
	lastMillis = currentMillis

	left, center, right := computeSpeeds(speed, direction)

	ml.Value = -left / 100
	mc.Value = -center / 100
	mr.Value = -right / 100
	ml.Sync()
	mc.Sync()
	mr.Sync()

	ledLG.Value = c.LedLeftGreen
	ledLR.Value = c.LedLeftRed
	ledRG.Value = c.LedRightGreen
	ledRR.Value = c.LedRightRed
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	setFrontActive(c.FrontActive)
}

// Loop contains the io loop
func Loop() {
	start = time.Now()
	for {
		now := time.Now()
		millis := ev3.TimespanAsMillis(start, now)

		pf.Sync()
		colR.Sync()
		colL.Sync()
		irR.Sync()
		irL.Sync()

		// vision.Process(millis, pf.Value, irR.Value, irL.Value)

		data <- logic.Data{
			Start:            start,
			Millis:           millis,
			CornerRightIsOut: colorIsOut(colR.Value),
			CornerLeftIsOut:  colorIsOut(colL.Value),
			CornerRight:      colR.Value,
			CornerLeft:       colL.Value,
			VisionIntensity:  0,
			VisionAngle:      0,
		}
	}
}

// Close terminates and cleans up the io module
func Close() {
	defer ev3.RunCommand(devs.OutA, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutB, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutC, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutD, ev3.CmdReset)

	defer ev3.RunCommand(devs.OutA, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutB, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutC, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutD, ev3.CmdStop)

	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	// TODO: close all files
	// pf, mf, ml, mc, mr
	// colR, colL, irR, irL
}
