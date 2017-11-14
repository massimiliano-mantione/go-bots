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
var cF, cL, cR, cB *ev3.Attribute
var buttons *ev3.Buttons

var conf config.Config

func closeSensors() {
	cF.Close()
	cL.Close()
	cR.Close()
	cB.Close()
}

func setSensorsMode() {
	ev3.SetMode(devs.In1, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In2, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In3, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In4, ev3.ColorModeRgbRaw)

	cF = ev3.OpenBinaryR(devs.In1, ev3.BinData, 3, 2)
	cL = ev3.OpenBinaryR(devs.In2, ev3.BinData, 3, 2)
	cR = ev3.OpenBinaryR(devs.In3, ev3.BinData, 3, 2)
	cB = ev3.OpenBinaryR(devs.In4, ev3.BinData, 3, 2)
}

func initializeTime() {
	initializationTime = time.Now()
}

func initialize() {
	initializeTime()

	buttons = ev3.OpenButtons(true)

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
const accelSpeedFactor int = 10000

func move(left int, right int, now int) {
	ticks := now - lastMoveTicks
	lastMoveTicks = now
	right *= accelSpeedFactor
	left *= accelSpeedFactor

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

	motorL1.Value = nextSpeedLeft / accelSpeedFactor
	motorL2.Value = nextSpeedLeft / accelSpeedFactor
	motorR1.Value = -nextSpeedRight / accelSpeedFactor
	motorR2.Value = -nextSpeedRight / accelSpeedFactor

	motorL1.Sync()
	motorL2.Sync()
	motorR1.Sync()
	motorR2.Sync()
}

func read() {
	cF.Sync()
	cL.Sync()
	cR.Sync()
	cB.Sync()
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

func waitEnter() {
	// Let the button be released if needed
	if buttons.Enter == true {
		print("wait enter release")
		for buttons.Enter == true {
			now := currentTicks()
			move(0, 0, now)
		}
	}

	// Wait for it to be pressed
	print("wait enter")
	for buttons.Enter == false {
		now := currentTicks()
		move(0, 0, now)
		if buttons.Back {
			newConf, err := config.FromFile("greyhound.toml")
			if err != nil {
				print("Error reading conf:", err)
			} else {
				conf = newConf
				print("Configuration reloaded:", conf)
			}
		}
	}
}

func waitOneSecond() int {
	initializeTime()
	print("wait one second")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(0, 0, now)
		if buttons.Enter && buttons.Back {
			quit("Done")
		}
		if elapsed >= 1000000 {
			return now
		}
	}
}

func trimSensor(attr *ev3.Attribute) int {
	value := attr.Value + attr.Value1 + attr.Value2
	if value < conf.SensorMin {
		value = conf.SensorMin
	}
	value -= conf.SensorMin
	if value > conf.SensorSpan {
		value = conf.SensorSpan
	}
	return value
}
func isOnTrack(value int) bool {
	return value < conf.SensorSpan
}
func distanceFromSensor(value int) int {
	return value * conf.SensorRadius / conf.SensorSpan
}
func positionBetweenSensors(value1 int, value2 int) int {
	return (value1 - value2) * conf.SensorRadius / conf.SensorSpan
}

func sign(value int) int {
	if value > 0 {
		return 1
	} else if value < 0 {
		return -1
	} else {
		return 0
	}
}

type sensorReadType int

const (
	bitB sensorReadType = 1 << iota
	bitR
	bitL
	bitF
)

const (
	sensorReadZero sensorReadType = iota
	sensorReadB
	sensorReadR
	sensorReadRB
	sensorReadL
	sensorReadLB
	sensorReadLR
	sensorReadLRB
	sensorReadF
	sensorReadFB
	sensorReadFR
	sensorReadFRB
	sensorReadFL
	sensorReadFLB
	sensorReadFLR
	sensorReadFLRB
)

var sensorReadNames = [16]string{
	"---",
	"-v-",
	"-->",
	"-v>",
	"<--",
	"<v-",
	"<->",
	"<v>",
	"-^-",
	"-X-",
	"-^>",
	"-X>",
	"<^-",
	"<X-",
	"<^>",
	"<X>",
}

// const lineStatusStraight = "-|-"
// const lineStatusStraightLeft = "<|-"
// const lineStatusLeft = "<--"
// const lineStatusStraightRight = "-|>"
// const lineStatusRight = "-->"
// const lineStatusFrontLeft = "<^-"
// const lineStatusFrontRight = "-^>"
// const lineStatusBackLeft = "<v-"
// const lineStatusBackRight = "-v>"
// const lineStatusOut = "---"
// const lineStatusCross = "-+-"

func processSensorData() (sensorRead sensorReadType, pos int, hint int, cross bool, out bool) {
	read()
	f, l, r, b := trimSensor(cB), trimSensor(cL), trimSensor(cR), trimSensor(cB)

	sensorRead = sensorReadZero
	if isOnTrack(b) {
		sensorRead |= bitB
	}
	if isOnTrack(r) {
		sensorRead |= bitR
	}
	if isOnTrack(l) {
		sensorRead |= bitL
	}
	if isOnTrack(f) {
		sensorRead |= bitF
	}

	switch sensorRead {
	case sensorReadZero, sensorReadB, sensorReadF:
		// Out
		out = true
		pos, hint, cross = 0, 0, false
	case sensorReadR:
		pos = conf.SensorRadius*2 + distanceFromSensor(r)
		hint = 0
		cross, out = false, false
	case sensorReadRB:
		pos = conf.SensorRadius + positionBetweenSensors(b, r)
		hint = 1
		cross, out = false, false
	case sensorReadL:
		pos = -conf.SensorRadius*2 - distanceFromSensor(l)
		hint = 0
		cross, out = false, false
	case sensorReadLB:
		pos = -conf.SensorRadius + positionBetweenSensors(l, b)
		hint = 1
		cross, out = false, false
	case sensorReadLR, sensorReadLRB, sensorReadFLRB, sensorReadFLR:
		// Cross
		cross = true
		pos, hint, out = 0, 0, false
	case sensorReadFB:
		pos = 0
		hint = 0
		cross, out = false, false
	case sensorReadFR:
		pos = conf.SensorRadius + positionBetweenSensors(f, r)
		hint = -1
		cross, out = false, false
	case sensorReadFRB:
		pos = conf.SensorRadius + positionBetweenSensors((f+b)/2, r)
		hint = 0
		cross, out = false, false
	case sensorReadFL:
		pos = -conf.SensorRadius + positionBetweenSensors(l, f)
		hint = 1
		cross, out = false, false
	case sensorReadFLB:
		pos = -conf.SensorRadius + positionBetweenSensors(l, (f+b)/2)
		hint = 0
		cross, out = false, false
	default:
		print("Error: reading", sensorRead)
	}

	return
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

func followLine(lastGivenTicks int) {
	print("following line")

	lastTicks := lastGivenTicks
	lastPos := 0
	lastPosD := 0

	for {
		now := currentTicks()
		sr, pos, hint, cross, out := processSensorData()
		posD := 0

		if out {
			pos = conf.SensorRadius * 3 * sign(lastPos)
			hint = sign(pos)
			posD = lastPosD
		} else if cross {
			pos = lastPos
			posD = 0
		} else {
			dTicks := now - lastTicks
			if dTicks < 5 {
				dTicks = 5
			}
			if dTicks > 100000 {
				dTicks = 100000
			}
			posD = ((pos - lastPos) * 100000) / dTicks
			posD += hint
		}

		pos2 := sign(pos) * pos * pos
		posD2 := sign(posD) * posD * posD

		factorP := (pos * conf.KP * conf.MaxSpeed) / (conf.MaxPos * 100)
		factorP2 := (pos2 * conf.KP2 * conf.MaxSpeed) / (conf.MaxPos2 * 100)
		factorD := (posD * conf.KD * conf.MaxSpeed) / (conf.MaxPosD * 100)
		factorD2 := (posD2 * conf.KD2 * conf.MaxSpeed) / (conf.MaxPosD2 * 100)

		steering := factorP + factorP2 - factorD - factorD2

		print(sensorReadNames[sr], "pos", pos, "d", posD, "f", factorP, factorP2, factorD, factorD2, "t", (now-lastTicks)/1000, "s", steering)

		lastTicks, lastPos, lastPosD = now, pos, posD

		if steering > 0 {
			if steering > conf.MaxSteering {
				steering = conf.MaxSteering
			}
			move(conf.MaxSpeed, conf.MaxSpeed-steering, now)
		} else if steering < 0 {
			steering = -steering
			if steering > conf.MaxSteering {
				steering = conf.MaxSteering
			}
			move(conf.MaxSpeed-steering, conf.MaxSpeed, now)
		} else {
			move(conf.MaxSpeed, conf.MaxSpeed, now)
		}

		if buttons.Enter {
			print("stopping")
			break
		}
	}
}

func main() {
	handleSignals()
	initialize()
	defer close()

	newConf, err := config.FromFile("greyhound.toml")
	if err != nil {
		print("Error reading conf:", err)
		conf = config.Default()
		print("Using default conf", conf)
	} else {
		conf = newConf
		print("Configuration loaded:", conf)
	}

	waitEnter()
	lastGivenTicks := waitOneSecond()
	followLine(lastGivenTicks)
}
