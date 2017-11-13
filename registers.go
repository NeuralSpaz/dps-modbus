package main

// Note 1: This product namely Production equipment M0-M9 a total of 10 groups of
// storage data groups, each group has a total of 10-10-17 data, of which M0 The
// data group is the data group called by default for the products to be powered on.

// The data groups M1 and M2 are quickly called out for the product panel, and the
// M3- M9 is an ordinary storage array, and the starting address of a data group is
// calculated as 0050H + data group number * 0010H, for example: M3 The starting
// address of the data set is: 0050H + 3 * 0010H = 0080H.

// Note 2: Key lock function to write values 0 and 0 for non-locking, a lock.
// Note 3: 0-3,0-protected state is normal read operation, one of OVP, 2 as the
// OCP, 3 to OPP. Note 4: Constant Current Status of read values CV 0 state and 0,
// CC 1 for the state. Note 5: open Write function Off lose Readings of 0 and 1,
// 0 are off and 1 is on. Note 6: The backlight brightness level range of 0-5,0
// read darkest, level 5 the brightest. level Note 5: fast recall function data
// set is written to 0-9, after the corresponding write data is automatically
// transferred group of data

const (
	voltageSetRegister    = 0
	currentSetRegister    = 1
	voltageOutRegister    = 2
	powerOutRegister      = 3
	supplyVoltageRegister = 4
	lockRegister          = 6
	protectionRegister    = 7
	modeRegister          = 8 // Constant Current or Voltage
	onOffRegister         = 9
	ledBrightnessRegister = 10
	modelRegister         = 11
	versionRegister       = 12

	loadMemoryRegister = 35 // loads presets

	defaultVoltageSetRegister             = 80
	defaultCurrentSetRegister             = 81
	defaultOverVoltageProtectionRegister  = 82
	defaultOverCurrentProtectionRegister  = 83
	defaultOverPowerProtectionRegister    = 84
	defaultLedBrightnessRegister          = 85
	defaultDataRecallRegister             = 86
	defaultpowerOutputSwitchStateRegister = 87

	m1VoltageSetRegister             = 0x60
	m1CurrentSetRegister             = 0x61
	m1OverVoltageProtectionRegister  = 0x62
	m1OverCurrentProtectionRegister  = 0x63
	m1OverPowerProtectionRegister    = 0x64
	m1LedBrightnessRegister          = 0x65
	m1DataRecallRegister             = 0x66
	m1powerOutputSwitchStateRegister = 0x67

	m2VoltageSetRegister             = 0x70
	m2CurrentSetRegister             = 0x71
	m2OverVoltageProtectionRegister  = 0x72
	m2OverCurrentProtectionRegister  = 0x73
	m2OverPowerProtectionRegister    = 0x74
	m2LedBrightnessRegister          = 0x75
	m2DataRecallRegister             = 0x76
	m2powerOutputSwitchStateRegister = 0x77

	m3VoltageSetRegister             = 0x80
	m3CurrentSetRegister             = 0x81
	m3OverVoltageProtectionRegister  = 0x82
	m3OverCurrentProtectionRegister  = 0x83
	m3OverPowerProtectionRegister    = 0x84
	m3LedBrightnessRegister          = 0x85
	m3DataRecallRegister             = 0x86
	m3powerOutputSwitchStateRegister = 0x87

	m4VoltageSetRegister             = 0x90
	m4CurrentSetRegister             = 0x91
	m4OverVoltageProtectionRegister  = 0x92
	m4OverCurrentProtectionRegister  = 0x93
	m4OverPowerProtectionRegister    = 0x94
	m4LedBrightnessRegister          = 0x95
	m4DataRecallRegister             = 0x96
	m4powerOutputSwitchStateRegister = 0x97

	m5VoltageSetRegister             = 0xa0
	m5CurrentSetRegister             = 0xa1
	m5OverVoltageProtectionRegister  = 0xa2
	m5OverCurrentProtectionRegister  = 0xa3
	m5OverPowerProtectionRegister    = 0xa4
	m5LedBrightnessRegister          = 0xa5
	m5DataRecallRegister             = 0xa6
	m5powerOutputSwitchStateRegister = 0xa7

	m6VoltageSetRegister             = 0xb0
	m6CurrentSetRegister             = 0xb1
	m6OverVoltageProtectionRegister  = 0xb2
	m6OverCurrentProtectionRegister  = 0xb3
	m6OverPowerProtectionRegister    = 0xb4
	m6LedBrightnessRegister          = 0xb5
	m6DataRecallRegister             = 0xb6
	m6powerOutputSwitchStateRegister = 0xb7

	m7VoltageSetRegister             = 0xc0
	m7CurrentSetRegister             = 0xc1
	m7OverVoltageProtectionRegister  = 0xc2
	m7OverCurrentProtectionRegister  = 0xc3
	m7OverPowerProtectionRegister    = 0xc4
	m7LedBrightnessRegister          = 0xc5
	m7DataRecallRegister             = 0xc6
	m7powerOutputSwitchStateRegister = 0xc7

	m8VoltageSetRegister             = 0xd0
	m8CurrentSetRegister             = 0xd1
	m8OverVoltageProtectionRegister  = 0xd2
	m8OverCurrentProtectionRegister  = 0xd3
	m8OverPowerProtectionRegister    = 0xd4
	m8LedBrightnessRegister          = 0xd5
	m8DataRecallRegister             = 0xd6
	m8powerOutputSwitchStateRegister = 0xd7

	m9VoltageSetRegister             = 0xe0
	m9CurrentSetRegister             = 0xe1
	m9OverVoltageProtectionRegister  = 0xe2
	m9OverCurrentProtectionRegister  = 0xe3
	m9OverPowerProtectionRegister    = 0xe4
	m9LedBrightnessRegister          = 0xe5
	m9DataRecallRegister             = 0xe6
	m9powerOutputSwitchStateRegister = 0xe7
)
