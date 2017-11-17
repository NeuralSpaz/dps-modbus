//go:generate stringer -type=Mode
//go:generate stringer -type=Protection
//go:generate stringer -type=Lock
package main

import (
	"fmt"
	"sync"

	"github.com/goburrow/modbus"
)

type DPS struct {
	conn modbus.Client
	sync.RWMutex

	puber chan<- Status

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
	if s.SupplyVoltage < s.ActualVoltage+2 && s.Constant == CC {
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
	None                  Protection = 0
	OverVoltageProtection Protection = 1
	OverCurrentProtection Protection = 2
	OverPowerProtection   Protection = 3
)

type Mode uint16

const (
	CV Mode = 0
	CC Mode = 1
)
