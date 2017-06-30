package vision

import "go-bots/xl4/logic"

const maxIrValue = 100
const maxIrDistance = 100

// Process processes IR input values into vision data
func Process(millis int, leftValue int, rightValue int) (intensity int, angle int) {
	intensity = 0
	angle = 0
	if leftValue >= maxIrDistance && rightValue >= maxIrDistance {
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
	return
}
