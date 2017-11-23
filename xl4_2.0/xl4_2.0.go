package main

import (
	"fmt"
	"go-bots/ev3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-bots/beep"

	"go-bots/xl4_2.0/config"
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
var buttons *ev3.Buttons

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

	buttons = ev3.OpenButtons(false)

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
	motorL1 = ev3.OpenTextW(devs.OutB, ev3.DutyCycleSp)
	motorL2 = ev3.OpenTextW(devs.OutC, ev3.DutyCycleSp)
	motorR1 = ev3.OpenTextW(devs.OutA, ev3.DutyCycleSp)
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
	beep.CCC()

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
	closeIrProx()
	closeIrRemote()
}

var lastMoveTicks int
var lastSpeedLeft int
var lastSpeedRight int

const accelSpeedFactor int = 10000

func move(left int, right int, now int) {
	ticks := now - lastMoveTicks
	lastMoveTicks = now
	right *= accelSpeedFactor
	left *= accelSpeedFactor

	nextSpeedLeft := lastSpeedLeft
	nextSpeedRight := lastSpeedRight
	delta := ticks * conf.AccelPerTicks
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

	motorL1.Value = nextSpeedLeft / accelSpeedFactor
	motorL2.Value = -nextSpeedLeft / accelSpeedFactor
	motorR1.Value = nextSpeedRight / accelSpeedFactor
	motorR2.Value = -nextSpeedRight / accelSpeedFactor

	// motorL1.Value = 0
	// motorL2.Value = 0
	// motorR1.Value = 0
	// motorR2.Value = 0

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

func checkVision() bool {
	read()
	if irL.Value < conf.MaxIrSide || irFL.Value < conf.MaxIrFront || irFR.Value < conf.MaxIrFront || irR.Value < conf.MaxIrSide {
		print("checkvision true", irL.Value, irFL.Value, irFR.Value, irR.Value)
		return true
	}
	return false
}

func loadConfig() {
	newConf, err := config.FromFile("xl4_2.0.toml")
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
		channelNumber = 1
	}
	setIrRemoteMode(channelNumber)
	for {
		now := currentTicks()
		move(0, 0, now)
		remoteValue := readRemote()

		if remoteValue == 11 {
			return true
		} else if remoteValue == 3 {
			loadConfig()
			beep.GC()
			strategyDirection = -1
		} else if remoteValue == 4 {
			loadConfig()
			beep.CG()
			strategyDirection = 1
		} else if remoteValue == 2 {
			loadConfig()
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
		move(0, 0, now)
		if elapsed >= conf.WaitTime {
			return
		}
	}
}

func strategy() ev3.Direction {
	setIrProxMode()
	print("strategy")
	for {
		now := currentTicks()
		move(0, 0, now)

		if strategyDirection == -1 {
			return strategyLeft()
		} else if strategyDirection == 0 {
			return strategyStraight()
		} else if strategyDirection == 1 {
			return strategyRight()
		}
	}
}

func strategyLeft() ev3.Direction {
	print("strategy left")
	startR1 := currentTicks()
	for {
		now := currentTicks()
		if now-startR1 >= conf.StrategyR1Time {
			break
		}
		if checkVision() {
			return ev3.Right
		}
		move(-20, conf.MaxSpeed, now)
	}
	startS1 := currentTicks()
	for {
		now := currentTicks()
		if now-startS1 >= conf.StrategyS1Time {
			break
		}
		if checkVision() {
			return ev3.Right
		}
		move(conf.MaxSpeed, conf.MaxSpeed, now)
	}
	startR2 := currentTicks()
	for {
		now := currentTicks()
		if now-startR2 >= conf.StrategyR2Time {
			break
		}
		if checkVision() {
			return ev3.Right
		}
		move(conf.MaxSpeed, -20, now)
	}
	startS2 := currentTicks()
	for {
		now := currentTicks()
		if now-startS2 >= conf.StrategyS2Time {
			break
		}
		if checkVision() {
			return ev3.Right
		}
		move(conf.MaxSpeed, conf.MaxSpeed, now)
	}
	print("ho finito strategy left, seeeee!!!")
	return ev3.Right
}

func strategyStraight() ev3.Direction {
	print("strategy straight")
	start := currentTicks()
	for {
		now := currentTicks()
		if now-start >= conf.StrategyStraightTime {
			break
		}
		if checkVision() {
			return ev3.Left
		}
		move(conf.MaxSpeed, conf.MaxSpeed, now)
	}
	print("ho finito strategy straight, seeeee!!!")
	return ev3.Left
}

func strategyRight() ev3.Direction {
	print("strategy right")
	startR1 := currentTicks()
	for {
		now := currentTicks()
		if now-startR1 >= conf.StrategyR1Time {
			break
		}
		if checkVision() {
			return ev3.Left
		}
		move(conf.MaxSpeed, -20, now)
	}
	startS1 := currentTicks()
	for {
		now := currentTicks()
		if now-startS1 >= conf.StrategyS1Time {
			break
		}
		if checkVision() {
			return ev3.Left
		}
		move(conf.MaxSpeed, conf.MaxSpeed, now)
	}
	startR2 := currentTicks()
	for {
		now := currentTicks()
		if now-startR2 >= conf.StrategyR2Time {
			break
		}
		if checkVision() {
			return ev3.Left
		}
		move(-20, conf.MaxSpeed, now)
	}
	startS2 := currentTicks()
	for {
		now := currentTicks()
		if now-startS2 >= conf.StrategyS2Time {
			break
		}
		if checkVision() {
			return ev3.Left
		}
		move(conf.MaxSpeed, conf.MaxSpeed, now)
	}
	print("ho finito strategy right, seeeee!!!")
	return ev3.Left
}

func track(dir ev3.Direction) {
	print("track", irL.Value, irFL.Value, irFR.Value, irR.Value)
	for {
		if buttons.Up || buttons.Down || buttons.Left || buttons.Right || buttons.Back || buttons.Enter {
			return
		}
		now := currentTicks()
		read()
		// print(irL.Value, irFL.Value, irFR.Value, irR.Value)

		if irFL.Value < conf.MaxIrFront {
			move(conf.TrackSpeed, conf.TrackSpeed, now)
			dir = ev3.Left
			print("FRONT LEFT")
		} else if irFR.Value < conf.MaxIrFront {
			move(conf.TrackSpeed, conf.TrackSpeed, now)
			dir = ev3.Right
			print("FRONT RIGHT")
		} else if irL.Value < conf.MaxIrSide {
			move(-conf.TrackTurnSpeed, conf.TrackTurnSpeed, now)
			dir = ev3.Left
			print("LEFT")
		} else if irR.Value < conf.MaxIrSide {
			move(conf.TrackTurnSpeed, -conf.TrackTurnSpeed, now)
			dir = ev3.Right
			print("RIGHT")
		} else {
			if dir == ev3.Right {
				move(conf.SeekTurnSpeed, -conf.SeekTurnSpeed, now)
				print("SEEK RIGHT")
			} else if dir == ev3.Left {
				move(-conf.SeekTurnSpeed, conf.SeekTurnSpeed, now)
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
