package main

import (
	"fmt"
	"go-bots/beep"
	"go-bots/earl_grey/config"
	"go-bots/ev3"
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
var motorL, motorR, motorBL, motorBR *ev3.Attribute
var pmotorFU *ev3.Attribute
var irL, irF, irR *ev3.Attribute
var irRemote1, irRemote2, irRemote3 *ev3.Attribute
var buttons *ev3.Buttons

var conf config.Config

func closeIrProx() {
	if irL != nil {
		irL.Close()
		irL = nil
	}
	if irF != nil {
		irF.Close()
		irF = nil
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
}

func setIrProxMode() {
	closeIrRemote()

	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ev3.SetMode(devs.In2, ev3.IrModeProx)
	ev3.SetMode(devs.In4, ev3.IrModeProx)

	irL = ev3.OpenByteR(devs.In1, ev3.BinData)
	irF = ev3.OpenByteR(devs.In2, ev3.BinData)
	irR = ev3.OpenByteR(devs.In4, ev3.BinData)
}

func setIrRemoteMode(remoteChannel int) {
	closeIrProx()

	ev3.SetMode(devs.In1, ev3.IrModeRemote)
	ev3.SetMode(devs.In2, ev3.IrModeRemote)
	ev3.SetMode(devs.In4, ev3.IrModeRemote)

	if remoteChannel == 1 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value0)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value0)
		irRemote3 = ev3.OpenTextR(devs.In4, ev3.Value0)
	} else if remoteChannel == 2 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value1)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value1)
		irRemote3 = ev3.OpenTextR(devs.In4, ev3.Value1)
	} else if remoteChannel == 3 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value2)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value2)
		irRemote3 = ev3.OpenTextR(devs.In4, ev3.Value2)
	} else if remoteChannel == 4 {
		irRemote1 = ev3.OpenTextR(devs.In1, ev3.Value3)
		irRemote2 = ev3.OpenTextR(devs.In2, ev3.Value3)
		irRemote3 = ev3.OpenTextR(devs.In4, ev3.Value3)
	} else {
		quit("Invalid remote channel number", remoteChannel)
	}
}

func initialize() {
	initializationTime = time.Now()

	buttons = ev3.OpenButtons(false)

	devs = ev3.Scan(&ev3.OutPortModes{
		OutA: ev3.OutPortModeAuto,
		OutB: ev3.OutPortModeAuto,
		OutC: ev3.OutPortModeAuto,
		OutD: ev3.OutPortModeAuto,
	})

	// Check motors
	ev3.CheckDriver(devs.OutA, ev3.DriverTachoMotorLarge, ev3.OutA)
	ev3.CheckDriver(devs.OutB, ev3.DriverTachoMotorLarge, ev3.OutB)
	ev3.CheckDriver(devs.OutC, ev3.DriverTachoMotorLarge, ev3.OutC)
	ev3.CheckDriver(devs.OutD, ev3.DriverTachoMotorLarge, ev3.OutD)

	// Check sensors
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	ev3.CheckDriver(devs.In2, ev3.DriverIr, ev3.In2)
	ev3.CheckDriver(devs.In4, ev3.DriverIr, ev3.In4)

	// Set sensors mode
	setIrProxMode()

	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)

	// Open motors
	motorL = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
	motorR = ev3.OpenTextW(devs.OutD, ev3.DutyCycleSp)
	motorBL = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	motorBR = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)

	pmotorFU = ev3.OpenTextR(devs.OutD, ev3.Position)

	// Reset motor speed
	motorL.Value = 0
	motorR.Value = 0
	motorBL.Value = 0
	motorBR.Value = 0

	motorL.Sync()
	motorR.Sync()
	motorBL.Sync()
	motorBR.Sync()

	// Put motors in direct mode
	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(devs.OutD, ev3.CmdRunDirect)
}

