package main

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/super_red/config"
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
var motorL, motorR, motorFU, motorFD *ev3.Attribute
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
		OutA: ev3.OutPortModeAuto,
		OutB: ev3.OutPortModeAuto,
		OutC: ev3.OutPortModeAuto,
		OutD: ev3.OutPortModeAuto,
	})

	// Check motors
	ev3.CheckDriver(devs.OutA, ev3.DriverTachoMotorMedium, ev3.OutA)
	ev3.CheckDriver(devs.OutB, ev3.DriverTachoMotorLarge, ev3.OutB)
	ev3.CheckDriver(devs.OutC, ev3.DriverTachoMotorLarge, ev3.OutC)
	ev3.CheckDriver(devs.OutD, ev3.DriverTachoMotorLarge, ev3.OutD)

	// Check sensors
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	ev3.CheckDriver(devs.In2, ev3.DriverIr, ev3.In2)
	ev3.CheckDriver(devs.In3, ev3.DriverIr, ev3.In3)
	ev3.CheckDriver(devs.In4, ev3.DriverIr, ev3.In4)

	// Set sensors mode
	setIrProxMode()

	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)

	// Open motors
	motorL = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	motorR = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	motorFU = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	motorFD = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)

	// Reset motor speed
	motorL.Value = 0
	motorR.Value = 0
	motorFU.Value = 0
	motorFD.Value = 0

	motorL.Sync()
	motorR.Sync()
	motorFU.Sync()
	motorFD.Sync()

	// Put motors in direct mode
	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)
}

func close() {
	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)

	// Close motors
	motorL.Close()
	motorR.Close()
	motorFU.Close()
	motorFD.Close()

	// Close sensor values
	closeIrProx()
	closeIrRemote()
}

var lastMoveTicks int
var lastSpeedLeft int
var lastSpeedRight int

const accelPerTicks int = 5

func move(left int, right int, frontUp int, frontDown int) {

	motorL.Value = -left / 10000
	motorR.Value = -right / 10000
	motorFU.Value = -frontUp / 10000
	motorFD.Value = frontDown / 10000

	// motorL1.Value = 0
	// motorL2.Value = 0
	// motorR1.Value = 0
	// motorR2.Value = 0

	motorL.Sync()
	motorR.Sync()
	motorFU.Sync()
	motorFD.Sync()
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

func checkVision() bool {
	read()
	if irL.Value < conf.MaxIrValue || irFL.Value < conf.MaxIrValue || irFR.Value < conf.MaxIrValue || irR.Value < conf.MaxIrValue {
		return true
	}
	return false
}

var strategyDirection = 0

func chooseStrategy(channelNumber int) {
	setIrRemoteMode(channelNumber)
	for {
		move(0, 0, 0, 0)
		remoteValue := readRemote()

		if remoteValue == 3 {
			strategyDirection = -1
		} else if remoteValue == 4 {
			strategyDirection = 1
		} else if remoteValue == 2 {
			strategyDirection = 0
		} else if remoteValue == 1 {
			waitBegin()
			return
		}
		print(strategyDirection)
	}
}

func waitBegin() {
	print("wait 5 seconds")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(0, 0, 0, 0)
		if elapsed >= 5000000 {
			strategy()
			return
		}
	}
}

func strategy() {
	setIrProxMode()
	print("strategy")
	for {
		move(0, 0, 0, 0)

		if strategyDirection == -1 {
			strategyLeft()
		} else if strategyDirection == 0 {
			strategyStraight()
		} else if strategyDirection == 1 {
			strategyRight()
		}
	}
}

