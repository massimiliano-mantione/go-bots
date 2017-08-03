package logic

import (
	"fmt"
	"go-bots/ev3"
	"go-bots/scooba/config"
	"os"
)

func seekMove(start int, dir ev3.Direction, leftSpeed int, rightSpeed int, duration int, ignoreBorder bool) (done bool, now int) {
	for {
		select {
		case d := <-data:
			now, elapsed := handleTime(d, start)
			if elapsed >= duration {
				return false, now
			}

			if checkVision(d, now) {
				return true, now
			}

			speed(leftSpeed, rightSpeed)
			ledsFromData(d)
			cmd()
		case k := <-keys:
			if checkDone(k) {
				return true, now
			}
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
		fmt.Fprintln(os.Stderr, "BACK TURN", dir)
		done, now = seekMove(now, dir, config.BackTurn2Speed, -config.BackTurn2Speed, config.BackTurn2Millis, false)
		if done {
			return
		}
	} else if dir == ev3.Left {
		done, now = seekMove(now, dir, -config.BackTurn1SpeedOuter, -config.BackTurn1SpeedInner, config.BackTurn1Millis, true)
		if done {
			return
		}
		fmt.Fprintln(os.Stderr, "BACK TURN", dir)
		done, now = seekMove(now, dir, -config.BackTurn2Speed, config.BackTurn2Speed, config.BackTurn2Millis, false)
		if done {
			return
		}
	} else {
		dir = ev3.Right
		done, now = seekMove(now, dir, -config.BackMoveSpeed, -config.BackMoveSpeed, config.BackMoveMillis, true)
		if done {
			return
		}
		fmt.Fprintln(os.Stderr, "BACK TURN", dir)
		done, now = seekMove(now, dir, config.BackTurn3Speed, -config.BackTurn3Speed, config.BackTurn3Millis, false)
		if done {
			return
		}
	}

	go seekMoving(now, dir)
}

func seekMoving(start int, dir ev3.Direction) {
	seek(start, dir, false)
}
func seekTurning(start int, dir ev3.Direction) {
	seek(start, dir, true)
}

func seek(start int, dir ev3.Direction, skipFirstMove bool) {

	fmt.Fprintln(os.Stderr, "SEEK", dir)

	done, now := false, start
	for {

		if !skipFirstMove {
			fmt.Fprintln(os.Stderr, "SEEK MOVE", dir, now)
			done, now = seekMove(now, dir, config.SeekMoveSpeed, config.SeekMoveSpeed, config.SeekMoveMillis, false)
			if done {
				return
			}
		} else {
			skipFirstMove = false
		}

		fmt.Fprintln(os.Stderr, "SEEK TURN", dir, now)
		done, now = seekMove(now, dir, config.SeekTurnSpeed*ev3.LeftTurnVersor(dir), config.SeekTurnSpeed*ev3.RightTurnVersor(dir), config.SeekTurnMillis, false)
		if done {
			return
		}
		dir = ev3.ChangeDirection(dir)
	}
}
