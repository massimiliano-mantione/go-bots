package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/scooba/config"
	"go-bots/ui"
	"os"
)

func pauseBeforeBegin(start int, strategy func(int, ev3.Direction), dir ev3.Direction) {
	for {
		select {
		case d := <-data:
			now, elapsed := handleTime(d, start)
			if elapsed >= config.StartTime {
				go strategy(now, dir)
				return
			}
			speed(0, 0)
			intensity := ((elapsed % 1000) * 255) / (config.StartTime / 5)
			if elapsed > (config.StartTime * 4 / 5) {
				leds(intensity, intensity, intensity, intensity)
			} else {
				leds(0, 0, intensity, intensity)
			}
			startCmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}

func chooseStrategy(start int) {
	strategy := goForward
	var dir ev3.Direction = ev3.Left
	leds(0, 0, 0, 0)
	speed(0, 0)
	startCmd()
	fmt.Fprintln(os.Stderr, "chooseStrategy START")

	for {
		select {
		case d := <-data:
			handleTime(d, start)
			speed(0, 0)
			startCmd()
		case k := <-keys:
			checkQuit(k)
			if k.Key == ui.Enter {
				go pauseBeforeBegin(k.Millis, strategy, dir)
				return
			} else if k.Key == ui.Left {
				dir = ev3.Left
				leds(255, 0, 255, 0)
				fmt.Fprintln(os.Stderr, "chooseStrategy circle left")
			} else if k.Key == ui.Right {
				dir = ev3.Right
				leds(0, 255, 0, 255)
				fmt.Fprintln(os.Stderr, "chooseStrategy circle right")
			} else if k.Key == ui.Up {
				strategy = goForward
				if dir == ev3.Left {
					leds(255, 0, 0, 0)
					fmt.Fprintln(os.Stderr, "chooseStrategy forward left")
				} else {
					leds(0, 255, 0, 0)
					fmt.Fprintln(os.Stderr, "chooseStrategy forward right")
				}
			} else if k.Key == ui.Down {
				strategy = turnBack
				if dir == ev3.Left {
					leds(0, 0, 255, 0)
					fmt.Fprintln(os.Stderr, "chooseStrategy back left")
				} else {
					leds(0, 0, 0, 255)
					fmt.Fprintln(os.Stderr, "chooseStrategy back right")
				}
			}
			speed(0, 0)
			startCmd()
		}
	}
}

func goForward(start int, dir ev3.Direction) {
	now, elapsed := start, 0

	fmt.Fprintln(os.Stderr, "goForward", now, dir)

	for elapsed < config.GoForwardMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			speed(config.GoForwardSpeed, config.GoForwardSpeed)

			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}

	fmt.Fprintln(os.Stderr, "goForward turn", now, dir)

	forwardDone := elapsed
	for elapsed-forwardDone < config.GoForwardTurnMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			if dir == ev3.Right {
				speed(config.GoForwardTurnOuterSpeed, config.GoForwardTurnInnerSpeed)
			} else {
				speed(config.GoForwardTurnInnerSpeed, config.GoForwardTurnOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}

func turnBack(start int, dir ev3.Direction) {
	now, elapsed := start, 0

	fmt.Fprintln(os.Stderr, "turnBack", now, dir)

	for elapsed < config.TurnBackPreMoveMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			speed(config.TurnBackPreMoveSpeed, config.TurnBackPreMoveSpeed)

			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}

	fmt.Fprintln(os.Stderr, "turnBack turn", now, dir)
	turnBackPreMoveDone := elapsed
	for elapsed-turnBackPreMoveDone < config.TurnBackMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			if dir == ev3.Right {
				speed(config.TurnBackOuterSpeed, config.TurnBackInnerSpeed)
			} else {
				speed(config.TurnBackInnerSpeed, config.TurnBackOuterSpeed)
			}

			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}

	fmt.Fprintln(os.Stderr, "turnBack move", now, dir)

	turnBackMoveDone := elapsed
	for elapsed-turnBackMoveDone < config.TurnBackMoveMillis {
		select {
		case d := <-data:
			now, elapsed = handleTime(d, start)

			if checkVision(d, now) {
				return
			}

			speed(config.TurnBackMoveSpeed, config.TurnBackMoveSpeed)

			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return
			}
		}
	}
}
