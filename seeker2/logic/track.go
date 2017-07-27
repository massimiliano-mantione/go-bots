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
			if d.VisionIntensity < config.VisionIgnoreBorderValue && checkBorder(d, now) {
				return
			}

			if d.VisionAngle > config.TrackFrontAngle {
				speedCorrectionAngle := config.VisionMaxAngle - d.VisionAngle
				speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle

				if (now / 500) >= printTick {
					printTick = (now / 500) + 1
					fmt.Fprintln(os.Stderr, "TRACK RIGHT", d.VisionIntensity, d.VisionAngle, speedCorrection)
				}

				speed(config.TrackOuterSpeed, config.TrackInnerSpeed+speedCorrection)
			} else if d.VisionAngle < -config.TrackFrontAngle {
				speedCorrectionAngle := config.VisionMaxAngle + d.VisionAngle
				speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle

				if (now / 500) >= printTick {
					printTick = (now / 500) + 1
					fmt.Fprintln(os.Stderr, "TRACK LEFT", d.VisionIntensity, d.VisionAngle, speedCorrection)
				}

				speed(config.TrackInnerSpeed+speedCorrection, config.TrackOuterSpeed)
			} else {

				if (now / 500) >= printTick {
					printTick = (now / 500) + 1
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
