package ev3

import (
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	fp "path/filepath"
	"strings"
	"time"
)

// Devices contains connected devices
type Devices struct {
	Port0         string
	Port1         string
	Port2         string
	Port3         string
	Port4         string
	Port5         string
	Port6         string
	Port7         string
	In1           string
	In2           string
	In3           string
	In4           string
	OutA          string
	OutB          string
	OutC          string
	OutD          string
	LedRightGreen string
	LedRightRed   string
	LedLeftGreen  string
	LedLeftRed    string
}

// DriverIr IR sensor driver constant
const DriverIr = "lego-ev3-ir"

// DriverColor color sensor driver constant
const DriverColor = "lego-ev3-color"

// DriverTachoMotorLarge large tacho motor driver constant
const DriverTachoMotorLarge = "lego-ev3-l-motor"

// DriverTachoMotorMedium medium tacho motor driver constant
const DriverTachoMotorMedium = "lego-ev3-m-motor"

// DriverDcMotor DC motor driver constant
const DriverDcMotor = "dc-motor"

// DriverRcxMotor rcx motor driver constant
const DriverRcxMotor = "rcx-motor"

// IrModeProx IR sensor proximity mode
const IrModeProx = "IR-PROX"

// IrModeRemote IR sensor remote control mode
const IrModeRemote = "IR-REMOTE"

// ColorModeReflect color sensor reflective mode
const ColorModeReflect = "COL-REFLECT"

// ColorModeColor color sensor RGB mode
const ColorModeColor = "COL-COLOR"

// ColorModeAmbient color sensor ambient light mode
const ColorModeAmbient = "COL-AMBIENT"

// ColorModeRaw color sensor raw mode
const ColorModeRaw = "COL-RAW"

// ColorModeRefRaw color sensor reflective raw mode
const ColorModeRefRaw = "REF-RAW"

// ColorModeRgbRaw color sensor RGB raw mode
const ColorModeRgbRaw = "RGB-RAW"

// In1 input port 1
const In1 = "in1"

// In2 input port 1
const In2 = "in2"

// In3 input port 1
const In3 = "in3"

// In4 input port 4
const In4 = "in4"

// OutA output port A
const OutA = "outA"

// OutB output port B
const OutB = "outB"

// OutC output port C
const OutC = "outC"

// OutD output port D
const OutD = "outD"

// Address attribute
const Address = "address"

// Value0 attribute
const Value0 = "value0"

// Value1 attribute
const Value1 = "value1"

// Value2 attribute
const Value2 = "value2"

// Value3 attribute
const Value3 = "value3"

// BinData attribute
const BinData = "bin_data"

// Mode attribute
const Mode = "mode"

// Modes attribute
const Modes = "modes"

// Command attribute
const Command = "command"

// Commands attribute
const Commands = "commands"

// DriverName attribute
const DriverName = "driver_name"

// CountPerRot attribute
const CountPerRot = "count_per_rot"

// DutyCycle attribute
const DutyCycle = "duty_cycle"

// DutyCycleSp attribute
const DutyCycleSp = "duty_cycle_sp"

// MaxSpeed attribute
const MaxSpeed = "max_speed"

// Polarity attribute
const Polarity = "polarity"

// Position attribute
const Position = "position"

// PositionSp attribute
const PositionSp = "position_sp"

// Speed attribute
const Speed = "speed"

// SpeedSp attribute
const SpeedSp = "speed_sp"

// State attribute
const State = "state"

// StopAction attribute
const StopAction = "stop_action"

// StopActions attribute
const StopActions = "stop_actions"

// TimeSp attribute
const TimeSp = "time_sp"

// Uevent attribute
const Uevent = "uevent"

// CmdRunForever motor command
const CmdRunForever = "run-forever"

// CmdRunToAbsPos motor command
const CmdRunToAbsPos = "run-to-abs-pos"

// CmdRunToRelPos motor command
const CmdRunToRelPos = "run-to-rel-pos"

// CmdRunTimed motor command
const CmdRunTimed = "run-timed"

