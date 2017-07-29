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
		if d.VisionAngle > config.TrackFrontAngle {
			go trackTurn(now, ev3.Right)
		} else if d.VisionAngle < -config.TrackFrontAngle {
			go trackTurn(now, ev3.Left)
		} else {
			go trackScan(now)
		}
	}
	return result
}

const trackPrintMillis = 250

func trackScan(start int) {
	now := start
	var dir ev3.Direction = ev3.Right

	fmt.Fprintln(os.Stderr, "TRACK SCAN")
	printTick := 0

	for {
		select {
		case d := <-data:
			now, _ = handleTime(d, start)

			if d.VisionIntensity == 0 {
				go seekTurning(now, dir)
				return
			}
			if d.VisionIntensity < config.VisionIgnoreBorderValue && checkBorder(d, now) {
				return
			}

			if (now / trackPrintMillis) >= printTick {
				printTick = (now / trackPrintMillis) + 1
				fmt.Fprintln(os.Stderr, "TRACK SCAN FRONT", now, d.VisionIntensity, d.VisionAngle)
			}

			speed(config.TrackMaxSpeed, config.TrackMaxSpeed)

			/*
				if d.VisionAngle > config.TrackFrontAngle {
					dir = ev3.Right
					speedCorrectionAngle := config.TrackScanAngle - d.VisionAngle
					speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle
					// speedCorrection = speedCorrection * 300 / 400

					if (now / trackPrintMillis) >= printTick {
						printTick = (now / trackPrintMillis) + 1
						fmt.Fprintln(os.Stderr, "TRACK SCAN RIGHT", d.VisionIntensity, d.VisionAngle, speedCorrection)
					}

					speed(config.TrackOuterSpeed, config.TrackInnerSpeed+speedCorrection)
				} else if d.VisionAngle < -config.TrackFrontAngle {
					dir = ev3.Left
					speedCorrectionAngle := config.TrackScanAngle + d.VisionAngle
					speedCorrection := config.TrackSpeedReductionMax * speedCorrectionAngle / config.TrackSpeedReductionAngle
					// speedCorrection = speedCorrection * 300 / 400

					if (now / trackPrintMillis) >= printTick {
						printTick = (now / trackPrintMillis) + 1
						fmt.Fprintln(os.Stderr, "TRACK SCAN LEFT", d.VisionIntensity, d.VisionAngle, speedCorrection)
					}

					speed(config.TrackInnerSpeed+speedCorrection, config.TrackOuterSpeed)
				} else {

					if (now / trackPrintMillis) >= printTick {
						printTick = (now / trackPrintMillis) + 1
						fmt.Fprintln(os.Stderr, "TRACK SCAN FRONT", d.VisionIntensity, d.VisionAngle)
					}

					speed(config.TrackMaxSpeed, config.TrackMaxSpeed)
				}
			*/

			// fmt.Fprintln(os.Stderr, "TRACK time ", now, ", speed", c.SpeedLeft, c.SpeedRight)
			ledsFromData(d)
			cmdTrackScan()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}

func trackTurn(start int, dir ev3.Direction) {
	now := start

	fmt.Fprintln(os.Stderr, "TRACK TURN", dir)
	printTick := 0

	for {
		select {
		case d := <-data:
			now, _ = handleTime(d, start)

			if d.VisionIntensity == 0 {
				go seekTurning(now, dir)
				return
			}
			if d.VisionIntensity < config.VisionIgnoreBorderValue && checkBorder(d, now) {
				return
			}
			if d.VisionAngle > -config.TrackFrontAngle && d.VisionAngle < config.TrackFrontAngle {
				go trackScan(now)
				return
			}

			if (now / trackPrintMillis) >= printTick {
				printTick = (now / trackPrintMillis) + 1
				fmt.Fprintln(os.Stderr, "TRACK TURN", now, dir, d.VisionIntensity, d.VisionAngle, config.TrackTurnSpeed*ev3.LeftTurnVersor(dir), config.TrackTurnSpeed*ev3.RightTurnVersor(dir))
			}
			speed(config.TrackTurnSpeed*ev3.LeftTurnVersor(dir), config.TrackTurnSpeed*ev3.RightTurnVersor(dir))

			ledsFromData(d)
			cmdTrackTurn(dir)
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}
