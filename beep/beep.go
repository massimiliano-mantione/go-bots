package beep

import (
	"os/exec"
)

func play(args ...string) {
	cmd := exec.Command("beep", args...)
	cmd.Start()
	go func() {
		cmd.Wait()
	}()
}

// Beep beeps
func Beep() {
	play()
}

// C beeps
func C() {
	play("-f", "261.6")
}

// G beeps
func G() {
	play("-f", "392.0")
}

// CG beeps
func CG() {
	play("-f", "261.6", "-n", "-f", "392.0")
}

// GC beeps
func GC() {
	play("-f", "392.0", "-n", "-f", "261.6")
}

// GG beeps
func GG() {
	play("-f", "392.0", "-n", "-f", "392.0")
}

// CCC beeps
func CCC() {
	play("-f", "261.6", "-n", "-f", "261.6", "-n", "-f", "261.6")
}
