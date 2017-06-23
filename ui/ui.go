package ui

import t "github.com/gizak/termui"
import "log"
import "golang.org/x/crypto/ssh/terminal"
import "time"

type Key int

const (
	None Key = iota
	Enter
	Back
	Up
	Down
	Right
	Left
	Quit
)

var keys chan<- Key
var state *terminal.State

var lastEnterTime time.Time

func Init(k chan<- Key) {
	var err error
	keys = k
	state, err = terminal.GetState(0)
	if err != nil {
		log.Fatalln("Error getting terminal state:", err)
	}
}

const quitEvent = "custom/quitEvent"

func Close() {
	log.Println("Closing terminal")
	if state != nil {
		terminal.Restore(0, state)
	}
	t.SendCustomEvt(quitEvent, nil)
}

func Loop() {
	err := t.Init()
	if err != nil {
		log.Fatalln("Error setting up ui:", err)
	}
	defer t.Close()

	t.Handle(quitEvent, func(t.Event) {
		keys <- Quit
		t.StopLoop()
	})
	t.Handle("/sys/kbd/C-c", func(t.Event) {
		keys <- Quit
		t.StopLoop()
	})
	t.Handle("/sys/kbd/<up>", func(t.Event) {
		keys <- Up
	})
	t.Handle("/sys/kbd/<down>", func(t.Event) {
		keys <- Down
	})
	t.Handle("/sys/kbd/<right>", func(t.Event) {
		keys <- Right
	})
	t.Handle("/sys/kbd/<left>", func(t.Event) {
		keys <- Left
	})
	t.Handle("/sys/kbd/<enter>", func(t.Event) {
		lastEnterTime = time.Now()
		keys <- Enter
	})
	t.Handle("/sys/kbd/C-8", func(t.Event) {
		backTime := time.Now()
		interval := backTime.Sub(lastEnterTime)
		if interval < time.Millisecond*400 {
			keys <- Quit
		} else {
			keys <- Back
		}
	})

	t.Loop()
}
