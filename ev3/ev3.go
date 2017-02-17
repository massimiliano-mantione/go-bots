package ev3

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	fp "path/filepath"
	"strings"
)

type Ev3Devices struct {
	In1  string
	In2  string
	In3  string
	In4  string
	OutA string
	OutB string
	OutC string
	OutD string
	LedL string
	LedR string
}

const DriverIr = "lego-ev3-ir"
const DriverColor = "lego-ev3-color"
const DriverTachoMotorLarge = "lego-ev3-l-motor"
const DriverTachoMotorMedium = "lego-ev3-m-motor"
const DriverDcMotor = "dc-motor"

const IrModeProx = "IR-PROX"
const IrModeRemote = "IR-REMOTE"

const ColorModeReflect = "COL-REFLECT"
const ColorModeColor = "COL-COLOR"
const ColorModeAmbient = "COL-AMBIENT"
const ColorModeRaw = "COL-RAW"

const In1 = "in1"
const In2 = "in2"
const In3 = "in3"
const In4 = "in4"
const OutA = "outA"
const OutB = "outB"
const OutC = "outC"
const OutD = "outD"

const Address = "address"
const Value0 = "value0"
const Value1 = "value1"
const Value2 = "value2"
const Value3 = "value3"
const BinData = "bin_data"
const Mode = "mode"
const Modes = "modes"
const Command = "command"
const Commands = "commands"
const DriverName = "driver_name"

const CountPerRot = "count_per_rot"
const DutyCycle = "duty_cycle"
const DutyCycleSp = "duty_cycle_sp"
const MaxSpeed = "max_speed"
const Polarity = "polarity"
const Position = "position"
const PositionSp = "position_sp"
const Speed = "speed"
const SpeedSp = "speed_sp"
const State = "state"
const StopAction = "stop_action"
const StopActions = "stop_actions"
const TimeSp = "time_sp"
const Uevent = "uevent"

const CmdRunForever = "run-forever"
const CmdRunToAbsPos = "run-to-abs-pos"
const CmdRunToRelPos = "run-to-rel-pos"
const CmdRunTimed = "run-timed"
const CmdRunDirect = "run-direct"
const CmdStop = "stop"
const CmdReset = "reset"

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

func ReadStringAttribute(dev string, attr string) string {
	return readString(fp.Join(dev, attr))
}
func WriteStringAttribute(dev string, attr string, v string) {
	writeString(fp.Join(dev, attr), v)
}

func SetMode(dev string, mode string) {
	modesText := ReadStringAttribute(dev, Modes)
	modes := strings.Split(modesText, " ")
	if !contains(modes, mode) {
		log.Fatalln("Device", dev, "does not support mode", mode)
	}
	WriteStringAttribute(dev, Mode, mode)
}

func SetStopAction(dev string, action string) {
	actionsText := ReadStringAttribute(dev, StopActions)
	actions := strings.Split(actionsText, " ")
	if !contains(actions, action) {
		log.Fatalln("Device", dev, "does not support stop action", action)
	}
	WriteStringAttribute(dev, StopAction, action)
}

func CheckCommand(dev string, cmd string) {
	commandsText := ReadStringAttribute(dev, Commands)
	commands := strings.Split(commandsText, " ")
	if !contains(commands, cmd) {
		log.Fatalln("Device", dev, "does not support command", cmd)
	}
}

func RunCommand(dev string, cmd string) {
	WriteStringAttribute(dev, Command, cmd)
}

func CheckDriver(dev string, driver string) {
	actualDriver := ReadStringAttribute(dev, DriverName)
	if actualDriver != driver {
		log.Fatalln("Device", dev, "has driver", actualDriver, "instead of", driver)
	}
}

const attributeBufSize int = 64

type Attribute struct {
	path         string
	file         *os.File
	currentValue int
	Value        int
	writable     bool
	text         bool
	buf          [attributeBufSize]byte
}

func (a *Attribute) Path() string {
	return a.path
}
func (a *Attribute) Writable() bool {
	return a.writable
}
func (a *Attribute) Close() {
	err := a.file.Close()
	if err != nil {
		log.Fatalln("Cannot close dev", a.path, ":", err)
	}
	a.file = nil
}

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
					digits += 1
					ceiling *= 10
				}
				length := digits
				if v < 0 {
					length += 1
				}
				toWrite = a.buf[0:length]
				index := 0
				if v < 0 {
					toWrite[0] = '-'
					index += 1
				}
				for digits > 0 {
					ceiling /= 10
					digitValue := abs / ceiling
					toWrite[index] = '0' + byte(digitValue)
					index += 1
					abs = abs % ceiling
					digits -= 1
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
				index -= 1
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

	log.Println("Open", path, "writable", writable, "text", text)

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

func OpenByteR(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, false, false)
}
func OpenByteW(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, true, false)
}
func OpenTextR(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, false, true)
}
func OpenTextW(dev string, attr string) *Attribute {
	return OpenAttribute(dev, attr, true, true)
}

func Scan() *Ev3Devices {
	devs := Ev3Devices{}
	classes := "/sys/class"

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
