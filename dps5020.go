package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goburrow/modbus"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats"
)

func dataLogger(dbname string, user string, password string, server string) {
	go func() {
		servers := "nats://127.0.0.1:4222"
		hostname, _ := os.Hostname()
		name := nats.Name(hostname + "logger")
		nc, err := nats.Connect(servers, name)
		if err != nil {
			log.Fatalln(err)
		}
		c, _ := nats.NewEncodedConn(nc, "json")
		defer c.Close()
		suber := make(chan Status)
		c.BindRecvChan("CellStatus", suber)
		dbConnectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", user, password, server, dbname)

		db, err := sqlx.Open("mysql", dbConnectString)
		if err != nil {
			log.Fatalln(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalln(err)
		}
		inserttmpl := fmt.Sprintf("INSERT INTO %s.cell (ts,SetVoltage,SetCurrent,ActualVoltage,ActualCurrent,Power,SupplyVoltage,ProtectionTrip,Constant,OutputOn) VALUES (?,?,?,?,?,?,?,?,?,?)", dbname)

		stmt, err := db.Preparex(inserttmpl)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			select {
			case s := <-suber:
				fmt.Println("chan: ", s)
				ts := time.Now()

				_, err := stmt.Exec(ts, s.SetVoltage, s.SetCurrent, s.ActualVoltage, s.ActualCurrent, s.Power, s.SupplyVoltage, s.ProtectionTrip, s.Constant, s.OutputOn)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

}

func main() {
	fmt.Println("starting dps5020 monitor in modbus mode")
	dbuser := os.Getenv("DPSUSER")
	dbpassword := os.Getenv("DPSPASS")
	dbname := os.Getenv("DPSDB")
	dbconn := os.Getenv("DPSDBCONN")
	dataLogger(dbname, dbuser, dbpassword, dbconn)
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	// handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	servers := "nats://127.0.0.1:4222"
	hostname, _ := os.Hostname()
	name := nats.Name(hostname)
	nc, err := nats.Connect(servers, name)
	if err != nil {
		log.Fatalln(err)
	}
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()
	puber := make(chan Status, 10)
	c.BindSendChan("CellStatus", puber)

	dps := new(DPS)
	dps.conn = client
	dps.puber = puber

	dps.readPresets()

	dps.readStatus()
	// time.Sleep(time.Second * 1)
	targetVoltage := 2.0

	if err := dps.setVoltage(targetVoltage); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := dps.setCurrent(4.0); err != nil {
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
		dps.readStatus()
		dps.RLock()
		fmt.Println(dps.Statuz)
		dps.RUnlock()
	}

	if err := dps.disableOutput(); err != nil {
		log.Println(err)
	}
}

func (d *DPS) readStatus() error {
	d.Lock()
	defer d.Unlock()
	statusRaw, err := d.conn.ReadHoldingRegisters(0, 12)
	if err != nil {
		return err
	}
	d.Statuz = parseStatus(statusRaw)
	go func() { d.puber <- d.Statuz }()

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
	resp, err := d.conn.WriteSingleRegister(voltageSetRegister, uint16(sv*100))
	setVoltage := floatFromBytes(resp)
	d.Statuz.SetVoltage = setVoltage
	if err != nil {
		return fmt.Errorf("failed to set voltage")
	}
	return err
}

func (d *DPS) setCurrent(sc float64) error {
	resp, err := d.conn.WriteSingleRegister(currentSetRegister, uint16(sc*100))
	scResp := floatFromBytes(resp)
	d.Statuz.SetCurrent = scResp
	if err != nil {
		return fmt.Errorf("failed to set voltage")
	}
	return err
}
