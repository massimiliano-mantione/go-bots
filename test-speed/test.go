package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ev3go/ev3"
	"github.com/ev3go/ev3dev"
)

func motor(port string, driver string) *ev3dev.TachoMotor {
	m, err := ev3dev.TachoMotorFor(port, driver)
	if err != nil {
		log.Fatalf("failed to find medium motor on port %v: %v", port, err)
	}
	err = m.SetStopAction("coast").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for medium motor on outA: %v", err)
	}
	return m
}

func sensorIR(port string) *ev3dev.Sensor {
	driver := "lego-ev3-ir"
	s, err := ev3dev.SensorFor(port, driver)
	if err != nil {
		log.Fatalf("failed to find sensor with driver %v on port %v: %v", driver, port, err)
	}
	s.SetMode("IR-PROX")
	return s
}

func sensorColor(port string) *ev3dev.Sensor {
	driver := "lego-ev3-color"
	s, err := ev3dev.SensorFor(port, driver)
	if err != nil {
		log.Fatalf("failed to find sensor with driver %v on port %v: %v", driver, port, err)
	}
	s.SetMode("COL-REFLECT")
	return s
}

func main() {
	ev3.LCD.Init(true)
	defer ev3.LCD.Close()

	motorF := motor("outD", "lego-ev3-m-motor")
	motorR := motor("outA", "lego-ev3-l-motor")
	motorC := motor("outB", "lego-ev3-l-motor")
	motorL := motor("outC", "lego-ev3-l-motor")

	checkErrors(motorF, motorR, motorC, motorL)

	motorF.Command("reset")
	motorR.Command("reset")
	motorC.Command("reset")
	motorL.Command("reset")

	motorF.Command("run-direct")
	motorR.Command("run-direct")
	motorC.Command("run-direct")
	motorL.Command("run-direct")

	setMotors := func(v float32) {
		motorF.SetDutyCycleSetpoint(-100)
		motorR.SetDutyCycleSetpoint(int(-v * 100))
		motorC.SetDutyCycleSetpoint(int(-v * 100))
		motorL.SetDutyCycleSetpoint(int(-v * 100))
	}

	irR := sensorIR("in1")
	irL := sensorIR("in2")
	borderR := sensorColor("in3")
	borderL := sensorColor("in4")

	readSensors := func() {
		GetValue(irR, 0)
		GetValue(irL, 0)
		GetValue(borderR, 0)
		GetValue(borderL, 0)
	}

	setMotors(0)
	readSensors()

	tick := time.Tick(10 * time.Millisecond)
	sample := time.Tick(1 * time.Second)

	log.Println("Sampling...")

	ticks := 0
	samples := 0
	var v float32 = 0
	for samples < 10 {
		select {
		case <-tick:
			ticks += 1
			if samples%2 == 1 {
				v = -1
			} else {
				v = 1
			}
			setMotors(v)
			readSensors()
		case <-sample:
			samples += 1
			log.Println("samples", samples, "ticks", ticks, "v", v)
			ticks = 0
		}
	}

	log.Println("Stop.")

	motorF.Command("stop")
	motorR.Command("stop")
	motorC.Command("stop")
	motorL.Command("stop")
}

func GetValue(s *ev3dev.Sensor, v int) string {
	val, err := s.Value(v)
	if err != nil {
		log.Fatalf("Cannot read from sensor %v", s)
	}
	return val
}

func checkErrors(devs ...ev3dev.Device) {
	for _, d := range devs {
		err := d.(*ev3dev.TachoMotor).Err()
		if err != nil {
			drv, dErr := ev3dev.DriverFor(d)
			if dErr != nil {
				drv = fmt.Sprintf("(missing driver name: %v)", dErr)
			}
			addr, aErr := ev3dev.AddressOf(d)
			if aErr != nil {
				drv = fmt.Sprintf("(missing port address: %v)", aErr)
			}
			log.Fatalf("motor error for %s:%s on port %s: %v", d, drv, addr, err)
		}
	}
}
