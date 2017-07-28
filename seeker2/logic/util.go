package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"go-bots/ui"
	"os"
)

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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func cmd(eyesActive bool, frontActive bool) {
	c.EyesActive = eyesActive
	c.FrontActive = frontActive
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
	green := 255 * d.VisionIntensity / config.VisionMaxIntensity
	if d.VisionAngle > 0 {
		c.LedLeftGreen = normalizeLedValue(green - (green * d.VisionAngle / config.VisionMaxAngle))
		c.LedRightGreen = normalizeLedValue(green)
	} else if d.VisionAngle < 0 {
		c.LedLeftGreen = normalizeLedValue(green)
		c.LedRightGreen = normalizeLedValue(green + (green * d.VisionAngle / config.VisionMaxAngle))
	} else {
		c.LedLeftGreen = normalizeLedValue(green)
		c.LedRightGreen = normalizeLedValue(green)
	}
	if d.CornerLeftIsOut {
		c.LedLeftRed = 0
	} else {
		c.LedLeftRed = 0
	}
	if d.CornerRightIsOut {
		c.LedRightRed = 0
	} else {
		c.LedRightRed = 0
	}
}

func checkDone(k ui.KeyEvent) bool {
	if k.Key == ui.Quit || k.Key == ui.Back {

		log(k.Millis, ev3.NoDirection, " *** DONE ***")

		go chooseStrategy(k.Millis)
		return true
	}
	return false
}
func checkQuit(k ui.KeyEvent) {
	if k.Key == ui.Quit || k.Key == ui.Back {
		quit <- true
		return
	}
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
