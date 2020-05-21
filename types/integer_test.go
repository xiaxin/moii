package types

import (
	"math"
	"testing"
)

func TestInt64_String(t *testing.T) {
	var I64 Int64 = math.MaxInt64

	if I64.String() != "9223372036854775807" {
		t.Error("Int64 Error")
	}

	var UI64 UInt64 = math.MaxUint64

	if UI64.String() != "18446744073709551615" {
		t.Error("UInt64 Error")
	}
}

func TestUInt64_Datetime(t *testing.T) {

	var UI64Max UInt64 = math.MaxInt64
	if UI64Max.Datetime() != "292277026596-12-04 23:30:07" {
		t.Error("UI64Max Error")
	}

	var UI64Min UInt64 = 0
	if UI64Min.Datetime() != "1970-01-01 08:00:00" {
		t.Error("UI64Max Error")
	}
}