func strategyLeft() {
	print("strategy left")
	startR1 := currentTicks()
	for {
		now := currentTicks()
		if now-startR1 >= conf.StrategyR1Time {
			break
		}
		if checkVision() {
			track(ev3.Right)
			return
		}
		move(-20, conf.MaxSpeed, 0, 0)
	}
	startS1 := currentTicks()
	for {
		now := currentTicks()
		if now-startS1 >= conf.StrategyS1Time {
			break
		}
		if checkVision() {
			track(ev3.Right)
			return
		}
		move(conf.MaxSpeed, conf.MaxSpeed, 0, 0)
	}
	startR2 := currentTicks()
	for {
		now := currentTicks()
		if now-startR2 >= conf.StrategyR2Time {
			break
		}
		if checkVision() {
			track(ev3.Right)
			return
		}
		move(conf.MaxSpeed, -20, 0, 0)
	}
	startS2 := currentTicks()
	for {
		now := currentTicks()
		if now-startS2 >= conf.StrategyS2Time {
			break
		}
		if checkVision() {
			track(ev3.Right)
			return
		}
		move(conf.MaxSpeed, conf.MaxSpeed, 0, 0)
	}
	track(ev3.Right)
	return
}

func strategyStraight() {
	print("strategy straight")
	start := currentTicks()
	for {
		now := currentTicks()
		if now-start >= conf.StrategyStraightTime {
			break
		}
		if checkVision() {
			track(ev3.Left)
			return
		}
		move(conf.MaxSpeed, conf.MaxSpeed, 0, 0)
	}
	track(ev3.Left)
	return
}

func strategyRight() {
	print("strategy right")
	startR1 := currentTicks()
	for {
		now := currentTicks()
		if now-startR1 >= conf.StrategyR1Time {
			break
		}
		if checkVision() {
			track(ev3.Left)
			return
		}
		move(conf.MaxSpeed, -20, 0, 0)
	}
	startS1 := currentTicks()
	for {
		now := currentTicks()
		if now-startS1 >= conf.StrategyS1Time {
			break
		}
		if checkVision() {
			track(ev3.Left)
			return
		}
		move(conf.MaxSpeed, conf.MaxSpeed, 0, 0)
	}
	startR2 := currentTicks()
	for {
		now := currentTicks()
		if now-startR2 >= conf.StrategyR2Time {
			break
		}
		if checkVision() {
			track(ev3.Left)
			return
		}
		move(-20, conf.MaxSpeed, 0, 0)
	}
	startS2 := currentTicks()
	for {
		now := currentTicks()
		if now-startS2 >= conf.StrategyS2Time {
			break
		}
		if checkVision() {
			track(ev3.Left)
			return
		}
		move(conf.MaxSpeed, conf.MaxSpeed, 0, 0)
	}
	track(ev3.Left)
	return
}

func track(dir ev3.Direction) {
	print("track", irL.Value, irFL.Value, irFR.Value, irR.Value)
	for {
		read()
		print(irL.Value, irFL.Value, irFR.Value, irR.Value)

		if irL.Value < conf.MaxIrValue {
			move(-conf.TrackTurnSpeed, conf.TrackTurnSpeed, conf.MaxSpeed, conf.MaxSpeed)
			dir = ev3.Left
			print("LEFT")
		} else if irR.Value < conf.MaxIrValue {
			move(conf.TrackTurnSpeed, -conf.TrackTurnSpeed, conf.MaxSpeed, conf.MaxSpeed)
			dir = ev3.Right
			print("RIGHT")
		} else if irFL.Value < conf.MaxIrValue {
			move(conf.TrackSpeed, conf.TrackSpeed, conf.MaxSpeed, conf.MaxSpeed)
			dir = ev3.Left
			print("FRONT LEFT")
		} else if irFR.Value < conf.MaxIrValue {
			move(conf.TrackSpeed, conf.TrackSpeed, conf.MaxSpeed, conf.MaxSpeed)
			dir = ev3.Right
			print("FRONT RIGHT")
		} else {
			if dir == ev3.Right {
				move(conf.SeekTurnSpeed, -conf.SeekTurnSpeed, 0, 0)
				print("SEEK RIGHT")
			} else if dir == ev3.Left {
				move(-conf.SeekTurnSpeed, conf.SeekTurnSpeed, 0, 0)
				print("SEEK LEFT")
			} else {
				print("SEEK NONE")
			}
		}
	}
}

func testRemote() {
	setIrRemoteMode(1)

	for {
		move(0, 0, 0, 0)
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
	chooseStrategy(2)
}
