package vision

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"os"
)

var maxIntensityLeft int

func hasSeenLeft() bool {
	return maxIntensityLeft > 0
}

var maxIntensityRight int

func hasSeenRight() bool {
	return maxIntensityRight > 0
}

type Point struct {
	Millis    int
	Intensity int
	Angle     int
}

var point0 Point
var point1 Point
var point2 Point
var point3 Point
var point4 Point
var point5 Point
var point6 Point
var point7 Point
var point8 Point

func recordPoint(millis int, intensity int, angle int) {
	point8 = point7
	point7 = point6
	point6 = point5
	point5 = point4
	point4 = point3
	point3 = point2
	point2 = point1
	point1 = point0
	point0.Millis = millis
	point0.Intensity = intensity
	point0.Angle = angle
}

func contributePoint(now int, intensitySum int, angleSum int, ic int, ac int, point Point) (nextIntensitySum int, nextAngleSum int, intensityCount int, angleCount int) {
	timeAttenuation := 1 + ((now - point.Millis) / 10)
	nextIntensitySum = intensitySum + (point.Intensity / timeAttenuation)
	intensityCount = ic + 1
	if point.Intensity > 0 {
		nextAngleSum = angleSum + (point.Angle / timeAttenuation)
		angleCount = ac + 1
	} else {
		nextAngleSum = angleSum
		angleCount = ac
	}
	return
}

func estimate(now int, d ev3.Direction) (intensity int, angle int, dir ev3.Direction) {
	intensity, angle, intensityCount, angleCount := 0, 0, 0, 0
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point0)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point1)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point2)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point3)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point4)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point5)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point6)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point7)
	intensity, angle, intensityCount, angleCount = contributePoint(now, intensity, angle, intensityCount, angleCount, point8)
	intensity /= intensityCount
	if angleCount > 0 {
		angle /= angleCount
	} else {
		angle = 0
	}

	dir = d
	return
}

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

func positionToAngle(pos int, dir ev3.Direction) int {
	return (pos * 9 / 25) + (int(dir) * 45)
}

// Reset resets the vision state
func Reset() {
	maxIntensityLeft = 0
	maxIntensityRight = 0
}

func switchDirection(dir ev3.Direction) ev3.Direction {
	maxIntensityLeft = 0
	maxIntensityRight = 0
	return ev3.ChangeDirection(dir)
}

// Process processes IR sensor data
func Process(now int, d ev3.Direction, pos int, leftValue int, rightValue int) (intensity int, angle int, dir ev3.Direction) {
	leftIntensity, rightIntensity := irValuesToIntensity(leftValue, rightValue, pos)
	if leftIntensity > maxIntensityLeft {
		maxIntensityLeft = leftIntensity
	}
	if rightIntensity > maxIntensityRight {
		maxIntensityRight = rightIntensity
	}
	recordPoint(now, leftIntensity, positionToAngle(pos, ev3.Left))
	recordPoint(now, rightIntensity, positionToAngle(pos, ev3.Right))

	if d == ev3.Right && pos >= config.VisionThresholdPosition {
		fmt.Fprintln(os.Stderr, "VISION SWITCH RIGHT")
		dir = switchDirection(d)
	} else if d == ev3.Left && pos <= -config.VisionThresholdPosition {
		fmt.Fprintln(os.Stderr, "VISION SWITCH LEFT")
		dir = switchDirection(d)
	} else if hasSeenLeft() && (leftIntensity == 0 || leftIntensity <= maxIntensityLeft-(maxIntensityLeft/config.VisionEstimateReductionRange)) {
		dir = switchDirection(d)
	} else if hasSeenRight() && (rightIntensity == 0 || rightIntensity <= maxIntensityRight-(maxIntensityRight/config.VisionEstimateReductionRange)) {
		dir = switchDirection(d)
	} else {
		dir = d
	}

	intens, ang, _ := estimate(now, dir)
	// fmt.Fprintln(os.Stderr, "VISION", dir, pos, leftValue, rightValue, "- I", leftIntensity, rightIntensity, "- L", currentIntensityLeft, currentPositionLeft, "- R", currentIntensityRight, currentPositionRight, "- RES", intens, ang)
	fmt.Fprintln(os.Stderr, "VISION", dir, pos, "- I", leftIntensity, rightIntensity, "-P", hasSeenLeft(), hasSeenRight(), "- RES", intens, ang)

	return estimate(now, dir)
}
