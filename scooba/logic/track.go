package logic

import (
	"fmt"
	"os"
)

func checkVision(d Data, now int) bool {
	if d.IrValueLeft < 100 || d.IrValueFrontLeft < 100 || d.IrValueFrontRight < 100 || d.IrValueRight < 100 {
		// REMOVE ME!!
		// go track(now)
		return true
	}
	return false
}

const trackPrintMillis = 250

func track(start int) {
	// now, elapsed := start, 0
	// var dir ev3.Direction = ev3.Right

	fmt.Fprintln(os.Stderr, "TRACK")

	for {
		select {
		case d := <-data:
			_, _ = handleTime(d, start)

			if d.IrValueLeft >= 100 || d.IrValueFrontLeft >= 100 || d.IrValueFrontRight >= 100 || d.IrValueRight >= 100 {
				return
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
