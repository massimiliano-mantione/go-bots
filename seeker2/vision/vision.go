package vision

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"os"
)

var currentIntensityLeft int
var currentPositionLeft int
var hasLeftEstimation bool
var currentIntensityRight int
var currentPositionRight int
var hasRightEstimation bool

var estimatedIntensityLeft int
var estimatedPositionLeft int
var estimatedIntensityRight int
var estimatedPositionRight int

func farValueAtPosition(pos int) (farValueLeft int, farValueRight int) {
	if pos > 0 {
		farValueLeft = config.VisionFarValueSide + (config.VisionFarValueDelta * pos / config.VisionMaxPosition)
		farValueRight = config.VisionFarValueFront - (config.VisionFarValueDelta * pos / config.VisionMaxPosition)
	} else {
		farValueLeft = config.VisionFarValueFront + (config.VisionFarValueDelta * pos / config.VisionMaxPosition)
		farValueRight = config.VisionFarValueSide - (config.VisionFarValueDelta * pos / config.VisionMaxPosition)
	}
	return
}

func irValuesToIntensity(leftValue int, rightValue int, pos int) (leftIntensity int, rightIntensity int) {
	leftLimit, rightLimit := farValueAtPosition(pos)
	if leftValue >= leftLimit {
		leftValue = 100
	}
	if rightValue >= rightLimit {
		rightValue = 100
	}
	return (100 - leftValue), (100 - rightValue)
}

func estimate(d ev3.Direction) (intensity int, angle int, dir ev3.Direction) {
	leftAngle := estimatedPositionLeft - config.VisionMaxAngle
	rightAngle := estimatedPositionRight + config.VisionMaxAngle
	if estimatedIntensityLeft > estimatedIntensityRight {
		return estimatedIntensityLeft, leftAngle, d
	} else if estimatedIntensityRight > estimatedIntensityLeft {
		return estimatedIntensityRight, rightAngle, d
	} else {
		return estimatedIntensityRight, (leftAngle + rightAngle) / 2, d
	}
	/*
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
	*/
}

// Reset resets the vision state
func Reset() {
	currentIntensityLeft = 0
	currentPositionLeft = 0
	hasLeftEstimation = false
	currentIntensityRight = 0
	currentPositionRight = 0
	hasRightEstimation = false
	estimatedIntensityLeft = 0
	estimatedPositionLeft = 0
	estimatedIntensityRight = 0
	estimatedPositionRight = 0
}

func switchDirection(pos int, leftIntensity int, rightIntensity int, dir ev3.Direction) ev3.Direction {
	if leftIntensity > 0 && !hasLeftEstimation {
		estimatedIntensityLeft = leftIntensity
		estimatedPositionLeft = pos
	} else if currentIntensityLeft == 0 {
		estimatedIntensityLeft = 0
	}
	currentIntensityLeft = 0
	currentPositionLeft = pos
	hasLeftEstimation = false

	if rightIntensity > 0 && !hasRightEstimation {
		estimatedIntensityRight = rightIntensity
		estimatedPositionRight = pos
	} else if currentIntensityRight == 0 {
		estimatedIntensityRight = 0
	}
	currentIntensityRight = 0
	currentPositionRight = pos
	hasRightEstimation = false

	return ev3.ChangeDirection(dir)
}

// Process processes IR sensor data
func Process(millis int, d ev3.Direction, pos int, leftValue int, rightValue int) (intensity int, angle int, dir ev3.Direction) {
	leftIntensity, rightIntensity := irValuesToIntensity(leftValue, rightValue, pos)

	if d == ev3.Right && pos >= config.VisionThresholdPosition {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else if d == ev3.Left && pos <= -config.VisionThresholdPosition {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
		// } else if hasLeftEstimation && pos-estimatedPositionLeft > config.VisionSpotWidth && estimatedPositionLeft < config.VisionSpotSearchWidth {
		//	dir = switchDirection(pos, leftIntensity, rightIntensity, d)
		// } else if hasRightEstimation && pos-estimatedPositionRight > config.VisionSpotWidth && estimatedPositionRight > -config.VisionSpotSearchWidth {
		//	dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else {
		dir = d

		if leftIntensity > currentIntensityLeft {
			currentIntensityLeft = leftIntensity
			currentPositionLeft = pos
		} else if leftIntensity <= currentIntensityLeft {
			estimatedIntensityLeft = currentIntensityLeft
			estimatedPositionLeft = currentPositionLeft
			hasLeftEstimation = true
		}

		if rightIntensity > currentIntensityRight {
			currentIntensityRight = rightIntensity
			currentPositionRight = pos
		} else if rightIntensity <= currentIntensityRight {
			estimatedIntensityRight = currentIntensityRight
			estimatedPositionRight = currentPositionRight
			hasRightEstimation = true
		}

		intens, ang, _ := estimate(dir)
		fmt.Fprintln(os.Stderr, " - VISION PROCESS", d, pos, leftIntensity, rightIntensity)
		fmt.Fprintln(os.Stderr, " - VISION   STATE", dir, estimatedIntensityLeft, estimatedPositionLeft, estimatedIntensityRight, estimatedPositionRight)
		fmt.Fprintln(os.Stderr, " - VISION     RES", intens, ang)
	}

	return estimate(dir)
}
