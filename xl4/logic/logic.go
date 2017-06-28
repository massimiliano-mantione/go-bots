package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/ui"
	"log"
	"os"
	"time"
)

type Data struct {
	Start            time.Time
	Millis           int
	CornerRightIsOut bool
	CornerLeftIsOut  bool
	CornerRight      int
	CornerLeft       int
	VisionIntensity  int
	VisionAngle      int
}

type Commands struct {
	Millis        int
	SpeedRight    int
	SpeedLeft     int
	LedRightRed   int
	LedRightGreen int
	LedLeftRed    int
	LedLeftGreen  int
}

var data <-chan Data
var commandProcessor func(*Commands)
var keys <-chan ui.Key

func Init(d <-chan Data, c func(*Commands), k <-chan ui.Key) {
	data = d
	commandProcessor = c
	keys = k
}

func Run() {
	c := Commands{}
	const fastSpeed = 10000
	const backSpeed = -10000

	for {
		select {
		case d := <-data:
			c.Millis = d.Millis
			if d.CornerRightIsOut {
				if c.SpeedRight > 0 {
					c.SpeedRight = backSpeed
				}
				if c.SpeedLeft > 0 {
					c.SpeedLeft = backSpeed
				}
				c.LedRightRed = 255
			} else if d.CornerLeftIsOut {
				if c.SpeedRight > 0 {
					c.SpeedRight = backSpeed
				}
				if c.SpeedLeft > 0 {
					c.SpeedLeft = backSpeed
				}
				c.LedLeftRed = 255
			} else {
				c.LedRightRed = 0
				c.LedLeftRed = 0
			}

			now := time.Now()
			// millis := ev3.TimespanAsMillis(start, t)
			millis := ev3.TimespanAsMillis(d.Start, now)
			fmt.Fprintln(os.Stderr, "DATA", c.Millis, millis, d.CornerLeftIsOut, d.CornerLeft, d.CornerRightIsOut, d.CornerRight, c.SpeedLeft, c.SpeedRight)

			commandProcessor(&c)
		case k := <-keys:
			switch k {
			case ui.Up:
				c.SpeedRight = fastSpeed
				c.SpeedLeft = fastSpeed
			case ui.Down:
				c.SpeedRight = -fastSpeed
				c.SpeedLeft = -fastSpeed
			case ui.Right:
				c.SpeedRight = -fastSpeed
				c.SpeedLeft = fastSpeed
			case ui.Left:
				c.SpeedRight = fastSpeed
				c.SpeedLeft = -fastSpeed
			case ui.Enter:
				c.SpeedRight = 0
				c.SpeedLeft = 0
			case ui.Quit:
				log.Println("Logic got clean quit (kill)")
				return
			case ui.Back:
				log.Println("Logic got clean quit (back)")
				return
			}
		}
	}
}
