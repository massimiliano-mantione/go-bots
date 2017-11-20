package data

import (
	"fmt"
)

// Point describes measures at one point in time
type Point struct {
	T      uint32
	DT     uint16
	Pos    int16
	PosD   int16
	PosI   int16
	PosE   int16
	FP     int16
	FD     int16
	FI     int16
	FE     int16
	SpeedL int16
	SpeedR int16
	Kind   uint8
	Left   int8
	Right  int8
}

var pointKinds [16]string

// Init initializes the data store
func Init(pk [16]string) {
	pointKinds = pk
	Reset()
}

const maxPoints = 20000

var startIndex int
var nextIndex int
var points = [maxPoints]Point{}

// Reset resets the data
func Reset() {
	fmt.Println("Data reset start")
	for i := 0; i < maxPoints; i++ {
		points[i] = Point{}
	}
	startIndex = 0
	nextIndex = 1
	fmt.Println("Data reset done")
}

// Store stores a data point
func Store(time uint32,
	dt uint16,
	pos int16,
	posD int16,
	posI int16,
	posE int16,
	fP int16,
	fD int16,
	fI int16,
	fE int16,
	kind uint8,
	speedL int16,
	speedR int16,
	left int8,
	right int8) {

	points[nextIndex] = Point{
		T:      time / 100,
		DT:     dt / 100,
		Pos:    pos,
		PosD:   posD,
		PosI:   posI,
		PosE:   posE,
		FP:     fP,
		FD:     fD,
		FI:     fI,
		FE:     fE,
		Kind:   kind,
		SpeedL: speedL,
		SpeedR: speedR,
		Left:   left,
		Right:  right,
	}
	nextIndex++
	if nextIndex >= maxPoints {
		nextIndex = 0
	}
	if startIndex == nextIndex {
		startIndex++
	}
}

func printPoint(p *Point) {
	fmt.Printf("%6d.%d %4d.%d %s   P %5d D %5d I %5d E %5d   F %4d %4d %4d %4d   Sl %4d Sr %4d L %4d R %4d\n",
		p.T/10,
		p.T%10,
		p.DT/10,
		p.DT%10,
		pointKinds[p.Kind],
		p.Pos,
		p.PosD,
		p.PosI,
		p.PosE,
		p.FP,
		p.FD,
		p.FI,
		p.FE,
		p.SpeedL,
		p.SpeedR,
		p.Left,
		p.Right)
}

// Print prints the data
func Print() {
	fmt.Println()
	fmt.Println("Data points")
	for i := startIndex; i < maxPoints; i++ {
		if i == nextIndex {
			fmt.Println()
			return
		}
		printPoint(&points[i])
	}
	for i := 0; i < nextIndex; i++ {
		printPoint(&points[i])
	}
	fmt.Println()
}
