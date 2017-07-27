package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"os"
)

func checkVision(d Data, now int) bool {
	result := d.VisionIntensity > 0
	if result {
		go track(now)
	}
	return result
}

func track(start int) {
	now, _ := start, 0
	var dir ev3.Direction = ev3.Right

	fmt.Fprintln(os.Stderr, "TRACK")
	printTick := 0

	for {
		select {
		case d := <-data:
			now, _ = handleTime(d, start)

			if d.VisionIntensity == 0 {
				go seek(now, dir)
				return
			}
			if checkBorder(d, now) {
				return
			}

			if d.VisionAngle > config.TrackFrontAngle {
				// speedReductionAngle := d.VisionAngle - config.TrackFrontAngle
				// speedReduction := config.TrackSpeedReductionMax * speedReductionAngle / config.TrackSpeedReductionAngle
				speedReduction := d.VisionAngle * 66

				if (now / 1000) >= printTick {
					printTick = (now / 1000) + 1
					fmt.Fprintln(os.Stderr, "TRACK RIGHT", d.VisionIntensity, d.VisionAngle, speedReduction)
				}

				speed(config.TrackMaxSpeed, config.TrackMaxSpeed-speedReduction)
			} else if d.VisionAngle < -config.TrackFrontAngle {
				// speedReductionAngle := config.TrackFrontAngle - d.VisionAngle
				// speedReduction := config.TrackSpeedReductionMax * speedReductionAngle / config.TrackSpeedReductionAngle
				speedReduction := -d.VisionAngle * 66

				if (now / 1000) >= printTick {
					printTick = (now / 1000) + 1
					fmt.Fprintln(os.Stderr, "TRACK LEFT", d.VisionIntensity, d.VisionAngle, speedReduction)
				}

				speed(config.TrackMaxSpeed-speedReduction, config.TrackMaxSpeed)
			} else {

				if (now / 1000) >= printTick {
					printTick = (now / 1000) + 1
					fmt.Fprintln(os.Stderr, "TRACK FRONT", d.VisionIntensity, d.VisionAngle)
				}

				speed(config.TrackMaxSpeed, config.TrackMaxSpeed)
			}
			// fmt.Fprintln(os.Stderr, "TRACK time ", now, ", speed", c.SpeedLeft, c.SpeedRight)
			ledsFromData(d)
			cmd(true, true)
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}
