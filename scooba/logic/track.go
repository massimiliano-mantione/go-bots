package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/scooba/config"
	"os"
)

func checkVision(d Data, now int) bool {
	result := d.VisionIntensity > 0
	if result {
		go track(now)
	}
	return result
}

const trackPrintMillis = 250

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
				go seekTurning(now, dir)
				return
			}

			if d.VisionAngle > config.TrackSemiFrontAngle {
				dir = ev3.Right
				speedCorrectionAngle := config.VisionMaxAngle - d.VisionAngle
				speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle
				// speedCorrection = speedCorrection * 200 / 400

				if (now / trackPrintMillis) >= printTick {
					printTick = (now / trackPrintMillis) + 1
					fmt.Fprintln(os.Stderr, "TRACK RIGHT", d.VisionIntensity, d.VisionAngle, speedCorrection)
				}

				speed(config.TrackOuterSpeed, config.TrackInnerSpeed+speedCorrection)
			} else if d.VisionAngle > config.TrackFrontAngle {
				if (now / trackPrintMillis) >= printTick {
					printTick = (now / trackPrintMillis) + 1
					fmt.Fprintln(os.Stderr, "TRACK FRONT RIGHT", d.VisionIntensity, d.VisionAngle)
				}
				speed(config.TrackOuterSpeed, config.TrackSemiFrontInnerSpeed)
			} else if d.VisionAngle < -config.TrackSemiFrontAngle {
				dir = ev3.Left
				speedCorrectionAngle := config.VisionMaxAngle + d.VisionAngle
				speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle
				// speedCorrection = speedCorrection * 200 / 400

				if (now / trackPrintMillis) >= printTick {
					printTick = (now / trackPrintMillis) + 1
					fmt.Fprintln(os.Stderr, "TRACK LEFT", d.VisionIntensity, d.VisionAngle, speedCorrection)
				}

				speed(config.TrackInnerSpeed+speedCorrection, config.TrackOuterSpeed)
			} else if d.VisionAngle < -config.TrackFrontAngle {
				if (now / trackPrintMillis) >= printTick {
					printTick = (now / trackPrintMillis) + 1
					fmt.Fprintln(os.Stderr, "TRACK FRONT LEFT", d.VisionIntensity, d.VisionAngle)
				}
				speed(config.TrackSemiFrontInnerSpeed, config.TrackOuterSpeed)
			} else {

				if (now / trackPrintMillis) >= printTick {
					printTick = (now / trackPrintMillis) + 1
					fmt.Fprintln(os.Stderr, "TRACK FRONT", d.VisionIntensity, d.VisionAngle)
				}

				speed(config.TrackMaxSpeed, config.TrackMaxSpeed)
			}
			// fmt.Fprintln(os.Stderr, "TRACK time ", now, ", speed", c.SpeedLeft, c.SpeedRight)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}