// CmdRunDirect motor command
const CmdRunDirect = "run-direct"

// CmdStop motor command
const CmdStop = "stop"

// CmdReset motor command
const CmdReset = "reset"

// Brightness led attribute
const Brightness = "brightness"

// MaxBrightness led attribute
const MaxBrightness = "max_brightness"

// OutPortModes is used to set out port modes
type OutPortModes struct {
	OutA string
	OutB string
	OutC string
	OutD string
}

// OutPortModeAuto auto mode
const OutPortModeAuto = "auto"

// OutPortModeTachoMotor tacho-motor mode
const OutPortModeTachoMotor = "tacho-motor"

// OutPortModeDcMotor dc-motor mode
const OutPortModeDcMotor = "dc-motor"

// OutPortModeLed led mode
const OutPortModeLed = "led"

// OutPortModeRaw raw mode
const OutPortModeRaw = "raw"

func readString(fileName string) string {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Cannot read from file", fileName, ":", err)
	}
	text := string(buf)
	text = strings.TrimSuffix(text, "\n")
	return text
}

func writeString(fileName string, v string) {
	err := ioutil.WriteFile(fileName, []byte(v), 0644)
	if err != nil {
		log.Fatalln("Cannot write to file", fileName, ":", err)
	}
}

func contains(stringSlice []string, search string) bool {
	for _, value := range stringSlice {
		if value == search {
			return true
		}
	}
	return false
}

// ReadStringAttribute reads the value of a string attribute on a device
func ReadStringAttribute(dev string, attr string) string {
	return readString(fp.Join(dev, attr))
}

// WriteStringAttribute writes the value of a string attribute on a device
func WriteStringAttribute(dev string, attr string, v string) {
	writeString(fp.Join(dev, attr), v)
}

// SetMode sets the mode attribute on a device
func SetMode(dev string, mode string) {
	modesText := ReadStringAttribute(dev, Modes)
	modes := strings.Split(modesText, " ")
	if !contains(modes, mode) {
		log.Fatalln("Device", dev, "does not support mode", mode)
	}
	WriteStringAttribute(dev, Mode, mode)
}

// SetStopAction sets the stop action on a motor device
func SetStopAction(dev string, action string) {
	actionsText := ReadStringAttribute(dev, StopActions)
	actions := strings.Split(actionsText, " ")
	if !contains(actions, action) {
		log.Fatalln("Device", dev, "does not support stop action", action)
	}
	WriteStringAttribute(dev, StopAction, action)
}

// CheckCommand checks that a device supports a given command
func CheckCommand(dev string, cmd string) {
	commandsText := ReadStringAttribute(dev, Commands)
	commands := strings.Split(commandsText, " ")
	if !contains(commands, cmd) {
		log.Fatalln("Device", dev, "does not support command", cmd)
	}
}

// RunCommand writes a value to the command attribute
func RunCommand(dev string, cmd string) {
	WriteStringAttribute(dev, Command, cmd)
}

// CheckDriver checks that a device has the given driver
func CheckDriver(dev string, driver string, port string) {
	if dev == "" {
		log.Fatalln("Port", port, "has no device instead of expected driver", driver)
	}
	actualDriver := ReadStringAttribute(dev, DriverName)
	if actualDriver != driver {
		log.Fatalln("Device", dev, "in port", port, "has driver", actualDriver, "instead of", driver)
	}
}

// CheckMode checks that a device is set to the given mode
func CheckMode(dev string, mode string) bool {
	actualMode := ReadStringAttribute(dev, Mode)
	return actualMode == mode
}

const attributeBufSize int = 64

// Attribute holds an open file on a given attribute so that subsequent operations are faster
type Attribute struct {
	path         string
	file         *os.File
	currentValue int
	Value        int
	writable     bool
	text         bool
	buf          [attributeBufSize]byte
}

// Path gets the attribute file full path
func (a *Attribute) Path() string {
	return a.path
}

// Writable checks if the attribute is writable
func (a *Attribute) Writable() bool {
	return a.writable
}

