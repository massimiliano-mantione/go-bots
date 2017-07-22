package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/ui"
	"go-bots/xl4/config"
	"os"
	"time"
)

// Data contains readings from sensors (adjusted for logic)
type Data struct {
	Start            time.Time
	Millis           int
	CornerRightIsOut bool
	CornerLeftIsOut  bool
	CornerRight      int
	CornerLeft       int
	IrLeftValue      int
	IrRightValue     int
}

// Commands contains commands for motors and leds
type Commands struct {
	Millis        int
	SpeedRight    int
	SpeedLeft     int
	LedRightRed   int
	LedRightGreen int
	LedLeftRed    int
	LedLeftGreen  int
}

var data <-chan Data
var commandProcessor func(*Commands)
var keys <-chan ui.KeyEvent
var quit chan<- bool

// Init initializes the logic module
func Init(d <-chan Data, c func(*Commands), k <-chan ui.KeyEvent, q chan<- bool) {
	data = d
	commandProcessor = c
	keys = k
	quit = q
}

var c = Commands{}

func log(now int, dir ev3.Direction, msg string) {
	dirString := ""
	if dir == ev3.Left {
		dirString = "LEFT"
	} else if dir == ev3.Right {
		dirString = "RIGHT"
	} else {
		dirString = "NONE"
	}

	fmt.Fprintln(os.Stderr, now, dirString, msg)
}

func cmd() {
	commandProcessor(&c)
}

func handleTime(d Data, start int) (now int, elapsed int) {
	now = d.Millis
	c.Millis = now
	elapsed = now - start
	return
}

func speed(left int, right int) {
	c.SpeedLeft = left
	c.SpeedRight = right
}

func normalizeLedValue(v int) int {
	if v > 255 {
		v = 255
	}
	if v < 0 {
		v = 0
	}
	return v
}

func leds(leftGreen int, rightGreen int, leftRed int, rightRed int) {
	leftGreen = normalizeLedValue(leftGreen)
	rightGreen = normalizeLedValue(rightGreen)
	leftRed = normalizeLedValue(leftRed)
	rightRed = normalizeLedValue(rightRed)
	c.LedLeftGreen = leftGreen
	c.LedRightGreen = rightGreen
	c.LedLeftRed = leftRed
	c.LedRightRed = rightRed
}

func ledsFromData(d Data) {
	c.LedLeftGreen = 255 * d.IrLeftValue / 100
	c.LedRightGreen = 255 * d.IrRightValue / 100
	if d.CornerLeftIsOut {
		c.LedLeftRed = 255
	} else {
		c.LedLeftRed = 0
	}
	if d.CornerRightIsOut {
		c.LedRightRed = 255
	} else {
		c.LedRightRed = 0
	}
}

func checkEnd(k ui.KeyEvent) {
	if k.Key == ui.Quit || k.Key == ui.Back {
		quit <- true
		return
	}
}

func chooseStrategy(start int) {
	strategy := seek
	strategyIsGoForward := false
	var dir ev3.Direction = ev3.Left
	leds(0, 0, 0, 0)
	speed(0, 0)
	cmd()

	for {
		select {
		case d := <-data:
			now := d.Millis
			c.Millis = now
		case k := <-keys:
			checkEnd(k)
			if k.Key == ui.Enter {
				go pauseBeforeBegin(k.Millis, strategy, dir)
				return
			} else if k.Key == ui.Left {
				dir = ev3.Left
				strategy = circle
				strategyIsGoForward = false
				leds(255, 0, 255, 0)
				fmt.Fprintln(os.Stderr, "chooseStrategy circle left")
			} else if k.Key == ui.Right {
				dir = ev3.Right
				strategy = circle
				strategyIsGoForward = false
				leds(0, 255, 0, 255)
				fmt.Fprintln(os.Stderr, "chooseStrategy circle right")
			} else if k.Key == ui.Up {
				if strategyIsGoForward {
					strategy = seek
					strategyIsGoForward = false
					leds(0, 0, 0, 0)
					fmt.Fprintln(os.Stderr, "chooseStrategy seek")
				} else {
					strategy = goForward
					strategyIsGoForward = true
					if dir == ev3.Left {
						leds(255, 0, 0, 0)
						fmt.Fprintln(os.Stderr, "chooseStrategy forward left")
					} else {
						leds(0, 255, 0, 0)
						fmt.Fprintln(os.Stderr, "chooseStrategy forward right")
					}
				}
			} else if k.Key == ui.Down {
				strategy = turnBack
				strategyIsGoForward = false
				if dir == ev3.Left {
					leds(0, 0, 255, 0)
					fmt.Fprintln(os.Stderr, "chooseStrategy back left")
				} else {
					leds(0, 0, 0, 255)
					fmt.Fprintln(os.Stderr, "chooseStrategy back right")
				}
			}
			speed(0, 0)
			cmd()
		}
	}
}

