package main

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/greyhound/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		quit("Terminated by signal", sig)
	}()
}

var devs *ev3.Devices
var initializationTime time.Time
var motorL1, motorL2, motorR1, motorR2 *ev3.Attribute
var cLL, cL, cR, cRR *ev3.Attribute
var buttons *ev3.Buttons

var conf config.Config

func closeSensors() {
	cLL.Close()
	cL.Close()
	cR.Close()
	cRR.Close()
}

func setSensorsMode() {
	ev3.SetMode(devs.In1, ev3.ColorModeReflect)
	ev3.SetMode(devs.In2, ev3.ColorModeReflect)
	ev3.SetMode(devs.In3, ev3.ColorModeReflect)
	ev3.SetMode(devs.In4, ev3.ColorModeReflect)

	cLL = ev3.OpenByteR(devs.In4, ev3.BinData)
	cL = ev3.OpenByteR(devs.In3, ev3.BinData)
	cR = ev3.OpenByteR(devs.In2, ev3.BinData)
	cRR = ev3.OpenByteR(devs.In1, ev3.BinData)
}

func initialize() {
	initializationTime = time.Now()

	buttons = ev3.OpenButtons()

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
	ev3.CheckDriver(devs.In1, ev3.DriverColor, ev3.In1)
	ev3.CheckDriver(devs.In2, ev3.DriverColor, ev3.In2)
	ev3.CheckDriver(devs.In3, ev3.DriverColor, ev3.In3)
	ev3.CheckDriver(devs.In4, ev3.DriverColor, ev3.In4)

	// Set sensors mode
	setSensorsMode()

	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutB, ev3.CmdStop)
	ev3.RunCommand(devs.OutC, ev3.CmdStop)
	ev3.RunCommand(devs.OutD, ev3.CmdStop)

	// Open motors
	motorL1 = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	motorL2 = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	motorR1 = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	motorR2 = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)

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
}

func close() {
	// Close buttons
	buttons.Close()

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

	// Close sensor values
	closeSensors()
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

	motorL1.Value = nextSpeedLeft / 10000
	motorL2.Value = nextSpeedLeft / 10000
	motorR1.Value = -nextSpeedRight / 10000
	motorR2.Value = -nextSpeedRight / 10000

	motorL1.Sync()
	motorL2.Sync()
	motorR1.Sync()
	motorR2.Sync()
}

func read() {
	cLL.Sync()
	cL.Sync()
	cR.Sync()
	cRR.Sync()
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
	fmt.Fprintln(os.Stderr, data...)
}

func quit(data ...interface{}) {
	close()
	log.Fatalln(data...)
}

func waitOneSecond() {
	print("wait one second")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(0, 0, now)
		if elapsed >= 1000000 {
			break
		}
	}
}

func moveOneSecond() {
	print("move one second")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(conf.MaxSpeed, conf.MaxSpeed, now)
		if elapsed >= 1000000 {
			break
		}
	}
}

func main() {
	handleSignals()
	initialize()
	defer close()

	conf = config.Default()

	lastPrint := 0
	for {
		now := currentTicks()
		if now-lastPrint > 1000000 {
			lastPrint = now
			print("Tock...")
		}
	}

	// waitOneSecond()
	// moveOneSecond()
}
