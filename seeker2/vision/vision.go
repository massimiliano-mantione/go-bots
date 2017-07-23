package vision

import (
	"go-bots/ev3"
	"go-bots/seeker2/config"
)

// Process processes IR sensor data
func Process(millis int, d ev3.Direction, pos int, rightValue int, leftValue int) (intensity int, angle int, dir ev3.Direction) {
	intensity, angle = 0, 0

	if d == ev3.Right && pos >= config.VisionThresholdPosition {
		dir = ev3.Left
	} else if d == ev3.Left && pos <= -config.VisionThresholdPosition {
		dir = ev3.Right
	} else {
		dir = d
	}

	// fmt.Fprintln(os.Stderr, " - VISION PROCESS", d, pos, rightValue, leftValue)
	// fmt.Fprintln(os.Stderr, " - VISION  RESULT", intensity, angle, dir)

	return
}
