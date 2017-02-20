package vision

import (
	"log"
)

func Process(direction int, rightValue int, leftValue int) {
	log.Println("vision.Process", direction, rightValue, leftValue)
}

func Estimate() (intensity int, angle int) {
	intensity = 0
	angle = 0
	return
}