// Close closes the attribute
func (a *Attribute) Close() {
	err := a.file.Close()
	if err != nil {
		log.Fatalln("Cannot close dev", a.path, ":", err)
	}
	a.file = nil
}

// Sync reads or writes the attribute value (according to the "writable" status)
func (a *Attribute) Sync() {
	if a.file == nil {
		log.Fatalln("Cannot write to closed attribute", a.path)
	}
	if a.writable {
		if a.Value != a.currentValue {
			var toWrite []byte
			if a.text {
				v := a.Value
				abs := v
				if abs < 0 {
					abs = -abs
				}
				digits := 1
				ceiling := 10
				for abs >= ceiling {
					digits++
					ceiling *= 10
				}
				length := digits
				if v < 0 {
					length++
				}
				toWrite = a.buf[0:length]
				index := 0
				if v < 0 {
					toWrite[0] = '-'
					index++
				}
				for digits > 0 {
					ceiling /= 10
					digitValue := abs / ceiling
					toWrite[index] = '0' + byte(digitValue)
					index++
					abs = abs % ceiling
					digits--
				}
			} else {
				a.buf[0] = byte(a.Value)
				toWrite = a.buf[0:1]
			}
			n, err := a.file.WriteAt(toWrite, 0)
			if err != nil {
				log.Fatalln("Cannot write to attribute file", a.path, ":", err)
			}
			if n != len(toWrite) {
				log.Fatalln("Cannot write bytes to attribute file", a.path, ":", err)
			}
			a.currentValue = a.Value
		}
	} else {
		if a.text {
			_, err := a.file.Seek(0, 0)
			if err != nil {
				log.Fatalln("Cannot rewind text attribute file", a.path, ":", err)
			}
			toRead, err := ioutil.ReadAll(a.file)
			if err != nil {
				log.Fatalln("Cannot read from text attribute file", a.path, ":", err)
			}
			if len(toRead) == 0 {
				log.Fatalln("Cannot read one byte from text attribute file", a.path, ":", err)
			}
			if toRead[len(toRead)-1] == '\n' {
				toRead = toRead[0 : len(toRead)-1]
			}
			isNegative := false
			if toRead[0] == '-' {
				toRead = toRead[1:]
				isNegative = true
			}
			v := 0
			digitValue := 1
			index := len(toRead) - 1
			for index >= 0 {
				digit := int(toRead[index] - '0')
				v += digitValue * digit
				index--
				digitValue *= 10
			}
			if isNegative {
				v = -v
			}
			a.Value = v
		} else {
			n, err := a.file.ReadAt(a.buf[0:1], 0)
			if err != nil {
				log.Fatalln("Cannot read from attribute file", a.path, ":", err)
			}
			if n != 1 {
				log.Fatalln("Cannot read one byte from attribute file", a.path, ":", err)
			}
			a.Value = int(a.buf[0])
		}
	}
}

// OpenAttribute creates an attribute struct opening the relevant file
func OpenAttribute(dev string, attr string, writable bool, text bool) *Attribute {
	path := fp.Join(dev, attr)
	flag := os.O_RDONLY
	if writable {
		flag = os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		log.Fatalln("Cannot open device", dev, ":", err)
	}
	return &Attribute{
		path:         path,
		file:         f,
		currentValue: math.MaxInt32,
		Value:        0,
		writable:     writable,
		text:         text,
		buf:          [attributeBufSize]byte{},
	}
}

// OpenByteR opens a byte attribute for reading
func OpenByteR(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, false, false)
}

// OpenByteW opens a byte attribute for writing
func OpenByteW(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, true, false)
}

// OpenTextR opens a string attribute for reading
func OpenTextR(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, false, true)
}

// OpenTextW opens a string attribute for writing
func OpenTextW(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, true, true)
}

