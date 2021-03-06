package main

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/greyhound/config"
	"go-bots/greyhound/data"
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

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

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
	data.Init(sensorReadNames)

	buttons = ev3.OpenButtons(true)

	devs = ev3.Scan(&ev3.OutPortModes{
		OutA: ev3.OutPortModeDcMotor,
		OutB: ev3.OutPortModeDcMotor,
		OutC: ev3.OutPortModeDcMotor,
		OutD: ev3.OutPortModeDcMotor,
	})

	ledLG = ev3.OpenTextW(devs.LedLeftGreen, ev3.Brightness)
	ledLR = ev3.OpenTextW(devs.LedLeftRed, ev3.Brightness)
	ledRG = ev3.OpenTextW(devs.LedRightGreen, ev3.Brightness)
	ledRR = ev3.OpenTextW(devs.LedRightRed, ev3.Brightness)
	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

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

	// Set leds off
	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()

	// Close sensor values
	closeSensors()
}

func leds(rl int, rr int, gr int, gl int) {
	ledLG.Value = gl
	ledLR.Value = rl
	ledRG.Value = gr
	ledRR.Value = rr
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()
}

var lastMoveTicks int
var lastSpeedLeft int
var lastSpeedRight int

const accelSpeedFactor int = 10000

func stop() {
	motorL1.Value = 0
	motorL2.Value = 0
	motorR1.Value = 0
	motorR2.Value = 0

	motorL1.Sync()
	motorL2.Sync()
	motorR1.Sync()
	motorR2.Sync()
}

func move(left int, right int, now int) {
	ticks := now - lastMoveTicks
	lastMoveTicks = now

	if left > 100 {
		left = 100
	} else if left < -100 {
		left = -100
	}
	if right > 100 {
		right = 100
	} else if right < -100 {
		right = -100
	}

	right *= accelSpeedFactor
	left *= accelSpeedFactor

	nextSpeedLeft := lastSpeedLeft
	nextSpeedRight := lastSpeedRight
	delta := ticks * conf.AccelPerTicksN
	delta /= conf.AccelPerTicksD

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
	data.Reset()
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
func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
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
		pos = conf.SensorHole + conf.SensorRadius*2 + distanceFromSensor(r)
		hint = 0
		cross, out = false, false
	case sensorReadRB:
		pos = conf.SensorHole + conf.SensorRadius + positionBetweenSensors(b, r)
		hint = 1
		cross, out = false, false
	case sensorReadL:
		pos = -conf.SensorHole - conf.SensorRadius*2 - distanceFromSensor(l)
		hint = 0
		cross, out = false, false
	case sensorReadLB:
		pos = -conf.SensorHole - conf.SensorRadius + positionBetweenSensors(l, b)
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
		pos = conf.SensorHole + conf.SensorRadius + positionBetweenSensors(f, r)
		hint = -1
		cross, out = false, false
	case sensorReadFRB:
		pos = conf.SensorHole + conf.SensorRadius + positionBetweenSensors((f+b)/2, r)
		hint = 0
		cross, out = false, false
	case sensorReadFL:
		pos = -conf.SensorHole - conf.SensorRadius + positionBetweenSensors(l, f)
		hint = 1
		cross, out = false, false
	case sensorReadFLB:
		pos = -conf.SensorHole - conf.SensorRadius + positionBetweenSensors(l, (f+b)/2)
		hint = 0
		cross, out = false, false
	default:
		print("Error: reading", sensorRead)
	}

	if out {
		if pos < 0 {
			leds(255, 0, 0, 0)
		} else if pos > 0 {
			leds(0, 255, 0, 0)
		} else {
			leds(255, 255, 0, 0)
		}
	} else {
		if pos < 0 {
			leds(0, 0, 255, (-pos*254)/conf.MaxPos)
		} else if pos > 0 {
			leds(0, 0, (pos*254)/conf.MaxPos, 255)
		} else {
			leds(0, 0, 255, 255)
		}
	}

	return
}

