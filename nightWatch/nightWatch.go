package main

import (
	"go-bots/nightWatch/io"
	"go-bots/nightWatch/logic"
	"go-bots/ui"
)

var data = make(chan logic.Data)
var keys = make(chan ui.KeyEvent)
var quit = make(chan bool)

func main() {
	io.Init(data)
	defer io.Close()
	go io.Loop()

	ui.Init(keys, io.StartTime())
	defer ui.Close()
	go ui.Loop()

	logic.Init(data, io.ProcessCommand, keys, quit)
	go logic.Run()
	<-quit
}