// Scan scans the EV3 for devices and returns the structure describing them
func Scan(outModes *OutPortModes) *Devices {
	devs := Devices{}
	classes := "/sys/class"

	ports := fp.Join(classes, "lego-port")
	devs.Port0 = fp.Join(ports, "port0")
	devs.Port1 = fp.Join(ports, "port1")
	devs.Port2 = fp.Join(ports, "port2")
	devs.Port3 = fp.Join(ports, "port3")
	devs.Port4 = fp.Join(ports, "port4")
	devs.Port5 = fp.Join(ports, "port5")
	devs.Port6 = fp.Join(ports, "port6")
	devs.Port7 = fp.Join(ports, "port7")

	if outModes == nil {
		outModes = &OutPortModes{}
	}
	if outModes.OutA == "" {
		outModes.OutA = OutPortModeAuto
	}
	if outModes.OutB == "" {
		outModes.OutB = OutPortModeAuto
	}
	if outModes.OutC == "" {
		outModes.OutC = OutPortModeAuto
	}
	if outModes.OutD == "" {
		outModes.OutD = OutPortModeAuto
	}

	sleep := false
	if !CheckMode(devs.Port4, outModes.OutA) {
		SetMode(devs.Port4, outModes.OutA)
		sleep = true
	}
	if !CheckMode(devs.Port5, outModes.OutB) {
		SetMode(devs.Port5, outModes.OutB)
		sleep = true
	}
	if !CheckMode(devs.Port6, outModes.OutC) {
		SetMode(devs.Port6, outModes.OutC)
		sleep = true
	}
	if !CheckMode(devs.Port7, outModes.OutD) {
		SetMode(devs.Port7, outModes.OutD)
		sleep = true
	}
	if sleep {
		log.Println("Sleep...")
		time.Sleep(500 * time.Millisecond)
		log.Println("Slept.")
	}

	leds := fp.Join(classes, "leds")
	devs.LedRightGreen = fp.Join(leds, "ev3:right:green:ev3dev")
	devs.LedRightRed = fp.Join(leds, "ev3:right:red:ev3dev")
	devs.LedLeftGreen = fp.Join(leds, "ev3:left:green:ev3dev")
	devs.LedLeftRed = fp.Join(leds, "ev3:left:red:ev3dev")

	sensors, _ := fp.Glob(fp.Join(classes, "/*/sensor*"))
	for _, s := range sensors {
		port := ReadStringAttribute(s, Address)
		switch port {
		case In1:
			devs.In1 = s
		case In2:
			devs.In2 = s
		case In3:
			devs.In3 = s
		case In4:
			devs.In4 = s
		default:
			log.Fatalln("Unknown port", port, "for sensor", s)
		}
	}

	tachoMotors, _ := fp.Glob(fp.Join(classes, "/tacho-motor/*"))
	for _, m := range tachoMotors {
		port := ReadStringAttribute(m, Address)
		switch port {
		case OutA:
			devs.OutA = m
		case OutB:
			devs.OutB = m
		case OutC:
			devs.OutC = m
		case OutD:
			devs.OutD = m
		default:
			log.Fatalln("Unknown port", port, "for tacho motor", m)
		}
	}

	dcMotors, _ := fp.Glob(fp.Join(classes, "/dc-motor/*"))
	for _, m := range dcMotors {
		port := ReadStringAttribute(m, Address)
		switch port {
		case OutA:
			devs.OutA = m
		case OutB:
			devs.OutB = m
		case OutC:
			devs.OutC = m
		case OutD:
			devs.OutD = m
		default:
			log.Fatalln("Unknown port", port, "for dc motor", m)
		}
	}
	return &devs
}

// DurationToMillis converts a time.Duration to milliseconds
func DurationToMillis(d time.Duration) int {
	return int(d / 1000000)
}

// TimespanAsMillis converts a time interval to milliseconds
func TimespanAsMillis(start time.Time, end time.Time) int {
	return DurationToMillis(end.Sub(start))
}

// Direction represents left or right
type Direction int

const (
	// NoDirection means no direction
	NoDirection Direction = 0
	// Left direction
	Left = -1
	// Right direction
	Right = 1
)

// ChangeDirection returns the opposite direction
func ChangeDirection(dir Direction) Direction {
	return -dir
}

