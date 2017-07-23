package logic

import (
	"fmt"
	"go-bots/ev3"
	"os"
)

func checkVision(d Data, now int) bool {
	result := false
	if result {
		fmt.Fprintln(os.Stderr, "CheckVision")
		go track(now)
	}
	return result
}

func track(start int) {
	now, _ := start, 0
	var dir ev3.Direction = ev3.Right

	for {
		select {
		case d := <-data:
			now, _ = handleTime(d, start)

			if true {
				go seek(now, dir)
				return
			}
			if checkBorder(d, now) {
				return
			}

			/*
				if d.IrLeftValue >= config.MaxIrDistance {
					speed(config.TrackOnly1SensorOuterSpeed, config.TrackOnly1SensorInnerSpeed)
					dir = ev3.Right
				} else if d.IrRightValue >= config.MaxIrDistance {
					speed(config.TrackOnly1SensorInnerSpeed, config.TrackOnly1SensorOuterSpeed)
					dir = ev3.Left
				} else {
					difference := d.IrLeftValue - d.IrRightValue

					if difference > config.TrackCenterZone {
						speed(config.TrackSpeed, config.TrackSpeed-(difference*config.TrackDifferenceCoefficent))
						dir = ev3.Right
					} else if difference < -config.TrackCenterZone {
						speed(config.TrackSpeed+(difference*config.TrackDifferenceCoefficent), config.TrackSpeed)
						dir = ev3.Left
					} else {
						speed(config.TrackSpeed, config.TrackSpeed)
					}

				}
			*/
			// fmt.Fprintln(os.Stderr, "TRACK time ", now, ", speed", c.SpeedLeft, c.SpeedRight, ", IRsensors", d.IrLeftValue, d.IrRightValue)
			ledsFromData(d)
			cmd(true, true)
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}
