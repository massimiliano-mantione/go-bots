package io

import (
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"go-bots/seeker2/logic"
	"go-bots/seeker2/vision"
	"time"
)

func colorIsOut(v int) bool {
	return v > 5
}

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var pme, pmesp, ml, mr, mf *ev3.Attribute
var dme, dmf string
var colR, colL, irR, irL *ev3.Attribute

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

var start time.Time

var eyesAreActive bool
var eyesAreScanning bool
var eyesDirection ev3.Direction

func setEyesState(active bool, scanning bool, dir ev3.Direction) {
	eyesAreActive, eyesAreScanning, eyesDirection = active, scanning, dir
	desiredSetPosition := config.VisionStartPosition
	if eyesAreActive {
		if eyesAreScanning {
			desiredSetPosition = int(config.VisionMaxPosition * dir)
		} else {
			desiredSetPosition = int(config.VisionTurnPosition * dir)
		}
	}
	if pmesp.Value != desiredSetPosition {
		pmesp.Value = desiredSetPosition
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
	// C left direct
	ev3.CheckDriver(devs.OutC, ev3.DriverRcxMotor, ev3.OutC)
	// D right inverted
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
	// C left direct
	ml = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	// D right inverted
	mr = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
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
	ev3.WriteStringAttribute(dme, ev3.Position, config.VisionStartPositionString)
	ev3.WriteStringAttribute(dme, ev3.SpeedSp, config.VisionSpeed)
	ev3.WriteStringAttribute(dme, ev3.StopAction, "hold")
	setEyesState(false, false, ev3.NoDirection)
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

	mlValue := speedL / 100
	mrValue := -speedR / 100
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

	if c.FrontActive {
		mf.Value = config.FrontWheelsSpeed
	} else {
		mf.Value = 0
	}
	mf.Sync()

	// fmt.Fprintln(os.Stderr, "DATA EYES ACTIVE", c.EyesActive)

	if c.EyesActive {
		if c.EyesDirection == ev3.NoDirection {
			if eyesDirection == ev3.NoDirection {
				eyesDirection = ev3.Right
			}
			setEyesState(true, true, eyesDirection)
		} else {
			setEyesState(true, false, c.EyesDirection)
		}
	} else {
		vision.Reset()
		setEyesState(false, false, ev3.NoDirection)
	}
}

// Loop contains the io loop
func Loop() {
	lastEyesIntensity := 0
	for {
		now := time.Now()
		millis := ev3.TimespanAsMillis(start, now)

		pme.Sync()
		colR.Sync()
		colL.Sync()
		irR.Sync()
		irL.Sync()

		visionIntensity, visionAngle := 0, 0
		if eyesAreActive {
			if eyesAreScanning {
				var desiredEyesDirection ev3.Direction
				if lastEyesIntensity > 0 {
					visionIntensity, visionAngle, desiredEyesDirection = vision.ProcessScan(millis, eyesDirection, pme.Value, irL.Value, irR.Value)
				} else {
					visionIntensity, visionAngle, desiredEyesDirection = vision.ProcessSeekScan(millis, eyesDirection, pme.Value, irL.Value, irR.Value)
				}
				setEyesState(true, true, desiredEyesDirection)
			} else {
				visionIntensity, visionAngle = vision.ProcessTurn(millis, eyesDirection, pme.Value, irL.Value, irR.Value, lastEyesIntensity > 0)
				setEyesState(true, false, eyesDirection)
			}
			lastEyesIntensity = visionIntensity
		} else {
			lastEyesIntensity = 0
		}

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
			VisionIntensity:  visionIntensity,
			VisionAngle:      visionAngle,
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
