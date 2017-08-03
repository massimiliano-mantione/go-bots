package vision

import (
	"go-bots/ev3"
	"go-bots/scooba/config"
)

var firstIntensityLeft int
var firstPositionLeft int
var currentIntensityLeft int
var currentPositionLeft int
var hasLeftEstimation bool

var firstIntensityRight int
var firstPositionRight int
var currentIntensityRight int
var currentPositionRight int
var hasRightEstimation bool

var estimatedIntensityLeft int
var estimatedPositionLeft int
var estimatedIntensityRight int
var estimatedPositionRight int

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

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

func positionToAngle(pos int) int {
	return pos * 9 / 25
}

func estimate(d ev3.Direction) (intensity int, angle int, dir ev3.Direction) {
	leftAngle := positionToAngle(estimatedPositionLeft) - 45
	rightAngle := positionToAngle(estimatedPositionRight) + 45
	if estimatedIntensityLeft > estimatedIntensityRight {
		return estimatedIntensityLeft, leftAngle, d
	} else if estimatedIntensityRight > estimatedIntensityLeft {
		return estimatedIntensityRight, rightAngle, d
	} else {
		return estimatedIntensityRight, (leftAngle + rightAngle) / 2, d
	}
}

// Reset resets the vision state
func Reset() {
	firstIntensityLeft = 0
	firstPositionLeft = 0
	currentIntensityLeft = 0
	currentPositionLeft = 0
	hasLeftEstimation = false
	firstIntensityRight = 0
	firstPositionRight = 0
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
	firstIntensityLeft = 0
	firstPositionLeft = 0
	currentIntensityLeft = 0
	currentPositionLeft = pos
	hasLeftEstimation = false

	if rightIntensity > 0 && !hasRightEstimation {
		estimatedIntensityRight = rightIntensity
		estimatedPositionRight = pos
	} else if currentIntensityRight == 0 {
		estimatedIntensityRight = 0
	}
	firstIntensityRight = 0
	firstPositionRight = 0
	currentIntensityRight = 0
	currentPositionRight = pos
	hasRightEstimation = false

	return ev3.ChangeDirection(dir)
}

func estimationIsOld(estimationPosition int, pos int) bool {
	return abs(pos-estimationPosition) > config.VisionSpotWidth && abs(estimationPosition) > config.VisionSpotSearchWidth
}

func computeEstimatedPositionCorrection(firstPosition int, firstIntensity int, currentPosition int, currentIntensity int, pos int, intensity int) int {
	risingPositionDelta := currentPosition - firstPosition
	descendingPositionDelta := pos - currentPosition

	if risingPositionDelta == 0 || descendingPositionDelta == 0 {
		return 0
	}

	risingIntensityDelta := currentIntensity - firstIntensity
	descendingIntensityDelta := currentIntensity - intensity

	risingRatio := risingIntensityDelta / risingPositionDelta
	descendingRatio := descendingIntensityDelta / descendingPositionDelta

	if risingRatio >= descendingRatio || descendingRatio == 0 {
		return 0
	}

	return risingPositionDelta - (risingPositionDelta * risingRatio / descendingRatio)
}

// Process processes IR sensor data
func Process(millis int, d ev3.Direction, pos int, leftValue int, rightValue int) (intensity int, angle int, dir ev3.Direction) {
	leftIntensity, rightIntensity := irValuesToIntensity(leftValue, rightValue, pos)

	if d == ev3.Right && pos >= config.VisionThresholdPosition {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else if d == ev3.Left && pos <= -config.VisionThresholdPosition {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else if hasLeftEstimation && (rightIntensity == 0 || hasRightEstimation) && estimationIsOld(estimatedPositionLeft, pos) {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else if hasRightEstimation && (leftIntensity == 0 || hasLeftEstimation) && estimationIsOld(estimatedPositionRight, pos) {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else if (hasLeftEstimation || hasRightEstimation) && leftIntensity == 0 && rightIntensity == 0 {
		dir = switchDirection(pos, leftIntensity, rightIntensity, d)
	} else {
		dir = d

		if leftIntensity > currentIntensityLeft {
			if firstIntensityLeft == 0 {
				firstIntensityLeft = leftIntensity
				firstPositionLeft = pos
			}
			currentIntensityLeft = leftIntensity
			currentPositionLeft = pos
		} else if leftIntensity < currentIntensityLeft-(currentIntensityLeft/config.VisionEstimateReductionRange) {
			estimatedIntensityLeft = currentIntensityLeft
			positionCorrection := computeEstimatedPositionCorrection(abs(firstPositionLeft), firstIntensityLeft, currentPositionLeft, currentIntensityLeft, abs(pos), leftIntensity)
			estimatedPositionLeft = currentPositionLeft - (int(dir) * positionCorrection)
			hasLeftEstimation = true
		}

		if rightIntensity > currentIntensityRight {
			if firstIntensityRight == 0 {
				firstIntensityRight = rightIntensity
				firstPositionRight = pos
			}
			currentIntensityRight = rightIntensity
			currentPositionRight = pos
		} else if rightIntensity < currentIntensityRight-(currentIntensityRight/config.VisionEstimateReductionRange) {
			estimatedIntensityRight = currentIntensityRight
			positionCorrection := computeEstimatedPositionCorrection(abs(firstPositionRight), firstIntensityRight, currentPositionRight, currentIntensityRight, abs(pos), rightIntensity)
			estimatedPositionRight = currentPositionRight - (int(dir) * positionCorrection)
			hasRightEstimation = true
		}

		// intens, ang, _ := estimate(dir)
		// fmt.Fprintln(os.Stderr, "VISION", dir, pos, leftValue, rightValue, "- I", leftIntensity, rightIntensity, "- L", currentIntensityLeft, currentPositionLeft, "- R", currentIntensityRight, currentPositionRight, "- RES", intens, ang)
		// fmt.Fprintln(os.Stderr, "VISION", dir, pos, "- I", leftIntensity, rightIntensity, "- L", currentIntensityLeft, currentPositionLeft, "- R", currentIntensityRight, currentPositionRight, "- RES", intens, ang)
	}

	return estimate(dir)
}
