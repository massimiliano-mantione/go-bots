package main

import (
	"go-bots/ui"
	"go-bots/xl4/io"
	"go-bots/xl4/logic"
)

var data = make(chan logic.Data)
var keys = make(chan ui.Key)
var quit = make(chan bool)

func main() {
	io.Init(data)
	defer io.Close()
	go io.Loop()

	ui.Init(keys)
	defer ui.Close()
	go ui.Loop()

	logic.Init(data, io.ProcessCommand, keys, quit)
	go logic.Run()
	<-quit
}
