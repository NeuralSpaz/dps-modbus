//go:generate stringer -type=Mode
//go:generate stringer -type=Protection
//go:generate stringer -type=Lock
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	fmt.Println("starting dps5020 monitor in modbus mode")

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	// handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	dps := new(DPS)
	dps.conn = client

	// dps.readPresets()
	dps.readStatus()

	dps.RLock()
	// fmt.Printf("%v", dps)
	// for k, preset := range dps.PreSets {
	fmt.Printf("PresetM0:\n%s\n", dps.PreSets[0])
	// 	if k > 2 {
	// 		break
	// 	}
	// }
	dps.RUnlock()
	targetVoltage := 12.0
	initalVolage := 1.0

	if err := dps.setVoltage(0.0); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := dps.enableOutput(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	end := time.Now().Add(time.Second * 60)
	for {
		if time.Now().After(end) {
			break
		}
		if targetVoltage > initalVolage {
			initalVolage += 1.0
			if err := dps.setVoltage(initalVolage); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}
		dps.readStatus()
		dps.RLock()
		fmt.Println(dps.Statuz)
		dps.RUnlock()
	}
	if err := dps.disableOutput(); err != nil {
		log.Println(err)
	}
}

type DPS struct {
	conn modbus.Client
	sync.RWMutex

	Statuz        Status
	PreSets       [10]Preset
	CurrentPreset int
	debug         bool
}

type Preset struct {
	VoltageSet            float64
	CurrentSet            float64
	OverVoltageProtection float64
	OverCurrentProtection float64
	OverPowerProtection   float64
	LedBrightness         uint16
	DataRecall            uint16
	PowerOutput           bool
}

func (p Preset) String() string {
	return fmt.Sprintf("\tSetVoltage: %2.2fV\n\tSetCurrent:%2.2fA\n\tOVP:%2.2fV\n\tOCP:%2.2fA\n\tOPP:%2.2fW\n\tLED Brightness:%v\n\tDataRecall:%v\n\tOutput Enabled On Start:%t\n",
		p.VoltageSet,
		p.CurrentSet,
		p.OverVoltageProtection,
		p.OverCurrentProtection,
		p.OverPowerProtection,
		p.LedBrightness,
		p.DataRecall,
		p.PowerOutput)
}

type Status struct {
	SetVoltage       float64
	SetCurrent       float64
	ActualVoltage    float64
	ActualCurrent    float64
	Power            float64
	SupplyVoltage    float64
	LockOut          Lock
	ProtectionTrip   Protection
	Constant         Mode
	OutputOn         Output
	DisplayBightness uint16
	Model            uint16
	Version          uint16
}

func (s Status) String() string {
	if s.ProtectionTrip != 0 {
		return fmt.Sprint(s.ProtectionTrip)
	}
	if s.SupplyVoltage < s.ActualVoltage+2 && s.Constant == ConstantCurrent {
		fmt.Println("Current limiting due to low supply voltage")
	}
	return fmt.Sprintf("V:%6.2fV\tAC:%6.2fA\tP: %6.2fW \t%v",
		s.ActualVoltage, s.ActualCurrent, s.Power, s.Constant)
}

type Lock uint16

const (
	Unlocked Lock = 0
	Locked   Lock = 1
)

type Output uint16

const (
	Off Output = 0
	On  Output = 1
)

type Protection uint16

const (
	OverVoltageProtection Protection = 1
	OverCurrentProtection Protection = 2
	OverPowerProtection   Protection = 3
)

type Mode uint16

const (
	ConstantVoltage Mode = 0
	ConstantCurrent Mode = 1
)

func (d *DPS) readStatus() error {
	d.Lock()
	defer d.Unlock()
	statusRaw, err := d.conn.ReadHoldingRegisters(0, 12)
	if err != nil {
		return err
	}
	d.Statuz = parseStatus(statusRaw)

	err = d.readPreset(0)
	return err

}

func parseStatus(raw []byte) Status {
	var s Status

	s.SetVoltage = floatFromBytes(raw[0:2])
	s.SetCurrent = floatFromBytes(raw[2:4])
	s.ActualVoltage = floatFromBytes(raw[4:6])
	s.ActualCurrent = floatFromBytes(raw[6:8])
	s.Power = floatFromBytes(raw[8:10])
	s.SupplyVoltage = floatFromBytes(raw[10:12])
	s.LockOut = Lock(binary.BigEndian.Uint16(raw[12:14]))
	s.ProtectionTrip = Protection(binary.BigEndian.Uint16(raw[14:16]))
	s.Constant = Mode(binary.BigEndian.Uint16(raw[16:18]))
	s.OutputOn = Output(binary.BigEndian.Uint16(raw[18:20]))
	s.DisplayBightness = binary.BigEndian.Uint16(raw[20:22])
	s.Model = binary.BigEndian.Uint16(raw[22:24])
	s.Version = binary.BigEndian.Uint16(raw[24:26])

	return s
}

func (d *DPS) readPresets() error {
	d.Lock()
	defer d.Unlock()
	for i := range d.PreSets {
		if err := d.readPreset(i); err != nil {
			return err
		}
	}
	return nil
}

func (d *DPS) readPreset(n int) error {
	presetRaw, err := d.conn.ReadHoldingRegisters(uint16(0x50+(n*0x10)), 8)
	if err != nil {
		return err
	}
	d.PreSets[n] = parsePresetBytes(presetRaw)
	if d.debug {
		log.Printf("M%d presets: %x\n", n, presetRaw)
		log.Println(d.PreSets[n])
	}

	return nil
}

func parsePresetBytes(raw []byte) Preset {
	var p Preset
	p.VoltageSet = floatFromBytes(raw[0:2])
	p.CurrentSet = floatFromBytes(raw[2:4])
	p.OverVoltageProtection = floatFromBytes(raw[4:6])
	p.OverCurrentProtection = floatFromBytes(raw[6:8])
	p.OverPowerProtection = floatFromBytes(raw[8:10]) * 10
	p.LedBrightness = binary.BigEndian.Uint16(raw[10:12])
	p.DataRecall = binary.BigEndian.Uint16(raw[12:14])
	if raw[15] > 0 {
		p.PowerOutput = true
	}
	return p
}

func floatFromBytes(b []byte) float64 {
	return float64(binary.BigEndian.Uint16(b)) / 100
}

func (d *DPS) enableOutput() error {
	resp, err := d.conn.WriteSingleRegister(onOffRegister, 1)
	status := Output(binary.BigEndian.Uint16(resp))
	d.Statuz.OutputOn = status
	if status != On {
		return fmt.Errorf("failed to turn on output")
	}
	return err
}

func (d *DPS) disableOutput() error {
	resp, err := d.conn.WriteSingleRegister(onOffRegister, 0)
	status := Output(binary.BigEndian.Uint16(resp))
	d.Statuz.OutputOn = status
	if status != Off {
		return fmt.Errorf("failed to turn off output")
	}
	return err
}

func (d *DPS) setVoltage(sv float64) error {
	resp, err := d.conn.WriteSingleRegister(voltageSetRegister, uint16(sv)*100)
	setVoltage := floatFromBytes(resp)
	d.Statuz.SetVoltage = setVoltage
	if err != nil {
		return fmt.Errorf("failed to set voltage")
	}
	return err
}
