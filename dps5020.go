package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	fmt.Println("starting dps5020 monitor in modbus mode")
	// client := modbus.RTUClient("/dev/ttyUSB0")
	// // results, err := client.ReadCoils(2, 1)
	// results, err := client.ReadHoldingRegisters(0, 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(results)
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
	defer handler.Close()
	results, err := client.ReadHoldingRegisters(0, 10)
	// results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		log.Fatal(err, results)
	}
	log.Println(results)
	// var amps uint16
	for i := 0; i < len(results); i += 2 {
		fmt.Printf("%2d: %02x %02x ", i, results[i], results[i+1])
		// if k == 3 {
		// 	// fmt.Println("bytes ", results[2:5])
		// 	amps = binary.BigEndian.Uint16(results[2:4])
		// 	fmt.Printf("Amps %2.2f\n", float64(amps)/100)
		// }
	}
	fmt.Println("")
	svolts := binary.BigEndian.Uint16(results[0:2])
	fmt.Printf("Set volts %2.2f\n", float64(svolts)/100)
	samps := binary.BigEndian.Uint16(results[2:4])
	fmt.Printf("Set Amps %2.2f\n", float64(samps)/100)
	avolts := binary.BigEndian.Uint16(results[4:6])
	fmt.Printf("Actual volts %2.2f\n", float64(avolts)/100)
	aamps := binary.BigEndian.Uint16(results[6:8])
	fmt.Printf("Actual Amps %2.2f\n", float64(aamps)/100)
	apower := binary.BigEndian.Uint16(results[8:10])
	fmt.Printf("Actual power %2.2f\n", float64(apower)/100)
	supplyVolts := binary.BigEndian.Uint16(results[10:12])
	fmt.Printf("Supply Volts %2.2f\n", float64(supplyVolts)/100)

	var on bool
	if results[19] == 01 {
		on = true
	}
	if on {
		r, err := client.WriteSingleRegister(9, 0)
		// fmt.Println(r)
		if err != nil || results == nil {
			log.Fatal(err, r)
		}
	}
	r, err := client.WriteSingleRegister(0, 0)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}

	r, err = client.WriteSingleRegister(6, 1)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}

	r, err = client.WriteSingleRegister(1, 0)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}

	r, err = client.WriteSingleRegister(9, 1)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}
	r, err = client.WriteSingleRegister(1, 2000)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}

	for v := 6.00; v < 13.0; v += 0.5 {

		sv := v * 100
		value := uint16(sv)
		r, err := client.WriteSingleRegister(0, value)
		// fmt.Println(r)
		svolts := binary.BigEndian.Uint16(r)
		fmt.Println(float64(svolts)/100, " volts")
		if err != nil || results == nil {
			log.Fatal(err, r)
		}
		r, err = client.ReadHoldingRegisters(3, 1)
		if err != nil || results == nil {
			log.Fatal(err, r)
		}
		aamps := binary.BigEndian.Uint16(r)
		fmt.Printf("Actual Amps %2.2f\n", float64(aamps)/100)
		// time.Sleep(time.Millisecond * 10)
	}
	start := time.Now()
	for {
		if time.Now().After(start.Add(time.Second * 10)) {
			break
		}
		// r, err := client.ReadHoldingRegisters(3, 1)
		// if err != nil || results == nil {
		// 	log.Fatal(err, r)
		// }
		// aamps := binary.BigEndian.Uint16(r)
		// fmt.Printf("Actual Amps %2.2f\n", float64(aamps)/100)
		results, err := client.ReadHoldingRegisters(0, 30)
		fmt.Println(results)
		// results, err := client.ReadDiscreteInputs(15, 2)
		if err != nil || results == nil {
			log.Fatal(err, results)
		}
		for i := 0; i < len(results); i += 2 {
			fmt.Printf("%2d: %02x %02x ", i, results[i], results[i+1])
			// if k == 3 {
			// 	// fmt.Println("bytes ", results[2:5])
			// 	amps = binary.BigEndian.Uint16(results[2:4])
			// 	fmt.Printf("Amps %2.2f\n", float64(amps)/100)
			// }
		}

		svolts := binary.BigEndian.Uint16(results[0:2])
		fmt.Printf("\nSet volts %2.2f\n", float64(svolts)/100)
		samps := binary.BigEndian.Uint16(results[2:4])
		fmt.Printf("Set Amps %2.2f\n", float64(samps)/100)
		avolts := binary.BigEndian.Uint16(results[4:6])
		fmt.Printf("Actual volts %2.2f\n", float64(avolts)/100)
		aamps := binary.BigEndian.Uint16(results[6:8])
		fmt.Printf("Actual Amps %2.2f\n", float64(aamps)/100)
		apower := binary.BigEndian.Uint16(results[8:10])
		fmt.Printf("Actual power %2.2f\n", float64(apower)/100)
		supplyVolts := binary.BigEndian.Uint16(results[10:12])
		fmt.Printf("Supply Volts %2.2f\n", float64(supplyVolts)/100)

	}
	r, err = client.WriteSingleRegister(9, 0)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}
	r, err = client.WriteSingleRegister(6, 0)
	// fmt.Println(r)
	if err != nil || results == nil {
		log.Fatal(err, r)
	}
}

type Control struct {
	client           modbus.Client
	SetVoltage       float64
	SetCurrent       float64
	ActualVoltage    float64
	ActualCurrent    float64
	Power            float64
	SupplyVoltage    float64
	Locked           bool
	ProtectionTrip   Overload
	Constant         Mode
	DisplayBightness uint16
	Model            uint16
	Version          uint16
}

type Overload uint16

const (
	OverVoltageProtection Overload = 1
	OverCurrentProtection          = 2
	OverPowerProtection            = 3
)

type Mode uint16

const (
	ConstantCurrent Mode = 0
	ConstantVoltage      = 1
)

func (p PowerSupply) GetStatus() error {
	res, err := p.client.ReadHoldingRegisters(0, 12)
	if err != nil {
		return err
	}
}

func floatFromBytes(b []byte) float64 {
	if len(b) != 2 {
		return -1.0
	}
	return float64(binary.BigEndian.Uint16(b)) / 100
}