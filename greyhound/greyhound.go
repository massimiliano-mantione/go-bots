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
	ev3.SetMode(devs.In1, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In2, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In3, ev3.ColorModeRgbRaw)
	ev3.SetMode(devs.In4, ev3.ColorModeRgbRaw)

	cLL = ev3.OpenBinaryR(devs.In1, ev3.BinData, 3, 2)
	cL = ev3.OpenBinaryR(devs.In2, ev3.BinData, 3, 2)
	cR = ev3.OpenBinaryR(devs.In3, ev3.BinData, 3, 2)
	cRR = ev3.OpenBinaryR(devs.In4, ev3.BinData, 3, 2)
}

func initialize() {
	initializationTime = time.Now()

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

func waitOneSecond() {
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
			break
		}
	}
}

func sumSensor(attr *ev3.Attribute) int {
	return attr.Value + attr.Value1 + attr.Value2
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
func isWhite(value int) bool {
	return value >= conf.SensorSpan
}
func distanceFromSensor(value int) int {
	return value * conf.SensorRadius / conf.SensorSpan
}
func positionBetweenSensors(value1 int, value2 int) int {
	return (value1 - value2) * conf.SensorRadius / conf.SensorSpan
}

func hintCenter() {
	nextPifLost = 0
}
func hintLeft() {
	nextPifLost = -conf.SensorRadius * 4
}
func hintRight() {
	nextPifLost = conf.SensorRadius * 4
}
func hintBetween(left int, right int) {
	if right > left {
		hintLeft()
	} else if left > right {
		hintRight()
	} else {
		hintCenter()
	}
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

var currentP int
var currentPTicks int
var lastP int
var lastPTicks int
var nextPifLost int

func processSensorData() {
	read()
	ll, l, r, rr := trimSensor(cLL), trimSensor(cL), trimSensor(cR), trimSensor(cRR)
	radius := conf.SensorRadius

	lastP = currentP
	lastPTicks = currentPTicks
	currentPTicks = currentTicks()

	status := "XXXX"

	if isWhite(ll) {
		if isWhite(l) {
			// left - -
			if isWhite(r) {
				if isWhite(rr) {
					// We lost the line, use last hint
					status = "----"
					currentP = nextPifLost
				} else {
					// Far right
					status = "---X"
					currentP = (radius * 3) + distanceFromSensor(rr)
					hintRight()
				}
			} else {
				if isWhite(rr) {
					// Centered on r
					status = "--X-"
					currentP = radius - distanceFromSensor(r)
					hintCenter()
				} else {
					// Angle (r, rr)
					status = "--XX"
					currentP = (radius * 2) + positionBetweenSensors(r, rr)
					hintRight()
				}
			}
		} else {
			// left - l
			if isWhite(r) {
				if isWhite(rr) {
					// Centered on l
					status = "-X--"
					currentP = -radius + distanceFromSensor(l)
					hintCenter()
				} else {
					// Inconclusive (l, rr), keep last value
					status = "-X-X"
					currentP = lastP
					hintRight()
				}
			} else {
				if isWhite(rr) {
					// Between l and r
					status = "-XX-"
					currentP = lastP
					hintBetween(r, l)
				} else {
					// Angle (l, r, rr)
					status = "-XXX"
					currentP = lastP
					hintRight()
				}
			}
		}
	} else {
		if isWhite(l) {
			// left ll -
			if isWhite(r) {
				if isWhite(rr) {
					// Far left
					status = "X---"
					currentP = -(radius * 3) - distanceFromSensor(ll)
					hintLeft()
				} else {
					// Inconclusive (ll, rr)
					status = "X--X"
					currentP = lastP
					hintBetween(ll, rr)
				}
			} else {
				if isWhite(rr) {
					// Inconclusive (ll, r)
					status = "X-X-"
					currentP = lastP
					hintLeft()
				} else {
					// Inconclusive (ll, r, rr)
					status = "X-XX"
					currentP = lastP
					hintLeft()
				}
			}
		} else {
			// left ll l
			if isWhite(r) {
				if isWhite(rr) {
					// Angle (ll, l)
					status = "XX--"
					currentP = -(radius * 2) + positionBetweenSensors(ll, l)
					hintLeft()
				} else {
					// Inconclusive (ll, l, rr)
					status = "XX-X"
					currentP = lastP
					hintLeft()
				}
			} else {
				if isWhite(rr) {
					// Angle (ll, l, r)
					status = "XXX-"
					currentP = lastP
					hintLeft()
				} else {
					// Inconclusive (ll, l, r, rr)
					status = "XXXX"
					currentP = lastP
					hintBetween(ll, rr)
				}
			}
		}
	}

	deltaTicks := currentPTicks - lastPTicks
	if deltaTicks < conf.MinDTicks {
		deltaTicks = conf.MinDTicks
	} else if deltaTicks > conf.MaxDTicks {
		deltaTicks = conf.MaxDTicks
	}
	deltaP := currentP - lastP
	currentD := (conf.DTicksBoost * deltaP) / deltaTicks

	factorP := ((currentP * conf.ParamP1) + (sign(currentP) * currentP * currentP * conf.ParamP2)) / conf.ParamPR
	factorD := ((currentD * conf.ParamD1) + (sign(currentD) * currentD * currentD * conf.ParamD2)) / conf.ParamDR
	steering := (factorP + factorD) / conf.InnerBrakeFactor

	// print(status, ll, "[", sumSensor(cLL), "]", l, "[", sumSensor(cL), "]", r, "[", sumSensor(cR), "]", rr, "[", sumSensor(cRR), "]", "p", currentP, "dP", deltaP, "dP/dT", currentD, "dT", deltaTicks/1000, (currentPTicks-lastPTicks)/1000)
	// print(status, "P", currentP, "dP", deltaP, "dP/dT", currentD, "dT", deltaTicks/1000, "T", (currentPTicks-lastPTicks)/1000)
	print(status, "P", currentP, "f", factorP, factorD, "s", steering, "T", (currentPTicks-lastPTicks)/1000)

	if steering > 0 {
		if steering > conf.MaxSteering {
			steering = conf.MaxSteering
		}
		move(conf.MaxSpeed, conf.MaxSpeed-steering, currentPTicks)
	} else if steering < 0 {
		steering = -steering
		if steering > conf.MaxSteering {
			steering = conf.MaxSteering
		}
		move(conf.MaxSpeed-steering, conf.MaxSpeed, currentPTicks)
	} else {
		move(conf.MaxSpeed, conf.MaxSpeed, currentPTicks)
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

func followLine() {
	print("following line")
	for {
		processSensorData()
		if buttons.Enter {
			print("stopping")
			break
		}
		if buttons.Back {
			print("reloading config")
			newConf, err := config.FromFile("greyhound.toml")
			if err != nil {
				print("Error reading conf:", err)
			} else {
				conf = newConf
				print("Configuration reloaded:", conf)
			}
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
	waitOneSecond()
	// moveOneSecond()
	followLine()
}