// LeftTurnVersor is the left turn speeed coefficient
func LeftTurnVersor(dir Direction) int {
	return int(dir)
}

// RightTurnVersor is the right turn speeed coefficient
func RightTurnVersor(dir Direction) int {
	return -int(dir)
}

const keyUp = 103
const keyDown = 108
const keyLeft = 105
const keyRight = 106
const keyEnter = 28
const keyBackspace = 14

// Buttons contains the current state of buttons
type Buttons struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
	Enter bool
	Back  bool
	stop  chan bool
}

const keyUpTime = time.Second / 10

// OpenButtons starts listening for button changes
func OpenButtons(readStdin bool) *Buttons {
	result := Buttons{}
	result.stop = make(chan bool)

	go func() {
		buttonDev := "/dev/input/by-path/platform-gpio-keys.0-event"
		f, err := os.OpenFile(buttonDev, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatalln("Cannot open file", buttonDev, ":", err)
		}

		data := make(chan [16]byte)
		go func() {
			var event [16]byte
			for {
				if f == nil {
					return
				}
				n, err := f.Read(event[0:16])
				if err != nil {
					if err == io.EOF || err == io.ErrClosedPipe {
						return
					}
					log.Fatalln("Error reading button events:", err)
				}
				if n != 16 {
					log.Fatalln("Event length is not 16 bytes:", n)
				}
				data <- event
			}
		}()

		if readStdin {
			go func() {
				consoleReader := bufio.NewReaderSize(os.Stdin, 1)
				for {
					key, _ := consoleReader.ReadByte()
					if f == nil {
						return
					}
					if key == 9 {
						// Tab -> Back
						result.Back = true
						time.AfterFunc(keyUpTime, func() {
							result.Back = false
						})
					} else if key == 32 {
						// Space -> Enter
						result.Enter = true
						time.AfterFunc(keyUpTime, func() {
							result.Enter = false
						})
					} else if key == 'a' || key == 'A' {
						// A -> Left
						result.Left = true
						time.AfterFunc(keyUpTime, func() {
							result.Left = false
						})
					} else if key == 'd' || key == 'D' {
						// D -> Right
						result.Right = true
						time.AfterFunc(keyUpTime, func() {
							result.Right = false
						})
					} else if key == 'w' || key == 'W' {
						// W -> Up
						result.Up = true
						time.AfterFunc(keyUpTime, func() {
							result.Up = false
						})
					} else if key == 's' || key == 'S' {
						// S -> Down
						result.Down = true
						time.AfterFunc(keyUpTime, func() {
							result.Down = false
						})
					}
				}
			}()
		}

		for {
			select {
			case event := <-data:
				eventType := binary.LittleEndian.Uint16(event[8:])
				eventCode := binary.LittleEndian.Uint16(event[10:])
				eventValue := binary.LittleEndian.Uint32(event[12:])
				if eventType == 1 {
					if eventCode == keyUp {
						if eventValue == 1 {
							result.Up = true
						} else {
							result.Up = false
						}
					} else if eventCode == keyUp {
						if eventValue == 1 {
							result.Up = true
						} else {
							result.Up = false
						}
					} else if eventCode == keyDown {
						if eventValue == 1 {
							result.Down = true
						} else {
							result.Down = false
						}
					} else if eventCode == keyLeft {
						if eventValue == 1 {
							result.Left = true
						} else {
							result.Left = false
						}
					} else if eventCode == keyRight {
						if eventValue == 1 {
							result.Right = true
						} else {
							result.Right = false
						}
					} else if eventCode == keyEnter {
						if eventValue == 1 {
							result.Enter = true
						} else {
							result.Enter = false
						}
					} else if eventCode == keyBackspace {
						if eventValue == 1 {
							result.Back = true
						} else {
							result.Back = false
						}
					}
				}
			case <-result.stop:
				f.Close()
				f = nil
				return
			}
		}

	}()

	return &result
}

// Close stops listening for button changes
func (b *Buttons) Close() {
	b.stop <- true
}
