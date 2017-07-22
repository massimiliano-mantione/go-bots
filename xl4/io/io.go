package io

import (
	"go-bots/ev3"
	"go-bots/xl4/config"
	"go-bots/xl4/logic"
	"time"
)

func colorIsOut(v int) bool {
	return v > config.ColorIsOut
}

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var mr1, mr2, ml1, ml2 *ev3.Attribute
var colR, colL, irR, irL *ev3.Attribute

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

var start time.Time

// StartTime gets the time when the bot started
func StartTime() time.Time {
	return start
}

// Init initializes the io module
func Init(d chan<- logic.Data, s time.Time) {
	devs = ev3.Scan(&ev3.OutPortModes{
		OutA: ev3.OutPortModeDcMotor,
		OutB: ev3.OutPortModeDcMotor,
		OutC: ev3.OutPortModeDcMotor,
		OutD: ev3.OutPortModeDcMotor,
	})
	data = d
	start = s

	// IR left
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	// Color left
	ev3.CheckDriver(devs.In2, ev3.DriverColor, ev3.In2)
	// IR right
	ev3.CheckDriver(devs.In3, ev3.DriverIr, ev3.In3)
	// Color right
	ev3.CheckDriver(devs.In4, ev3.DriverColor, ev3.In4)

	// Right back inverted
	ev3.CheckDriver(devs.OutA, ev3.DriverRcxMotor, ev3.OutA)
	// Right front direct
	ev3.CheckDriver(devs.OutB, ev3.DriverRcxMotor, ev3.OutB)
	// Left front direct
	ev3.CheckDriver(devs.OutC, ev3.DriverRcxMotor, ev3.OutC)
	// Left back direct
	ev3.CheckDriver(devs.OutD, ev3.DriverRcxMotor, ev3.OutD)

	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.ColorModeReflect)
	ev3.SetMode(devs.In3, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.ColorModeReflect)

	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutB, ev3.CmdStop)
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	irR = ev3.OpenByteR(devs.In3, ev3.BinData)
	irL = ev3.OpenByteR(devs.In1, ev3.BinData)
	colR = ev3.OpenByteR(devs.In4, ev3.BinData)
	colL = ev3.OpenByteR(devs.In2, ev3.BinData)
	mr2 = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	mr1 = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	ml1 = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	ml2 = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)

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

	mr1.Value = 0
	mr2.Value = 0
	ml1.Value = 0
	ml2.Value = 0

	mr1.Sync()
	mr2.Sync()
	ml1.Sync()
	ml2.Sync()

	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)
}

func computeSpeed(currentSpeed int, targetSpeed int, millis int) int {
	if currentSpeed < targetSpeed {
		currentSpeed += (config.ForwardAcceleration * millis)
		if currentSpeed > targetSpeed {
			currentSpeed = targetSpeed
		}
	}
	if currentSpeed > targetSpeed {
		currentSpeed -= (config.ReverseAcceleration * millis)
		if currentSpeed < targetSpeed {
			currentSpeed = targetSpeed
		}
	}
	return currentSpeed
}

var speedRight, speedLeft int
var lastMillis, currentMillis int

func ProcessCommand(c *logic.Commands) {
	currentMillis = c.Millis
	millis := currentMillis - lastMillis
	speedRight = computeSpeed(speedRight, c.SpeedRight, millis)
	speedLeft = computeSpeed(speedLeft, c.SpeedLeft, millis)
	lastMillis = currentMillis

	mr1.Value = speedRight / 100
	mr2.Value = -speedRight / 100
	ml1.Value = speedLeft / 100
	ml2.Value = speedLeft / 100
	mr1.Sync()
	mr2.Sync()
	ml1.Sync()
	ml2.Sync()

	ledLG.Value = c.LedLeftGreen
	ledLR.Value = c.LedLeftRed
	ledRG.Value = c.LedRightGreen
	ledRR.Value = c.LedRightRed
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()
}

// Loop contains the io loop
func Loop() {
	for {
		now := time.Now()
		millis := ev3.TimespanAsMillis(start, now)

		colR.Sync()
		colL.Sync()
		irR.Sync()
		irL.Sync()
		// fmt.Fprintln(os.Stderr, "DATA", irL.Value, irR.Value)
		// intensity, angle := vision.Process(millis, irL.Value, irR.Value)

		data <- logic.Data{
			Start:            start,
			Millis:           millis,
			CornerRightIsOut: colorIsOut(colR.Value),
			CornerLeftIsOut:  colorIsOut(colL.Value),
			CornerRight:      colR.Value,
			CornerLeft:       colL.Value,
			IrLeftValue:      irL.Value,
			IrRightValue:     irR.Value,
		}
	}
}

// Close terminates and cleans up the io module
func Close() {
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