func pauseBeforeBegin(start int, strategy func(int, ev3.Direction), dir ev3.Direction) {
	for {
		select {
		case d := <-data:
			now := d.Millis
			c.Millis = now
			elapsed := now - start
			if elapsed >= config.StartTime {
				go strategy(now, dir)
				return
			}
			speed(0, 0)
			intensity := ((elapsed % 1000) * 255) / (config.StartTime / 5)
			if elapsed > (config.StartTime * 4 / 5) {
				leds(intensity, intensity, intensity, intensity)
			} else {
				leds(0, 0, intensity, intensity)
			}
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}
}

func checkVision(d Data, now int) bool {
	//	return d.IrLeftValue < config.MaxIrDistance || d.IrRightValue < config.MaxIrDistance
	result := d.IrLeftValue < config.MaxIrDistance || d.IrRightValue < config.MaxIrDistance
	if result {
		fmt.Fprintln(os.Stderr, "CheckVision")
		go track(now)
	}
	return result
}

func checkBorder(d Data, now int) bool {
	if d.CornerLeftIsOut {
		if d.CornerRightIsOut {
			go back(now, ev3.NoDirection)
			return true
		}
		go back(now, ev3.Left)
		return true
	}
	if d.CornerRightIsOut {
		if d.CornerLeftIsOut {
			go back(now, ev3.NoDirection)
			return true
		}
		go back(now, ev3.Right)
		return true
	}
	return false
}

func seekMove(start int, dir ev3.Direction, leftSpeed int, rightSpeed int, duration int, ignoreBorder bool) (done bool, now int) {

	fmt.Fprintln(os.Stderr, "seekMove", dir, leftSpeed, rightSpeed, duration)

	for {
		select {
		case d := <-data:
			now, elapsed := handleTime(d, start)
			if elapsed >= duration {

				fmt.Fprintln(os.Stderr, "seekMove elapsed", now, duration)

				return false, now
			}

			if checkVision(d, now) {
				return true, now
			}
			if (!ignoreBorder) && checkBorder(d, now) {
				return true, now
			}

			speed(leftSpeed, rightSpeed)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}
}

func back(start int, dir ev3.Direction) {

	done, now := false, start
	fmt.Fprintln(os.Stderr, "BACK", dir)

	if dir == ev3.Right {
		done, now = seekMove(now, dir, -config.BackTurn1SpeedInner, -config.BackTurn1SpeedOuter, config.BackTurn1Millis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, config.BackTurn2Speed, -config.BackTurn2Speed, config.BackTurn2Millis, true)
		if done {
			return
		}
	} else if dir == ev3.Left {
		done, now = seekMove(now, dir, -config.BackTurn1SpeedOuter, -config.BackTurn1SpeedInner, config.BackTurn1Millis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, -config.BackTurn2Speed, config.BackTurn2Speed, config.BackTurn2Millis, true)
		if done {
			return
		}
	} else {
		dir = ev3.Right
		done, now = seekMove(now, dir, -config.BackMoveSpeed, -config.BackMoveSpeed, config.BackMoveMillis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, config.BackTurn3Speed, -config.BackTurn3Speed, config.BackTurn3Millis, true)
		if done {
			return
		}
	}

	go seek(now, dir)
}

func seek(start int, dir ev3.Direction) {

	fmt.Fprintln(os.Stderr, "SEEK", dir)

	done, now := false, start
	for {

		fmt.Fprintln(os.Stderr, "SEEK MOVE", dir, now)

		done, now = seekMove(now, dir, config.SeekMoveSpeed, config.SeekMoveSpeed, config.SeekMoveMillis, false)
		if done {
			return
		}

		fmt.Fprintln(os.Stderr, "SEEK TURN", dir, now)

		done, now = seekMove(now, dir, config.SeekTurnSpeed*ev3.LeftTurnVersor(dir), config.SeekTurnSpeed*ev3.RightTurnVersor(dir), config.SeekTurnMillis, false)
		if done {
			return
		}
		dir = ev3.ChangeDirection(dir)
	}
}

func circle(start int, dir ev3.Direction) {
	now, elapsed := start, 0

	log(now, dir, "CIRCLE find border")
findBorder:
	for {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)
			if dir == ev3.Right {
				if d.CornerRightIsOut {
					break findBorder
				}
				if elapsed < config.CircleFindBorderMillis {
					speed(config.CircleFindBorderOuterSpeed, -config.CircleFindBorderInnerSpeed)
				} else {
					speed(config.CircleFindBorderOuterSpeedSlow, -config.CircleFindBorderInnerSpeedSlow)
				}
			} else {
				if d.CornerLeftIsOut {
					break findBorder
				}
				if elapsed < config.CircleFindBorderMillis {
					speed(-config.CircleFindBorderInnerSpeed, config.CircleFindBorderOuterSpeed)
				} else {
					speed(-config.CircleFindBorderInnerSpeedSlow, config.CircleFindBorderOuterSpeedSlow)
				}
			}
			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	log(now, dir, "CIRCLE start")
	dir = ev3.ChangeDirection(dir)
	borderFoundTime := elapsed
	for elapsed-borderFoundTime < config.CircleMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			if dir == ev3.Right {
				cornerValue := d.CornerLeft
				adjustInner := cornerValue * config.CircleAdjustInnerMax / 100
				inner := config.CircleInnerSpeedRight - adjustInner
				speed(config.CircleOuterSpeed, inner)
			} else {
				cornerValue := d.CornerRight
				adjustInner := cornerValue * config.CircleAdjustInnerMax / 100
				inner := config.CircleInnerSpeedLeft - adjustInner
				speed(inner, config.CircleOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	log(now, dir, "CIRCLE spiral")
	circleDoneTime := elapsed
	for elapsed-circleDoneTime < config.CircleSpiralMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			if dir == ev3.Right {
				speed(config.CircleSpiralOuterSpeed, config.CircleSpiralInnerSpeed)
			} else {
				speed(config.CircleSpiralInnerSpeed, config.CircleSpiralOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	log(now, dir, "CIRCLE done")
	go seek(now, ev3.ChangeDirection(dir))
}

func goForward(start int, dir ev3.Direction) {
	now, elapsed := start, 0

	fmt.Fprintln(os.Stderr, "goForward", now, dir)

	for elapsed < config.GoForwardMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}
			if checkBorder(d, now) {
				return
			}

			speed(config.GoForwardSpeed, config.GoForwardSpeed)

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	fmt.Fprintln(os.Stderr, "goForward turn", now, dir)

	forwardDone := elapsed
	for elapsed-forwardDone < config.GoForwardTurnMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}
			if checkBorder(d, now) {
				return
			}

			if dir == ev3.Right {
				speed(config.GoForwardTurnOuterSpeed, config.GoForwardTurnInnerSpeed)
			} else {
				speed(config.GoForwardTurnInnerSpeed, config.GoForwardTurnOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	fmt.Fprintln(os.Stderr, "goForward done", now, dir)
	go seek(now, ev3.ChangeDirection(dir))
}

func turnBack(start int, dir ev3.Direction) {
	now, elapsed := start, 0

	fmt.Fprintln(os.Stderr, "turnBack", now, dir)

	for elapsed < config.TurnBackMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}
			if checkBorder(d, now) {
				return
			}

			if dir == ev3.Right {
				speed(config.TurnBackOuterSpeed, config.TurnBackInnerSpeed)
			} else {
				speed(config.TurnBackInnerSpeed, config.TurnBackOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	fmt.Fprintln(os.Stderr, "turnBack move", now, dir)

	turnBackDone := elapsed
	for elapsed-turnBackDone < config.TurnBackMoveMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}
			if checkBorder(d, now) {
				return
			}

			speed(config.TurnBackMoveSpeed, config.TurnBackMoveSpeed)

			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}

	fmt.Fprintln(os.Stderr, "turnBack done", now, dir)
	go seek(now, ev3.ChangeDirection(dir))
}

func track(start int) {
	now, _ := start, 0
	var dir ev3.Direction = ev3.Right

	for {
		select {
		case d := <-data:
			now, _ = handleTime(d, start)

			if d.IrLeftValue >= config.MaxIrDistance && d.IrRightValue >= config.MaxIrDistance {
				go seek(now, dir)
				return
			}
			if checkBorder(d, now) {
				return
			}

			if d.IrLeftValue >= config.MaxIrDistance {
				speed(config.TrackOnly1SensorOuterSpeed, config.TrackOnly1SensorInnerSpeed)
				dir = ev3.Right
			} else if d.IrRightValue >= config.MaxIrDistance {
				speed(config.TrackOnly1SensorInnerSpeed, config.TrackOnly1SensorOuterSpeed)
				dir = ev3.Left
			} else {
				difference := d.IrLeftValue - d.IrRightValue

				if difference > config.TrackCenterZone {
					speed(config.TrackSpeed, config.TrackSpeed-(difference*config.TrackDifferenceCoefficent))
					dir = ev3.Right
				} else if difference < -config.TrackCenterZone {
					speed(config.TrackSpeed+(difference*config.TrackDifferenceCoefficent), config.TrackSpeed)
					dir = ev3.Left
				} else {
					speed(config.TrackSpeed, config.TrackSpeed)
				}

			}
			// fmt.Fprintln(os.Stderr, "TRACK time ", now, ", speed", c.SpeedLeft, c.SpeedRight, ", IRsensors", d.IrLeftValue, d.IrRightValue)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}
}

// Run starts the logic module
func Run() {
	go chooseStrategy(0)
}
