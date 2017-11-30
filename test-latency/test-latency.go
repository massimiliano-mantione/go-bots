package main

import (
	"bufio"
	"fmt"
	"go-bots/ev3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		quit("Terminated by signal", sig)
	}()
}

var devs *ev3.Devices
var ir, c *ev3.Attribute

var ledRR, ledRG, ledLR, ledLG *ev3.Attribute

func leds(rl int, rr int, gr int, gl int) {
	ledLG.Value = gl
	ledLR.Value = rl
	ledRG.Value = gr
	ledRR.Value = rr
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()
}

func initialize() {
	initializeTime()

	devs = ev3.Scan(nil)

	// Sensors
	ev3.CheckDriver(devs.In1, ev3.DriverIr, ev3.In1)
	ev3.SetMode(devs.In1, ev3.IrModeProx)
	ir = ev3.OpenByteR(devs.In1, ev3.BinData)
	ev3.CheckDriver(devs.In2, ev3.DriverColor, ev3.In2)
	ev3.SetMode(devs.In2, ev3.ColorModeAmbient)
	c = ev3.OpenByteR(devs.In2, ev3.BinData)

	// Motor
	ev3.CheckDriver(devs.OutA, ev3.DriverTachoMotorMedium, ev3.OutA)
	ev3.WriteStringAttribute(devs.OutA, ev3.DutyCycleSp, "30")
	ev3.RunCommand(devs.OutA, ev3.CmdRunDirect)

	// Leds
	ledLG = ev3.OpenTextW(devs.LedLeftGreen, ev3.Brightness)
	ledLR = ev3.OpenTextW(devs.LedLeftRed, ev3.Brightness)
	ledRG = ev3.OpenTextW(devs.LedRightGreen, ev3.Brightness)
	ledRR = ev3.OpenTextW(devs.LedRightRed, ev3.Brightness)
	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()
}

func close() {
	// Sensors
	ir.Close()
	c.Close()

	// Motor
	ev3.RunCommand(devs.OutA, ev3.CmdStop)
	ev3.RunCommand(devs.OutA, ev3.CmdReset)

	// Leds
	ledLG.Value = 0
	ledLR.Value = 0
	ledRG.Value = 0
	ledRR.Value = 0
	ledLG.Sync()
	ledLR.Sync()
	ledRG.Sync()
	ledRR.Sync()
	ledLG.Close()
	ledLR.Close()
	ledRG.Close()
	ledRR.Close()
}

var initializationTime time.Time

func initializeTime() {
	initializationTime = time.Now()
}

func durationToTicks(d time.Duration) int {
	return int(d / 1000)
}
func timespanAsTicks(start time.Time, end time.Time) int {
	return durationToTicks(end.Sub(start))
}
func currentTicks() int {
	return timespanAsTicks(initializationTime, time.Now())
}
func ticksToMillis(ticks int) int {
	return ticks / 1000
}
func ticksDMillis(ticks int) int {
	return (ticks % 1000) / 100
}

func printf(s string, data ...interface{}) {
	fmt.Fprintf(os.Stderr, s, data...)
}

func quit(data ...interface{}) {
	PrintData()
	close()
	log.Fatalln(data...)
}

func trimIr() int {
	if ir.Value < 40 {
		return 1
	}
	return 0
}

// DataPoint holds one point of stats data
type DataPoint struct {
	Time    uint32
	Ticks   uint32
	Trigger uint32
	IR      uint8
	C       uint8
	Loops   uint8
}

// DataStats holds points stats data
type DataStats struct {
	Points           int
	PointsTicks      int
	Points00         int
	Points00Ticks    int
	Points10         int
	Points10Ticks    int
	Points20         int
	Points20Ticks    int
	Points30         int
	Points30Ticks    int
	Points40         int
	Points40Ticks    int
	Points50         int
	Points50Ticks    int
	Triggers         int
	TriggersTicks    int
	TriggersCounts   int
	Triggers00       int
	Triggers00Ticks  int
	Triggers00Counts int
	Triggers30       int
	Triggers30Ticks  int
	Triggers30Counts int
	Triggers60       int
	Triggers60Ticks  int
	Triggers60Counts int
	Triggers90       int
	Triggers90Ticks  int
	Triggers90Counts int
}

const maxPoints = 20000

var startIndex int
var nextIndex int
var points = [maxPoints]DataPoint{}

// DataReset resets the data
func DataReset() {
	fmt.Println("Data reset start")
	for i := 0; i < maxPoints; i++ {
		points[i] = DataPoint{}
	}
	startIndex = 0
	nextIndex = 1
	fmt.Println("Data reset done")
}

// DataStore stores a data point
func DataStore(time uint32,
	ticks uint32,
	ir uint8,
	c uint8,
	trigger uint32,
	loops uint8) {

	points[nextIndex] = DataPoint{
		Time:    time,
		Ticks:   ticks,
		IR:      ir,
		C:       c,
		Trigger: trigger,
		Loops:   loops,
	}
	nextIndex++
	if nextIndex >= maxPoints {
		nextIndex = 0
	}
	if startIndex == nextIndex {
		startIndex++
	}
}

func printDataPoint(w *bufio.Writer, p *DataPoint) {
	w.WriteString(fmt.Sprintf("NOW %6d.%d  TICKS %3d.%d  IR %1d  C %3d  D %4d.%d L %2d\n",
		ticksToMillis(int(p.Time)),
		ticksDMillis(int(p.Time)),
		ticksToMillis(int(p.Ticks)),
		ticksDMillis(int(p.Ticks)),
		p.IR,
		p.C,
		ticksToMillis(int(p.Trigger)),
		ticksDMillis(int(p.Trigger)),
		p.Loops))
}