func moveUntilOut() {
	print("move until out")
	leds(0, 0, 255, 255)
	for {
		now := currentTicks()
		cF.Sync()
		v := trimSensor(cF)
		if isOnTrack(v) {
			move(30, 30, now)
		} else {
			leds(0, 0, 255, 255)
			return
		}
	}
}
func turn50ms() {
	print("turn 50ms")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(-30, 30, now)
		if elapsed >= 50000 {
			stop()
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

func turnOneSecond() {
	print("turn one second")
	start := currentTicks()
	for {
		now := currentTicks()
		elapsed := now - start
		move(conf.MaxSpeed, -conf.MaxSpeed, now)
		if elapsed >= 50000 {
			break
		}
	}
}

func followLine(lastGivenTicks int) {
	print("following line")

	lastTicks := lastGivenTicks
	lastPos := 0
	lastNonZeroDirection := 0
	posI := 0
	posE := 0
	dt := 0
	dMillis := 0
	powerLeft := 0
	powerRight := 0
	estimatedSpeedLeft := 0
	estimatedSpeedRight := 0

	onTrackTicks := 0
	outLeftPowerTarget := 0
	outLeftPowerDelta := 0
	outRightPowerTarget := 0
	outRightPowerDelta := 0

	attenuationTicks := currentTicks()

	for {
		now := currentTicks()
		sr, pos, hint, cross, out := processSensorData()
		posD := 0
		posDAverage := 0

		var maxSpeed int
		if conf.SlowSpeed1 > 0 && now < conf.SlowStart1 && now > conf.SlowEnd1 {
			maxSpeed = conf.SlowSpeed1
		} else if conf.SlowSpeed2 > 0 && now < conf.SlowStart2 && now > conf.SlowEnd2 {
			maxSpeed = conf.SlowSpeed2
		} else if conf.SlowSpeed3 > 0 && now < conf.SlowStart3 && now > conf.SlowEnd3 {
			maxSpeed = conf.SlowSpeed3
		} else if conf.SlowSpeed4 > 0 && now < conf.SlowStart4 && now > conf.SlowEnd4 {
			maxSpeed = conf.SlowSpeed4
		} else if conf.Timeout > 0 && now > conf.Timeout {
			maxSpeed = 0
		} else {
			maxSpeed = conf.MaxSpeed
		}

		if out {
			pos = conf.SensorRadius * 6 * lastNonZeroDirection
			hint = sign(pos)
			posD = 0
		} else if cross {
			pos = lastPos
			posD = 0
		} else {
			dt = now - lastTicks
			// cdt is a "trimmed" dt
			cdt := dt
			if cdt < 2000 {
				cdt = 2000
				dMillis = 1
			} else if cdt > 100000 {
				cdt = 100000
				dMillis = cdt / 1000
			} else {
				dMillis = cdt / 1000
			}

			if onTrackTicks < 100 {
				posD = 0
			} else {
				posD = ((pos - lastPos) * 100000) / cdt
				posD += hint
				if posD > 2000 {
					posD = 2000
				} else if posD < -2000 {
					posD = -2000
				}
			}
		}

		// Update posD average
		posDAverageLost := (posDAverage * (dMillis)) / conf.DAvgMillis
		posDAverageGained := (posD * (dMillis)) / conf.DAvgMillis
		posDAverage = posDAverage + posDAverageGained - posDAverageLost

		// Apply attenuations and updates
		pos2e := pos * pos / conf.KEReduction
		for attenuationTicks < now {
			attenuationTicks += 1000

			// I
			posI *= conf.KIrn
			posI /= conf.KIrd
			posI += pos

			// E
			posE *= conf.KErn
			posE /= conf.KErd
			posE += pos2e
			if posE > 30000 {
				posE = 30000
			}

			// Estimate speed
			estimatedSpeedLeft *= conf.SpeedEstRn
			estimatedSpeedLeft /= conf.SpeedEstRd
			estimatedSpeedLeft += powerLeft
			estimatedSpeedRight *= conf.SpeedEstRn
			estimatedSpeedRight /= conf.SpeedEstRd
			estimatedSpeedRight += powerRight

			// Reduce out power deltas
			outLeftPowerDelta *= conf.OutPowerRn
			outLeftPowerDelta /= conf.OutPowerRd
			outRightPowerDelta *= conf.OutPowerRn
			outRightPowerDelta /= conf.OutPowerRd
		}

		// print(dMillis, "power", powerLeft, powerRight, "speed", estimatedSpeedLeft, estimatedSpeedRight)

		// Compute factors
		factorP := (pos * conf.KPn * maxSpeed) / (conf.MaxPos * conf.KPd)
		factorD := (posDAverage * conf.KDn * maxSpeed) / (conf.MaxPosD * conf.KDd)
		factorI := (posI * conf.KIn * maxSpeed) / conf.KId

		factorE := ((posE / conf.KELimit) * conf.KEn) / conf.KEd

		// Limit slowness factor
		if factorE > conf.MaxSlowPC {
			factorE = conf.MaxSlowPC
		}

		// Compute steering
		steering := factorP + factorD + factorI
		if out {
			if onTrackTicks > 0 {
				// estimatedDirection := sign(estimatedSpeedRight - estimatedSpeedLeft)

				initialOuterPowerDelta := 0
				initialInnerPowerDelta := 0
				if ticksToMillis(onTrackTicks) > conf.OutTimeMs {
					initialOuterPowerDelta = conf.OutPowerMax - conf.OutPowerMin
					initialInnerPowerDelta = -conf.OutPowerMax
				}

				if lastNonZeroDirection > 0 {
					outLeftPowerTarget = conf.OutPowerMin
					outRightPowerTarget = 0
					outLeftPowerDelta = initialOuterPowerDelta
					outRightPowerDelta = initialInnerPowerDelta
				} else {
					outLeftPowerTarget = 0
					outRightPowerTarget = conf.OutPowerMin
					outLeftPowerDelta = initialInnerPowerDelta
					outRightPowerDelta = initialOuterPowerDelta
				}

				print("OUT INIT", ticksToMillis(onTrackTicks), lastNonZeroDirection, outLeftPowerDelta, outRightPowerDelta)
			}

			onTrackTicks = 0

			powerLeft = outLeftPowerTarget + outLeftPowerDelta
			powerRight = outRightPowerTarget + outRightPowerDelta
		} else {
			onTrackTicks += dt

			// Apply slowdown
			actualMaxSpeed := (maxSpeed * (100 - factorE)) / 100
			actualSteering := (steering * (100 - factorE)) / 100
			maxSteering := (actualMaxSpeed * conf.MaxSteeringPC) / 100

			// Compute motor powers
			if actualSteering > 0 {
				if actualSteering > maxSteering {
					actualSteering = maxSteering
				}
				powerLeft = actualMaxSpeed
				powerRight = actualMaxSpeed - actualSteering
			} else if actualSteering < 0 {
				actualSteering = -actualSteering
				if actualSteering > maxSteering {
					actualSteering = maxSteering
				}
				powerLeft = actualMaxSpeed - actualSteering
				powerRight = actualMaxSpeed
			} else {
				powerLeft = actualMaxSpeed
				powerRight = actualMaxSpeed
			}
		}

		// Apply power
		move(powerLeft, powerRight, now)

		// Compute last values for next round
		if pos > 0 {
			lastNonZeroDirection = 1
		} else if pos < 0 {
			lastNonZeroDirection = -1
		}
		lastTicks, lastPos = now, pos

		// Store data
		data.Store(uint32(now),
			uint16(dt),
			int16(pos),
			int16(posD),
			int16(posI),
			int16(posE),
			int16(factorP),
			int16(factorD),
			int16(factorI),
			int16(factorE),
			uint8(sr),
			int16(estimatedSpeedLeft),
			int16(estimatedSpeedRight),
			int8(powerLeft),
			int8(powerRight))

		// Check stop command
		if buttons.Enter {
			print("stopping")
			stop()
			data.Print()
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

	// waitOneSecond()
	// moveUntilOut()
	// turn50ms()
	// waitOneSecond()
	// waitOneSecond()
	// waitOneSecond()

	// waitOneSecond()
	// moveOneSecond()
	// turnOneSecond()
}
