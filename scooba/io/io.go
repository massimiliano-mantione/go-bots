package io

import (
	"go-bots/ev3"
	"go-bots/scooba/config"
	"go-bots/scooba/logic"
	"time"
)

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var ml, mr, mfl, mfr *ev3.Attribute
var irL, irFL, irFR, irR *ev3.Attribute

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
		OutB: ev3.OutPortModeAuto,
		OutC: ev3.OutPortModeAuto,
		OutD: ev3.OutPortModeDcMotor,
	})
	data = d
	start = s

	// Ir Left
	ev3.CheckDriver(devs.In2, ev3.DriverIr, ev3.In1)
	// Ir FrontLeft
	ev3.CheckDriver(devs.In4, ev3.DriverIr, ev3.In2)
	// Ir FrontRight
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In3)
	// Ir Right
	ev3.CheckDriver(devs.In3, ev3.DriverIr, ev3.In4)

	// A front Left
	ev3.CheckDriver(devs.OutA, ev3.DriverRcxMotor, ev3.OutA)
	frontRight := devs.OutA
	// B left direct
	ev3.CheckDriver(devs.OutB, ev3.DriverTachoMotorLarge, ev3.OutB)
	// C right direct
	ev3.CheckDriver(devs.OutC, ev3.DriverTachoMotorLarge, ev3.OutC)
	// D front Right
	ev3.CheckDriver(devs.OutD, ev3.DriverRcxMotor, ev3.OutD)
	frontLeft := devs.OutD

	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.IrModeProx)
	ev3.SetMode(devs.In3, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.IrModeProx)

	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	irL = ev3.OpenByteR(devs.In2, ev3.BinData)
	irFL = ev3.OpenByteR(devs.In4, ev3.BinData)
	irFR = ev3.OpenByteR(devs.In1, ev3.BinData)
	irR = ev3.OpenByteR(devs.In3, ev3.BinData)
	// A front Left
	mfl = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	// B left direct
	ml = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	// C right direct
	mr = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	// D front Right
	mfr = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)

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

	// Wheels
	mr.Value = 0
	ml.Value = 0
	mr.Sync()
	ml.Sync()
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)

	// front Left
	mfl.Value = 0
	mfl.Sync()
	ev3.RunCommand(frontLeft, ev3.CmdStop)
	ev3.RunCommand(frontLeft, ev3.CmdRunDirect)

	// front Right
	mfr.Value = 0
	mfr.Sync()
	ev3.RunCommand(frontRight, ev3.CmdStop)
	ev3.RunCommand(frontRight, ev3.CmdRunDirect)
}

var speedL, speedR int
var lastMillis, currentMillis int

// ProcessCommand process the commands
func ProcessCommand(c *logic.Commands) {
	currentMillis = c.Millis
	lastMillis = currentMillis

	mlValue := speedL / 100
	mrValue := speedR / 100
	if mlValue > 100 {
		mlValue = 100
	}
	if mlValue < -100 {
		mlValue = -100
	}
	if mrValue > 100 {
		mrValue = 100
	}
	if mrValue < -100 {
		mrValue = -100
	}
	ml.Value = mlValue
	mr.Value = mrValue
	ml.Sync()
	mr.Sync()

	ledLG.Value = c.LedLeftGreen
	ledLR.Value = c.LedLeftRed
	ledRG.Value = c.LedRightGreen
	ledRR.Value = c.LedRightRed
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	mfl.Value = -config.FrontWheelsSpeed
	mfr.Value = config.FrontWheelsSpeed

	mfl.Sync()
	mfr.Sync()
}

// Loop contains the io loop
func Loop() {
	for {
		now := time.Now()
		millis := ev3.TimespanAsMillis(start, now)

		irL.Sync()
		irFL.Sync()
		irFR.Sync()
		irR.Sync()

		// fmt.Fprintln(os.Stderr, "DATA", irL.Value, irFL.Value, irFR.Value, irR.Value)

		data <- logic.Data{
			Start:             start,
			Millis:            millis,
			IrValueLeft:       irL.Value,
			IrValueFrontLeft:  irFL.Value,
			IrValueFrontRight: irFR.Value,
			IrValueRight:      irR.Value,
		}
	}
}

// Close terminates and cleans up the io module
func Close() {
	defer ev3.RunCommand(devs.OutA, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutB, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutC, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutD, ev3.CmdStop)

	defer ev3.RunCommand(devs.OutB, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutC, ev3.CmdStop)

	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	// TODO: close all files
	// mfl, mlr, ml, mr
	// irL, irFL, irFR, irR
}
