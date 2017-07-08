package vision

import (
	"fmt"
	"go-bots/xl4/logic"
	"os"
)

const maxIrValue = 100
const maxIrDistance = 100

var lastSecond = 0
var lastSecondChanged = false

func detectLastSecond(now int, lastSecond int) (int, bool) {
	lastSecondChanged := false
	if now%1000 != lastSecond {
		lastSecond = now % 1000
		lastSecondChanged = true
	}
	return lastSecond, lastSecondChanged
}

// Process processes IR input values into vision data
func Process(millis int, leftValue int, rightValue int) (intensity int, angle int) {
	intensity = 0
	angle = 0
	if leftValue >= maxIrDistance && rightValue >= maxIrDistance {

		lastSecond, lastSecondChanged = detectLastSecond(millis, lastSecond)
		if lastSecondChanged {
			fmt.Fprintln(os.Stderr, "DATA EMPTY")
		}

		return
	}
	intensityLeft := (maxIrValue - leftValue) * logic.VisionIntensityMax / maxIrValue
	intensityRight := (maxIrValue - rightValue) * logic.VisionIntensityMax / maxIrValue
	if intensityRight > intensityLeft {
		intensity = intensityRight
	} else {
		intensity = intensityLeft
	}
	if intensityLeft == 0 {
		angle = (logic.VisionIntensityMax - intensityRight) * logic.VisionAngleMax / logic.VisionIntensityMax
		angle /= 2
		angle += logic.VisionAngleMax / 2
	} else if intensityRight == 0 {
		angle = (logic.VisionIntensityMax - intensityLeft) * logic.VisionAngleMax / logic.VisionIntensityMax
		angle /= 2
		angle += logic.VisionAngleMax / 2
		angle = -angle
	} else {
		angle = (intensityRight - intensityLeft) / 2
	}

	lastSecond, lastSecondChanged = detectLastSecond(millis, lastSecond)
	if lastSecondChanged {
		fmt.Fprintln(os.Stderr, "DATA", leftValue, rightValue, intensity, angle)
	}

	return
}
