package io

import (
	"go-bots/ev3"
	"go-bots/xl4/logic"
	"go-bots/xl4/vision"
	"time"
)

func colorIsOut(v int) bool {
	return v > 20
}

var devs *ev3.Devices
var data chan<- logic.Data
var commands <-chan logic.Commands

var mr1, mr2, ml1, ml2 *ev3.Attribute
var colR, colL, irR, irL *ev3.Attribute

var start time.Time

func StartTime() time.Time {
	return start
}

func Init(d chan<- logic.Data, c <-chan logic.Commands) {
	devs = ev3.Scan()
	data = d
	commands = c

	// IR left
	ev3.CheckDriver(devs.In1, ev3.DriverIr)
	// Color left
	ev3.CheckDriver(devs.In2, ev3.DriverColor)
	// IR right
	ev3.CheckDriver(devs.In3, ev3.DriverIr)
	// Color right
	ev3.CheckDriver(devs.In4, ev3.DriverColor)

	// Left back inverted
	ev3.CheckDriver(devs.OutA, ev3.DriverRcxMotor)
	// Left front direct
	ev3.CheckDriver(devs.OutB, ev3.DriverRcxMotor)
	// Right front direct
	ev3.CheckDriver(devs.OutC, ev3.DriverRcxMotor)
	// Right back direct
	ev3.CheckDriver(devs.OutD, ev3.DriverRcxMotor)

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
	mr1 = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	mr2 = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	ml1 = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	ml2 = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)

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

func executor() {
	for c := range commands {
		// log.Println("Execute:", c)

		mr1.Value = c.SpeedRight
		mr2.Value = c.SpeedRight
		ml1.Value = c.SpeedLeft
		ml2.Value = -c.SpeedLeft
		mr1.Sync()
		mr2.Sync()
		ml1.Sync()
		ml2.Sync()
	}
}

func Loop() {
	start = time.Now()
	sensorTicks := time.Tick(10 * time.Millisecond)

	go executor()

	for t := range sensorTicks {
		millis := ev3.TimespanAsMillis(start, t)

		colR.Sync()
		colL.Sync()
		irR.Sync()
		irL.Sync()

		vision.Process(millis, irR.Value, irL.Value)

		data <- logic.Data{
			Start:            start,
			Millis:           millis,
			CornerRightIsOut: colorIsOut(colR.Value),
			CornerLeftIsOut:  colorIsOut(colL.Value),
		}
	}
}

func Close() {
	defer ev3.RunCommand(devs.OutA, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutB, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutC, ev3.CmdStop)
	defer ev3.RunCommand(devs.OutD, ev3.CmdStop)

	// TODO: close all files
	// pf, mf, ml, mc, mr
	// colR, colL, irR, irL
}
