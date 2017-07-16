package main

import (
	"go-bots/nightWatch/io"
	"go-bots/nightWatch/logic"
	"go-bots/ui"
	"time"
)

var data = make(chan logic.Data)
var keys = make(chan ui.KeyEvent)
var quit = make(chan bool)

func main() {
	start := time.Now()

	io.Init(data, start)
	defer io.Close()
	go io.Loop()

	ui.Init(keys, start)
	defer ui.Close()
	go ui.Loop()

	logic.Init(data, io.ProcessCommand, keys, quit)
	go logic.Run()
	<-quit
}
