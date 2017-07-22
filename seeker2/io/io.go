package io

import (
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"go-bots/seeker2/logic"
	"time"
)

func colorIsOut(v int) bool {
	return v > 20
}

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var pme, pmesp, ml, mr, mf *ev3.Attribute
var dme, dmf string
var colR, colL, irR, irL *ev3.Attribute

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

var start time.Time

var eyesTargetPosition int

func goToEyesPosition(p int) {
	if p != eyesTargetPosition {
		eyesTargetPosition = p
		pmesp.Value = p
		pmesp.Sync()
		ev3.RunCommand(dme, ev3.CmdRunToAbsPos)
	}
}

// StartTime gets the time when the bot started
func StartTime() time.Time {
	return start
}

// Init initializes the io module
func Init(d chan<- logic.Data, s time.Time) {
	devs = ev3.Scan(&ev3.OutPortModes{
		OutA: ev3.OutPortModeAuto,
		OutB: ev3.OutPortModeAuto,
		OutC: ev3.OutPortModeDcMotor,
		OutD: ev3.OutPortModeDcMotor,
	})
	data = d
	start = s

	// Col L
	ev3.CheckDriver(devs.In1, ev3.DriverColor, ev3.In1)
	// Col R
	ev3.CheckDriver(devs.In2, ev3.DriverColor, ev3.In2)
	// Ir L
	ev3.CheckDriver(devs.In3, ev3.DriverIr, ev3.In3)
	// Ir R
	ev3.CheckDriver(devs.In4, ev3.DriverIr, ev3.In4)

	// A front
	ev3.CheckDriver(devs.OutA, ev3.DriverTachoMotorMedium, ev3.OutA)
	// B eyes
	ev3.CheckDriver(devs.OutB, ev3.DriverTachoMotorMedium, ev3.OutB)
	// C right direct
	ev3.CheckDriver(devs.OutC, ev3.DriverRcxMotor, ev3.OutC)
	// D left inverted
	ev3.CheckDriver(devs.OutD, ev3.DriverRcxMotor, ev3.OutD)

	ev3.SetMode(devs.In1, ev3.ColorModeReflect)
	ev3.SetMode(devs.In2, ev3.ColorModeReflect)
	ev3.SetMode(devs.In3, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.IrModeProx)

	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	colL = ev3.OpenByteR(devs.In1, ev3.BinData)
	colR = ev3.OpenByteR(devs.In2, ev3.BinData)
	irL = ev3.OpenByteR(devs.In3, ev3.BinData)
	irR = ev3.OpenByteR(devs.In4, ev3.BinData)
	// C right direct
	mr = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	// D left inverted
	ml = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	// B eyes
	dme = devs.OutB
	pme = ev3.OpenTextR(devs.OutB, ev3.Position)
	pmesp = ev3.OpenTextW(devs.OutB, ev3.PositionSp)
	// A front
	dmf = devs.OutA
	mf = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)

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
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)

	// Front
	ev3.RunCommand(dmf, ev3.CmdReset)
	mf.Value = 0
	mf.Sync()
	ev3.RunCommand(dmf, ev3.CmdRunDirect)

	// Eyes
	ev3.RunCommand(dme, ev3.CmdReset)
	ev3.RunCommand(dme, ev3.CmdRunForever)
	eyesTargetPosition = 0
	pmesp.Value = 0
	pmesp.Sync()
	ev3.RunCommand(dme, ev3.CmdRunToAbsPos)
}

var speedL, speedR int
var lastMillis, currentMillis int

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

func ProcessCommand(c *logic.Commands) {
	currentMillis = c.Millis
	millis := currentMillis - lastMillis
	speedL = computeSpeed(speedL, c.SpeedLeft, millis)
	speedR = computeSpeed(speedR, c.SpeedRight, millis)
	lastMillis = currentMillis

	ml.Value = -speedL / 100
	mr.Value = speedR / 100
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

	if c.FrontActive {
		mf.Value = config.FrontWheelsSpeed
	} else {
		mf.Value = 0
	}
	mf.Sync()

	if !c.EyesActive {
		goToEyesPosition(0)
	}
}

// Loop contains the io loop
func Loop() {
	for {
		now := time.Now()
		millis := ev3.TimespanAsMillis(start, now)

		// pf.Sync()
		colR.Sync()
		colL.Sync()
		irR.Sync()
		irL.Sync()

		// vision.Process(millis, pf.Value, irR.Value, irL.Value)
		// fmt.Fprintln(os.Stderr, "DATA", colL.Value, colR.Value, irL.Value, irR.Value)

		data <- logic.Data{
			Start:            start,
			Millis:           millis,
			CornerRightIsOut: colorIsOut(colR.Value),
			CornerLeftIsOut:  colorIsOut(colL.Value),
			CornerRight:      colR.Value,
			CornerLeft:       colL.Value,
			IrValueRight:     irR.Value,
			IrValueLeft:      irL.Value,
			VisionIntensity:  0,
			VisionAngle:      0,
		}
	}
}

// Close terminates and cleans up the io module
func Close() {
	defer ev3.RunCommand(devs.OutA, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutB, ev3.CmdReset)
	defer ev3.RunCommand(devs.OutC, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutD, ev3.CmdStop)

	defer ev3.RunCommand(devs.OutA, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutB, ev3.CmdStop)

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
