// Code generated by "stringer -type=Lock"; DO NOT EDIT.

package main

import "fmt"

const _Lock_name = "UnlockedLocked"

var _Lock_index = [...]uint8{0, 8, 14}

func (i Lock) String() string {
	if i >= Lock(len(_Lock_index)-1) {
		return fmt.Sprintf("Lock(%d)", i)
	}
	return _Lock_name[_Lock_index[i]:_Lock_index[i+1]]
}
