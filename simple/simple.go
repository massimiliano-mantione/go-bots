package main

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/simple/config"
	"os"
	"time"
)

var devs *ev3.Devices
var initializationTime time.Time
var motorL1, motorL2, motorR1, motorR2 *ev3.Attribute
var irL, irFL, irFR, irR *ev3.Attribute

func initialize() {
	initializationTime = time.Now()

	devs = ev3.Scan(&ev3.OutPortModes{
		OutA: ev3.OutPortModeDcMotor,
		OutB: ev3.OutPortModeDcMotor,
		OutC: ev3.OutPortModeDcMotor,
		OutD: ev3.OutPortModeDcMotor,
	})

	// Check motors
	ev3.CheckDriver(devs.OutA, ev3.DriverRcxMotor, ev3.OutA)
	ev3.CheckDriver(devs.OutB, ev3.DriverRcxMotor, ev3.OutB)
	ev3.CheckDriver(devs.OutC, ev3.DriverRcxMotor, ev3.OutC)
	ev3.CheckDriver(devs.OutD, ev3.DriverRcxMotor, ev3.OutD)

	// Check sensors
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	ev3.CheckDriver(devs.In2, ev3.DriverIr, ev3.In2)
	ev3.CheckDriver(devs.In3, ev3.DriverIr, ev3.In3)
	ev3.CheckDriver(devs.In4, ev3.DriverIr, ev3.In4)

	// Set sensors mode (for remote control)
	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.IrModeProx)
	ev3.SetMode(devs.In3, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.IrModeProx)

	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutB, ev3.CmdStop)
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	// Open motors
	motorL1 = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	motorL2 = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	motorR1 = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	motorR2 = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)

	// Reset motor speed
	motorL1.Value = 0
	motorL2.Value = 0
	motorR1.Value = 0
	motorR2.Value = 0

	motorL1.Sync()
	motorL2.Sync()
	motorR1.Sync()
	motorR2.Sync()

	// Put motors in direct mode
	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)

	// Open sensors
	irL = ev3.OpenByteR(devs.In1, ev3.BinData)
	irFL = ev3.OpenByteR(devs.In2, ev3.BinData)
	irFR = ev3.OpenByteR(devs.In3, ev3.BinData)
	irR = ev3.OpenByteR(devs.In4, ev3.BinData)
}

func close() {
	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutB, ev3.CmdStop)
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	// Close motors
	motorL1.Close()
	motorL2.Close()
	motorR1.Close()
	motorR2.Close()

	// Close sensors
	irL.Close()
	irFL.Close()
	irFR.Close()
	irR.Close()
}

var lastMoveTicks int
var lastSpeedLeft int
var lastSpeedRight int

const accelPerTicks int = 5

func move(left int, right int, now int) {
	ticks := now - lastMoveTicks
	lastMoveTicks = now

	nextSpeedLeft := lastSpeedLeft
	nextSpeedRight := lastSpeedRight
	delta := ticks * accelPerTicks
	// delta := ticks * ticks * accelPerTicks

	if left > nextSpeedLeft {
		nextSpeedLeft += delta
		if nextSpeedLeft > left {
			nextSpeedLeft = left
		}
	} else if left < nextSpeedLeft {
		nextSpeedLeft -= delta
		if nextSpeedLeft < left {
			nextSpeedLeft = left
		}
	}
	if right > nextSpeedRight {
		nextSpeedRight += delta
		if nextSpeedRight > right {
			nextSpeedRight = right
		}
	} else if right < nextSpeedRight {
		nextSpeedRight -= delta
		if nextSpeedRight < right {
			nextSpeedRight = right
		}
	}
	lastSpeedLeft = nextSpeedLeft
	lastSpeedRight = nextSpeedRight

	motorL1.Value = -nextSpeedLeft / 10000
	motorL2.Value = nextSpeedLeft / 10000
	motorR1.Value = -nextSpeedRight / 10000
	motorR2.Value = -nextSpeedRight / 10000

	motorL1.Sync()
	motorL2.Sync()
	motorR1.Sync()
	motorR2.Sync()
}

func read() {
	irL.Sync()
	irFL.Sync()
	irFR.Sync()
	irR.Sync()
}

func durationToTicks(d time.Duration) int {
	return int(d / 1000)
}
func timespanAsTicks(start time.Time, end time.Time) int {
	return durationToTicks(end.Sub(start))
}
func currentTicks() int {
	return timespanAsTicks(initializationTime, time.Now())
}
func ticksToMillis(ticks int) int {
	return ticks / 1000
}

func print(data ...interface{}) {
	fmt.Fprintln(os.Stderr, data)
}

func main() {
	initialize()
	defer close()

	/*
		start := currentTicks()
		for {
			now := currentTicks()
			if ticksToMillis(now-start) > 100 {
				break
			}
			read()
			move(0, 0, now)
		}
		start = currentTicks()
		for {
			now := currentTicks()
			if ticksToMillis(now-start) > 1000 {
				break
			}
			read()
			move(config.MaxSpeed, config.MaxSpeed, now)
	*/
	track(ev3.Right)
}

func track(dir ev3.Direction) {
	for {
		now := currentTicks()
		read()
		print(irL.Value, irFL.Value, irFR.Value, irR.Value)

		if irL.Value < config.MaxIrValue {
			move(-config.TrackTurnSpeed, config.TrackTurnSpeed, now)
			dir = ev3.Left
			print("LEFT")
		} else if irR.Value < config.MaxIrValue {
			move(config.TrackTurnSpeed, -config.TrackTurnSpeed, now)
			dir = ev3.Right
			print("RIGHT")
		} else if irFL.Value < config.MaxIrValue {
			move(config.TrackSpeed, config.TrackSpeed, now)
			dir = ev3.Left
			print("FRONT LEFT")
		} else if irFR.Value < config.MaxIrValue {
			move(config.TrackSpeed, config.TrackSpeed, now)
			dir = ev3.Right
			print("FRONT RIGHT")
		} else {
			if dir == ev3.Right {
				move(config.SeekTurnSpeed, -config.SeekTurnSpeed, now)
				print("SEEK RIGHT")
			} else if dir == ev3.Left {
				move(-config.SeekTurnSpeed, config.SeekTurnSpeed, now)
				print("SEEK LEFT")
			}
			print("SEEK NONE")
		}
	}
}