func processDataPoint(p *DataPoint, d *DataStats) {
	ticks := int(p.Ticks)
	ticksM := ticksToMillis(ticks)
	d.Points++
	d.PointsTicks += ticks
	if ticksM < 10 {
		d.Points00++
		d.Points10Ticks += ticks
	} else if ticksM < 10 {
		d.Points00++
		d.Points00Ticks += ticks
	} else if ticksM < 20 {
		d.Points10++
		d.Points10Ticks += ticks
	} else if ticksM < 30 {
		d.Points20++
		d.Points20Ticks += ticks
	} else if ticksM < 40 {
		d.Points30++
		d.Points30Ticks += ticks
	} else if ticksM < 50 {
		d.Points40++
		d.Points40Ticks += ticks
	} else {
		d.Points50++
		d.Points50Ticks += ticks
	}

	if p.Trigger != 0 {
		trigger := int(p.Trigger)
		triggerM := ticksToMillis(trigger)
		triggerCount := int(p.Loops)
		d.Triggers++
		d.TriggersTicks += trigger
		d.TriggersCounts += triggerCount
		if triggerM < 30 {
			d.Triggers00++
			d.Triggers00Ticks += trigger
			d.Triggers00Counts += triggerCount
		} else if triggerM < 60 {
			d.Triggers30++
			d.Triggers30Ticks += trigger
			d.Triggers30Counts += triggerCount
		} else if triggerM < 90 {
			d.Triggers60++
			d.Triggers60Ticks += trigger
			d.Triggers60Counts += triggerCount
		} else {
			d.Triggers90++
			d.Triggers90Ticks += trigger
			d.Triggers90Counts += triggerCount
		}
	}
}

func fmtPointsLine(kind string, count int, ticks int) string {
	avg := 0
	if count > 0 {
		avg = ticks / count
	}
	return fmt.Sprintf("%s %6d %3d.%1d\n", kind, count, ticksToMillis(avg), ticksDMillis(avg))
}
func fmtTriggersLine(kind string, count int, ticks int, loops int) string {
	avg := 0
	loopAvg := 0
	if count > 0 {
		avg = ticks / count
		loopAvg = loops / count
	}
	return fmt.Sprintf("%s %6d %3d.%1d L %2d\n", kind, count, ticksToMillis(avg), ticksDMillis(avg), loopAvg)
}

// PrintData prints the data
func PrintData() {
	file, err := os.OpenFile("data.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	defer w.Flush()

	var d DataStats
	w.WriteString("Data points START\n")
	for i := startIndex; i < maxPoints; i++ {
		if i == nextIndex {
			break
		}
		p := &points[i]
		if p.Ticks != 0 {
			printDataPoint(w, p)
			processDataPoint(p, &d)
		}
	}
	for i := 0; i < nextIndex; i++ {
		p := &points[i]
		if p.Ticks != 0 {
			printDataPoint(w, p)
			processDataPoint(p, &d)
		}
	}
	w.WriteString("Data points END\n")
	w.WriteString("Data stats START\n")
	w.WriteString(fmtPointsLine("P  ", d.Points, d.PointsTicks))
	w.WriteString(fmtPointsLine("P00", d.Points00, d.Points00Ticks))
	w.WriteString(fmtPointsLine("P10", d.Points10, d.Points10Ticks))
	w.WriteString(fmtPointsLine("P20", d.Points20, d.Points20Ticks))
	w.WriteString(fmtPointsLine("P30", d.Points30, d.Points30Ticks))
	w.WriteString(fmtPointsLine("P40", d.Points40, d.Points40Ticks))
	w.WriteString(fmtPointsLine("P50", d.Points50, d.Points50Ticks))
	w.WriteString(fmtTriggersLine("T  ", d.Triggers, d.TriggersTicks, d.TriggersCounts))
	w.WriteString(fmtTriggersLine("T00", d.Triggers00, d.Triggers00Ticks, d.Triggers00Counts))
	w.WriteString(fmtTriggersLine("T30", d.Triggers30, d.Triggers30Ticks, d.Triggers30Counts))
	w.WriteString(fmtTriggersLine("T60", d.Triggers60, d.Triggers60Ticks, d.Triggers60Counts))
	w.WriteString(fmtTriggersLine("T90", d.Triggers90, d.Triggers90Ticks, d.Triggers90Counts))
	w.WriteString("Data stats END\n")
	DataReset()
}

func main() {
	handleSignals()
	initialize()
	defer close()

	DataReset()
	defer PrintData()

	// profile.ProfilePath("/home/robot")
	// defer profile.Start().Stop()

	previousTicks := currentTicks()
	trigger := 0
	loops := 0
	for {
		now := currentTicks()
		ticks := now - previousTicks
		previousTicks = now

		ir.Sync()
		irT := trimIr()
		c.Sync()
		cV := c.Value
		cT := 0
		if cV > 20 {
			cT = 1
		}

		if irT == 0 {
			leds(0, 0, 0, 0)
		} else {
			leds(255, 255, 255, 255)
		}

		triggerV := 0
		loopsV := 0
		if irT != cT {
			trigger += ticks
			triggerV = 0
			loops++
			loopsV = 0
		} else {
			if trigger > 0 {
				triggerV = trigger + ticks
				loopsV = loops + 1
			} else {
				triggerV = 0
				loopsV = 0
			}
			trigger = 0
			loops = 0
		}

		DataStore(uint32(now), uint32(ticks), uint8(irT), uint8(cT), uint32(triggerV), uint8(loopsV))

		if ticksToMillis(now) > 10000 {
			break
		}

		time.Sleep(time.Microsecond)
	}

	leds(0, 0, 0, 0)
	ev3.RunCommand(devs.OutA, ev3.CmdStop)
}
