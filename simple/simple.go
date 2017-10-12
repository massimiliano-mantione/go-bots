package main

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/simple/config"
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
var irL, irFL, irFR, irR *ev3.Attribute
var irRemote1, irRemote2, irRemote3, irRemote4 *ev3.Attribute

var conf config.Config

func closeIrProx() {
	if irL != nil {
		irL.Close()
		irL = nil
	}
	if irFL != nil {
		irFL.Close()
		irFL = nil
	}
	if irFR != nil {
		irFR.Close()
		irFR = nil
	}
	if irR != nil {
		irR.Close()
		irR = nil
	}
}

func closeIrRemote() {
	if irRemote1 != nil {
		irRemote1.Close()
		irRemote1 = nil
	}
	if irRemote2 != nil {
		irRemote2.Close()
		irRemote2 = nil
	}
	if irRemote3 != nil {
		irRemote3.Close()
		irRemote3 = nil
	}
	if irRemote4 != nil {
		irRemote4.Close()
		irRemote4 = nil
	}
}

func setIrProxMode() {
	closeIrRemote()

	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.IrModeProx)
	ev3.SetMode(devs.In3, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.IrModeProx)

	irL = ev3.OpenByteR(devs.In1, ev3.BinData)
	irFL = ev3.OpenByteR(devs.In2, ev3.BinData)
	irFR = ev3.OpenByteR(devs.In3, ev3.BinData)
	irR = ev3.OpenByteR(devs.In4, ev3.BinData)
}

func setIrRemoteMode(remoteChannel int) {
	closeIrProx()

	ev3.SetMode(devs.In1, ev3.IrModeRemote)
	ev3.SetMode(devs.In2, ev3.IrModeRemote)
	ev3.SetMode(devs.In3, ev3.IrModeRemote)
	ev3.SetMode(devs.In4, ev3.IrModeRemote)

	if remoteChannel == 1 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value0)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value0)
		irRemote3 = ev3.OpenTextR(devs.In3, ev3.Value0)
		irRemote4 = ev3.OpenTextR(devs.In4, ev3.Value0)
	} else if remoteChannel == 2 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value1)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value1)
		irRemote3 = ev3.OpenTextR(devs.In3, ev3.Value1)
		irRemote4 = ev3.OpenTextR(devs.In4, ev3.Value1)
	} else if remoteChannel == 3 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value2)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value2)
		irRemote3 = ev3.OpenTextR(devs.In3, ev3.Value2)
		irRemote4 = ev3.OpenTextR(devs.In4, ev3.Value2)
	} else if remoteChannel == 4 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value3)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value3)
		irRemote3 = ev3.OpenTextR(devs.In3, ev3.Value3)
		irRemote4 = ev3.OpenTextR(devs.In4, ev3.Value3)
	} else {
		quit("Invalid remote channel number", remoteChannel)
	}
}

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

	// Set sensors mode
	setIrProxMode()

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

	// Close sensor values
	closeIrProx()
	closeIrRemote()
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

func readRemote() int {
	irRemote1.Sync()
	irRemote2.Sync()
	irRemote3.Sync()
	irRemote4.Sync()

	result := 0
	if irRemote1.Value != 0 {
		result = irRemote1.Value
	} else if irRemote2.Value != 0 {
		result = irRemote2.Value
	} else if irRemote3.Value != 0 {
		result = irRemote3.Value
	} else if irRemote4.Value != 0 {
		result = irRemote4.Value
	}
	return result
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

func track(dir ev3.Direction) {
	for {
		now := currentTicks()
		read()
		print(irL.Value, irFL.Value, irFR.Value, irR.Value)

		if irL.Value < conf.MaxIrValue {
			move(-conf.TrackTurnSpeed, conf.TrackTurnSpeed, now)
			dir = ev3.Left
			print("LEFT")
		} else if irR.Value < conf.MaxIrValue {
			move(conf.TrackTurnSpeed, -conf.TrackTurnSpeed, now)
			dir = ev3.Right
			print("RIGHT")
		} else if irFL.Value < conf.MaxIrValue {
			move(conf.TrackSpeed, conf.TrackSpeed, now)
			dir = ev3.Left
			print("FRONT LEFT")
		} else if irFR.Value < conf.MaxIrValue {
			move(conf.TrackSpeed, conf.TrackSpeed, now)
			dir = ev3.Right
			print("FRONT RIGHT")
		} else {
			if dir == ev3.Right {
				move(conf.SeekTurnSpeed, -conf.SeekTurnSpeed, now)
				print("SEEK RIGHT")
			} else if dir == ev3.Left {
				move(-conf.SeekTurnSpeed, conf.SeekTurnSpeed, now)
				print("SEEK LEFT")
			}
			print("SEEK NONE")
		}
	}
}

func testRemote() {
	setIrRemoteMode(1)

	for {
		now := currentTicks()

		move(0, 0, now)
		rem := readRemote()
		if rem != 0 {
			print("received", rem)
		}
	}
}

func main() {
	handleSignals()
	initialize()
	defer close()

	conf = config.Default()

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
	testRemote()
	// track(ev3.Right)
}
