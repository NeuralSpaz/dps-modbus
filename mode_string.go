// Code generated by "stringer -type=Mode"; DO NOT EDIT.

package main

import "fmt"

const _Mode_name = "ConstantVoltageConstantCurrent"

var _Mode_index = [...]uint8{0, 15, 30}

func (i Mode) String() string {
	if i >= Mode(len(_Mode_index)-1) {
		return fmt.Sprintf("Mode(%d)", i)
	}
	return _Mode_name[_Mode_index[i]:_Mode_index[i+1]]
}
