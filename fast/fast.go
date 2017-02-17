package main

import (
	"log"

	ev3 "go-bots/ev3"
	"time"
)

var bot = ev3.Scan()

func main() {
	log.Printf("Devices: %#v\n", bot)

	ev3.CheckDriver(bot.In1, ev3.DriverIr)
	ev3.CheckDriver(bot.In2, ev3.DriverIr)
	ev3.CheckDriver(bot.In3, ev3.DriverColor)
	ev3.CheckDriver(bot.In4, ev3.DriverColor)

	ev3.CheckDriver(bot.OutA, ev3.DriverTachoMotorLarge)
	ev3.CheckDriver(bot.OutB, ev3.DriverTachoMotorLarge)
	ev3.CheckDriver(bot.OutC, ev3.DriverTachoMotorLarge)
	ev3.CheckDriver(bot.OutD, ev3.DriverTachoMotorMedium)

	ev3.SetMode(bot.In1, ev3.IrModeProx)
	ev3.SetMode(bot.In2, ev3.IrModeProx)
	ev3.SetMode(bot.In3, ev3.ColorModeReflect)
	ev3.SetMode(bot.In4, ev3.ColorModeReflect)

	ev3.RunCommand(bot.OutA, ev3.CmdReset)
	ev3.RunCommand(bot.OutB, ev3.CmdReset)
	ev3.RunCommand(bot.OutC, ev3.CmdReset)
	ev3.RunCommand(bot.OutD, ev3.CmdReset)

	in1 := ev3.OpenByteR(bot.In1, ev3.BinData)
	defer in1.Close()
	in2 := ev3.OpenByteR(bot.In2, ev3.BinData)
	defer in2.Close()
	in3 := ev3.OpenByteR(bot.In3, ev3.BinData)
	defer in3.Close()
	in4 := ev3.OpenByteR(bot.In4, ev3.BinData)
	defer in4.Close()
	ma := ev3.OpenTextW(bot.OutA, ev3.DutyCycleSp)
	defer ma.Close()
	mb := ev3.OpenTextW(bot.OutB, ev3.DutyCycleSp)
	defer mb.Close()
	mc := ev3.OpenTextW(bot.OutC, ev3.DutyCycleSp)
	defer mc.Close()
	md := ev3.OpenTextW(bot.OutD, ev3.DutyCycleSp)
	defer md.Close()
	pd := ev3.OpenTextR(bot.OutD, ev3.Position)
	defer pd.Close()

	ev3.RunCommand(bot.OutA, ev3.CmdRunDirect)
	ev3.RunCommand(bot.OutB, ev3.CmdRunDirect)
	ev3.RunCommand(bot.OutC, ev3.CmdRunDirect)
	ev3.RunCommand(bot.OutD, ev3.CmdRunDirect)

	work := func(c int) {
		v := -50
		if c%2 == 1 {
			v = 50
		}
		ma.Value = v
		mb.Value = v
		mc.Value = v
		md.Value = -100
		in1.Sync()
		in2.Sync()
		in3.Sync()
		in4.Sync()
		ma.Sync()
		mb.Sync()
		mc.Sync()
		md.Sync()
		pd.Sync()
	}

	n, count := 0, 0
	ticks := time.Tick(10 * time.Millisecond)
	samples := time.Tick(1 * time.Second)

	for count < 10 {
		select {
		case <-ticks:
			work(count)
			n += 1
		case <-samples:
			log.Println("Sample", count, n)
			n = 0
			count += 1
		}
	}

	ev3.RunCommand(bot.OutA, ev3.CmdStop)
	ev3.RunCommand(bot.OutB, ev3.CmdStop)
	ev3.RunCommand(bot.OutC, ev3.CmdStop)
	ev3.RunCommand(bot.OutD, ev3.CmdStop)

	log.Println("Done.")
}
