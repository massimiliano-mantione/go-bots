package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/seeker2/config"
	"os"
)

func seekMove(start int, dir ev3.Direction, leftSpeed int, rightSpeed int, duration int, ignoreBorder bool) (done bool, now int) {

	fmt.Fprintln(os.Stderr, "seekMove", dir, leftSpeed, rightSpeed, duration)

	for {
		select {
		case d := <-data:
			now, elapsed := handleTime(d, start)
			if elapsed >= duration {

				fmt.Fprintln(os.Stderr, "seekMove elapsed", now, duration)

				return false, now
			}

			if checkVision(d, now) {
				return true, now
			}
			if (!ignoreBorder) && checkBorder(d, now) {
				return true, now
			}

			speed(leftSpeed, rightSpeed)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			checkEnd(k)
		}
	}
}

func back(start int, dir ev3.Direction) {

	done, now := false, start
	fmt.Fprintln(os.Stderr, "BACK", dir)

	if dir == ev3.Right {
		done, now = seekMove(now, dir, -config.BackTurn1SpeedInner, -config.BackTurn1SpeedOuter, config.BackTurn1Millis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, config.BackTurn2Speed, -config.BackTurn2Speed, config.BackTurn2Millis, true)
		if done {
			return
		}
	} else if dir == ev3.Left {
		done, now = seekMove(now, dir, -config.BackTurn1SpeedOuter, -config.BackTurn1SpeedInner, config.BackTurn1Millis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, -config.BackTurn2Speed, config.BackTurn2Speed, config.BackTurn2Millis, true)
		if done {
			return
		}
	} else {
		dir = ev3.Right
		done, now = seekMove(now, dir, -config.BackMoveSpeed, -config.BackMoveSpeed, config.BackMoveMillis, true)
		if done {
			return
		}
		done, now = seekMove(now, dir, config.BackTurn3Speed, -config.BackTurn3Speed, config.BackTurn3Millis, true)
		if done {
			return
		}
	}

	go seek(now, dir)
}

func seek(start int, dir ev3.Direction) {

	fmt.Fprintln(os.Stderr, "SEEK", dir)

	done, now := false, start
	for {

		fmt.Fprintln(os.Stderr, "SEEK MOVE", dir, now)

		done, now = seekMove(now, dir, config.SeekMoveSpeed, config.SeekMoveSpeed, config.SeekMoveMillis, false)
		if done {
			return
		}

		fmt.Fprintln(os.Stderr, "SEEK TURN", dir, now)

		done, now = seekMove(now, dir, config.SeekTurnSpeed*ev3.LeftTurnVersor(dir), config.SeekTurnSpeed*ev3.RightTurnVersor(dir), config.SeekTurnMillis, false)
		if done {
			return
		}
		dir = ev3.ChangeDirection(dir)
	}
}
