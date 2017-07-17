package vision

import "go-bots/xl4/config"

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
	if leftValue >= config.MaxIrDistance && rightValue >= config.MaxIrDistance {

		lastSecond, lastSecondChanged = detectLastSecond(millis, lastSecond)
		if lastSecondChanged {
			// fmt.Fprintln(os.Stderr, "DATA EMPTY")
		}

		return
	}
	intensityLeft := (config.MaxIrValue - leftValue) * config.VisionIntensityMax / config.MaxIrValue
	intensityRight := (config.MaxIrValue - rightValue) * config.VisionIntensityMax / config.MaxIrValue
	if intensityRight > intensityLeft {
		intensity = intensityRight
	} else {
		intensity = intensityLeft
	}
	if intensityLeft == 0 {
		angle = (config.VisionIntensityMax - intensityRight) * config.VisionAngleMax / config.VisionIntensityMax
		angle /= 2
		angle += config.VisionAngleMax / 2
	} else if intensityRight == 0 {
		angle = (config.VisionIntensityMax - intensityLeft) * config.VisionAngleMax / config.VisionIntensityMax
		angle /= 2
		angle += config.VisionAngleMax / 2
		angle = -angle
	} else {
		angle = (intensityRight - intensityLeft) / 2
	}

	lastSecond, lastSecondChanged = detectLastSecond(millis, lastSecond)
	if lastSecondChanged {
		// fmt.Fprintln(os.Stderr, "DATA", leftValue, rightValue, intensity, angle)
	}

	return
}
