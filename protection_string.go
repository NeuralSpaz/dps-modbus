// Code generated by "stringer -type=Protection"; DO NOT EDIT.

package main

import "fmt"

const _Protection_name = "OverVoltageProtectionOverCurrentProtectionOverPowerProtection"

var _Protection_index = [...]uint8{0, 21, 42, 61}

func (i Protection) String() string {
	i -= 1
	if i >= Protection(len(_Protection_index)-1) {
		return fmt.Sprintf("Protection(%d)", i+1)
	}
	return _Protection_name[_Protection_index[i]:_Protection_index[i+1]]
}