func close() {
	beep.CCC()

	// Close buttons
	buttons.Close()

	// Stop motors
	ev3.RunCommand(devs.OutA, ev3.CmdReset)
	ev3.RunCommand(devs.OutB, ev3.CmdReset)
	ev3.RunCommand(devs.OutC, ev3.CmdReset)
	ev3.RunCommand(devs.OutD, ev3.CmdReset)

	// Close motors
	motorL.Close()
	motorR.Close()
	motorBL.Close()
	motorBR.Close()

	// Close sensor values
	closeIrProx()
	closeIrRemote()
}

const accelPerTicks int = 5

func moveStop() {
	moveFull(0, 0, 0)
}

func move(left int, right int) {
	moveFull(left, right, 1)
}

var lastMoveTicks int
var lastSpeedLeft int
var lastSpeedRight int

const accelSpeedFactor int = 10000

func moveFull(left int, right int, useBack int) {
	now := currentTicks()
	ticks := now - lastMoveTicks
	lastMoveTicks = now

	right *= accelSpeedFactor
	left *= accelSpeedFactor

	nextSpeedLeft := lastSpeedLeft
	nextSpeedRight := lastSpeedRight
	delta := ticks * 30

	if left > nextSpeedLeft {
		nextSpeedLeft += delta
		if nextSpeedLeft > left {
			nextSpeedLeft = left
		}
	} else if left < nextSpeedLeft {
		nextSpeedLeft = left
	}
	if right > nextSpeedRight {
		nextSpeedRight += delta
		if nextSpeedRight > right {
			nextSpeedRight = right
		}
	} else if right < nextSpeedRight {
		nextSpeedRight = right
	}
	lastSpeedLeft = nextSpeedLeft
	lastSpeedRight = nextSpeedRight

	motorL.Value = nextSpeedLeft / accelSpeedFactor
	motorR.Value = nextSpeedRight / accelSpeedFactor

	if useBack == -1 {
		motorBL.Value = -100
		motorBR.Value = -100
	} else if useBack == 1 {
		motorBL.Value = 100
		motorBR.Value = 100
	} else {
		motorBL.Value = 0
		motorBR.Value = 0
	}

	motorL.Sync()
	motorR.Sync()
	motorBL.Sync()
	motorBR.Sync()
}

func read() {
	irL.Sync()
	irF.Sync()
	irR.Sync()
}

func readRemote() int {
	irRemote1.Sync()
	irRemote2.Sync()
	irRemote3.Sync()

	result := 0
	if irRemote1.Value != 0 {
		result = irRemote1.Value
	} else if irRemote2.Value != 0 {
		result = irRemote2.Value
	} else if irRemote3.Value != 0 {
		result = irRemote3.Value
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
	if irL.Value < conf.MaxIrSide || irF.Value < conf.MaxIrFront || irR.Value < conf.MaxIrSide {
		return true
	}
	return false
}

func loadConfig() {
	newConf, err := config.FromFile("earl_grey.toml")
	if err != nil {
		print("Error reading conf:", err)
		conf = config.Default()
		print("Using default conf", conf)
	} else {
		conf = newConf
		print("Configuration loaded:", conf)
	}
}

var strategyDirection = 0

var channelNumber int

func chooseStrategy() bool {
	if channelNumber == 0 {
		channelNumber = 2
	}
	ev3.WriteStringAttribute(devs.OutB, ev3.Position, "0")
	setIrRemoteMode(channelNumber)
	for {
		moveStop()
		remoteValue := readRemote()

		if remoteValue == 11 {
			return true
		} else if remoteValue == 3 {
			loadConfig()
			ev3.WriteStringAttribute(devs.OutB, ev3.Position, "0")
			beep.GC()
			strategyDirection = -1
		} else if remoteValue == 4 {
			loadConfig()
			ev3.WriteStringAttribute(devs.OutB, ev3.Position, "0")
			beep.CG()
			strategyDirection = 1
		} else if remoteValue == 2 {
			loadConfig()
			ev3.WriteStringAttribute(devs.OutB, ev3.Position, "0")
			beep.GG()
			strategyDirection = 0
		} else if remoteValue == 1 {
			beep.C()
			return false
		}

		if buttons.Up {
			channelNumber = 1
			setIrRemoteMode(channelNumber)
			beep.C()
		} else if buttons.Right {
			channelNumber = 2
			setIrRemoteMode(channelNumber)
			beep.CC()
		} else if buttons.Down {
			channelNumber = 3
			setIrRemoteMode(channelNumber)
			beep.CCC()
		} else if buttons.Left {
			channelNumber = 4
			setIrRemoteMode(channelNumber)
			beep.CCCC()
		}
	}
}

func waitBegin() {
	print("wait 5 seconds")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		moveStop()
		if elapsed >= conf.WaitTime {
			return
		}
	}
}

