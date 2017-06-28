package main

import (
	"go-bots/ui"
	"go-bots/xl4/io"
	"go-bots/xl4/logic"
)

var data chan logic.Data = make(chan logic.Data)
var keys chan ui.Key = make(chan ui.Key)

func main() {
	io.Init(data)
	defer io.Close()
	go io.Loop()

	ui.Init(keys)
	defer ui.Close()
	go ui.Loop()

	logic.Init(data, io.ProcessCommand, keys)
	logic.Run()
}
