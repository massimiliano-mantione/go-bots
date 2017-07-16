package main

import (
	"go-bots/ui"
	"go-bots/xl4/io"
	"go-bots/xl4/logic"
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