func strategy() ev3.Direction {
	setIrProxMode()
	print("strategy")

	if strategyDirection == -1 {
		strategyTurn(-1)
		return 1
	} else if strategyDirection == 1 {
		strategyTurn(1)
		return -1
	} else {
		strategyStraight()
		return -1
	}
}

func moveStrategyFull(dir ev3.Direction, duration int, useBack int, useVision bool) bool {
	start := currentTicks()
	for {
		now := currentTicks()
		if now-start >= duration {
			break
		}
		if useVision && checkVision() {
			return true
		}
		if dir == ev3.Left {
			moveFull(-10, conf.MaxSpeed, useBack)
		} else if dir == ev3.Right {
			moveFull(conf.MaxSpeed, -10, useBack)
		} else {
			moveFull(conf.MaxSpeed, conf.MaxSpeed, useBack)
		}
	}
	return false
}

func moveStrategy(dir ev3.Direction, duration int) bool {
	return moveStrategyFull(dir, duration, 0, true)
}

func strategyTurn(dir ev3.Direction) {

	if moveStrategyFull(dir, conf.StrategyR1Time, 0, false) {
		return
	}
	if moveStrategyFull(0, 100000, -1, false) {
		return
	}

	if moveStrategy(0, conf.StrategyS1Time) {
		return
	}
	if moveStrategy(-dir, conf.StrategyR2Time) {
		return
	}
	if moveStrategy(0, conf.StrategyS2Time) {
		return
	}
}

func strategyStraight() {

	if moveStrategyFull(0, conf.StrategyR1Time, 0, false) {
		return
	}
	if moveStrategyFull(0, 100000, -1, false) {
		return
	}

	if moveStrategy(0, conf.StrategyStraightTime) {
		return
	}
}

func track(dir ev3.Direction) {
	print("track", irL.Value, irF.Value, irR.Value)
	for {
		if buttons.Up || buttons.Down || buttons.Left || buttons.Right || buttons.Back || buttons.Enter {
			return
		}
		read()
		if irF.Value < conf.MaxIrFront {
			move(conf.TrackSpeed, conf.TrackSpeed)
			print("FRONT")
		} else if irL.Value < conf.MaxIrSide {
			move(-conf.TrackTurnSpeed, conf.TrackTurnSpeed)
			dir = ev3.Left
			print("LEFT")
		} else if irR.Value < conf.MaxIrSide {
			move(conf.TrackTurnSpeed, -conf.TrackTurnSpeed)
			dir = ev3.Right
			print("RIGHT")
		} else {
			if dir == ev3.Right {
				move(conf.SeekTurnSpeed, -conf.SeekTurnSpeed)
				print("SEEK RIGHT")
			} else if dir == ev3.Left {
				move(-conf.SeekTurnSpeed, conf.SeekTurnSpeed)
				print("SEEK LEFT")
			} else {
				print("SEEK NONE")
			}
		}
	}
}

func main() {
	handleSignals()
	initialize()
	defer close()

	beep.G()

	// conf = config.Default()

	for {
		if chooseStrategy() {
			return
		}
		waitBegin()
		trackDir := strategy()
		track(trackDir)
	}
}
