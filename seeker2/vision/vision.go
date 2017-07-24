package vision

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"os"
)

var currentIntensityLeft int
var currentPositionLeft int
var currentIntensityRight int
var currentPositionRight int

var estimatedIntensityLeft int
var estimatedPositionLeft int
var estimatedIntensityRight int
var estimatedPositionRight int

func irValueToIntensity(value int) int {
	if value > config.VisionFarValue {
		return 0
	}
	return config.VisionMaxValue - value
}

func estimate(d ev3.Direction) (intensity int, angle int, dir ev3.Direction) {
	intensity = (estimatedIntensityLeft + estimatedIntensityRight) / 2
	if intensity > 0 {
		leftAngle := estimatedPositionLeft - config.VisionMaxAngle
		rightAngle := estimatedPositionRight + config.VisionMaxAngle
		leftRatio := intensity - (estimatedIntensityRight / 2)
		rightRatio := intensity - (estimatedIntensityLeft / 2)
		angle = ((leftAngle * leftRatio) + (rightAngle * rightRatio)) / intensity
	} else {
		angle = 0
	}
	return intensity, angle, d
}

// Reset resets the vision state
func Reset() {
	currentIntensityLeft = 0
	currentPositionLeft = 0
	currentIntensityRight = 0
	currentPositionRight = 0
	estimatedIntensityLeft = 0
	estimatedPositionLeft = 0
	estimatedIntensityRight = 0
	estimatedPositionRight = 0
}

func switchDirection(pos int) {
	if currentIntensityLeft == 0 {
		estimatedIntensityLeft = 0
	}
	currentIntensityLeft = 0
	currentPositionLeft = pos
	if currentIntensityRight == 0 {
		estimatedIntensityRight = 0
	}
	currentIntensityRight = 0
	currentPositionRight = pos
}

// Process processes IR sensor data
func Process(millis int, d ev3.Direction, pos int, rightValue int, leftValue int) (intensity int, angle int, dir ev3.Direction) {
	if d == ev3.Right && pos >= config.VisionThresholdPosition {
		dir = ev3.Left
		switchDirection(pos)
		return estimate(dir)
	} else if d == ev3.Left && pos <= -config.VisionThresholdPosition {
		dir = ev3.Right
		switchDirection(pos)
		return estimate(dir)
	} else {
		dir = d
	}

	leftIntensity := irValueToIntensity(leftValue)
	if leftIntensity > currentIntensityLeft {
		currentIntensityLeft = leftIntensity
		currentPositionLeft = pos
	} else if leftIntensity <= currentIntensityLeft {
		estimatedIntensityLeft = currentIntensityLeft
		estimatedPositionLeft = currentPositionLeft
	}

	rightIntensity := irValueToIntensity(rightValue)
	if rightIntensity > currentIntensityRight {
		currentIntensityRight = rightIntensity
		currentPositionRight = pos
	} else if rightIntensity <= currentIntensityRight {
		estimatedIntensityRight = currentIntensityRight
		estimatedPositionRight = currentPositionRight
	}

	intens, ang, _ := estimate(dir)
	fmt.Fprintln(os.Stderr, " - VISION PROCESS", d, pos, rightValue, leftValue)
	fmt.Fprintln(os.Stderr, " - VISION   STATE", dir, estimatedIntensityLeft, estimatedPositionLeft, estimatedIntensityRight, estimatedPositionRight)
	fmt.Fprintln(os.Stderr, " - VISION     RES", intens, ang)

	return estimate(dir)
}
