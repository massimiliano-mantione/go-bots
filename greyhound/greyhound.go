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

func trimSensor(value int) int {
	if value > conf.SensorWhite {
		return conf.SensorWhite
	}
	return value
}
func isWhite(value int) bool {
	return value >= conf.SensorWhite
}
func distanceFromSensor(value int) int {
	return value * conf.SensorRadius / conf.SensorWhite
}
func positionBetweenSensors(value1 int, value2 int) int {
	return (value1 - value2) * conf.SensorRadius / conf.SensorWhite
}

var currentP int
var currentPTicks int
var lastP int
var lastPTicks int

func processSensorData() {
	read()
	ll, l, r, rr := trimSensor(cLL.Value), trimSensor(cL.Value), trimSensor(cR.Value), trimSensor(cRR.Value)
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
					// We lost the line, keep last direction set to max
					if lastP > 0 {
						status = "---+"
						currentP = radius * 4
					} else if lastP < 0 {
						status = "+---"
						currentP = -radius * 4
					} else {
						status = "----"
						currentP = 0
					}
				} else {
					// Far right
					status = "---X"
					currentP = (radius * 3) + distanceFromSensor(rr)
				}
			} else {
				if isWhite(rr) {
					// Centered on r
					status = "--X-"
					currentP = radius
				} else {
					// Between r and rr
					status = "--XX"
					currentP = (radius * 2) + positionBetweenSensors(r, rr)
				}
			}
		} else {
			// left - l
			if isWhite(r) {
				if isWhite(rr) {
					// Centered on l
					status = "-X--"
					currentP = -radius
				} else {
					// Inconclusive (l, rr), keep last value
					status = "-X-X"
					currentP = lastP
				}
			} else {
				if isWhite(rr) {
					// Between l and r
					status = "-XX-"
					currentP = positionBetweenSensors(l, r)
				} else {
					// Inconclusive (l, r, rr), keep r
					status = "-XXX"
					currentP = radius
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
					currentP = -(radius * 3) - distanceFromSensor(rr)
				} else {
					// Inconclusive (ll, rr), keep last value
					status = "X--X"
					currentP = lastP
				}
			} else {
				if isWhite(rr) {
					// Inconclusive (ll, r), keep last value
					status = "X-X-"
					currentP = lastP
				} else {
					// Inconclusive (ll, r, rr), keep last value
					status = "X-XX"
					currentP = lastP
				}
			}
		} else {
			// left ll l
			if isWhite(r) {
				if isWhite(rr) {
					// Between ll and l
					status = "XX--"
					currentP = -(radius * 2) - positionBetweenSensors(ll, l)
				} else {
					// Inconclusive (ll, l, rr), keep last value
					status = "XX-X"
					currentP = lastP
				}
			} else {
				if isWhite(rr) {
					// Inconclusive (ll, l, r), keep l
					status = "XXX-"
					currentP = -radius
				} else {
					// Inconclusive (ll, l, r, rr), keep last value
					status = "XXXX"
					currentP = lastP
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
	currentD := (currentP - lastP) / deltaTicks

	print(status, ll, "[", cLL.Value, "]", l, "[", cL.Value, "]", r, "[", cR.Value, "]", rr, "[", cRR.Value, "]", "p", currentP, "d", currentD, "t", deltaTicks/1000)

	steering := ((currentP * conf.ParamP1) + (currentP * currentP * conf.ParamP1)) / conf.ParamPR
	steering += ((currentD * conf.ParamD1) + (currentD * currentD * conf.ParamD1)) / conf.ParamDR
	steering /= conf.InnerBrakeFactor

	if steering > 0 {
		if steering > conf.MaxSpeed {
			steering = conf.MaxSpeed
		}
		move(conf.MaxSpeed, conf.MaxSpeed-steering, currentPTicks)
	} else if steering < 0 {
		steering = -steering
		if steering > conf.MaxSpeed {
			steering = conf.MaxSpeed
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
		if buttons.Back {
			print("stopping")
			break
		}
	}
}

func main() {
	handleSignals()
	initialize()
	defer close()

	conf = config.Default()

	// waitEnter()
	waitOneSecond()
	// moveOneSecond()
	followLine()
}
