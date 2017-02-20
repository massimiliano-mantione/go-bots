package main

import (
	"go-bots/nightWatch/io"
	"go-bots/nightWatch/logic"
)

var data chan logic.Data = make(chan logic.Data)
var commands chan logic.Commands = make(chan logic.Commands)

func main() {
	io.Init(data, commands)
	defer io.Close()

	logic.Init(data, commands)

	go io.Loop()
	logic.Run()
}
