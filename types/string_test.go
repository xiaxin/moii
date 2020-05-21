package types

import (
	"math"
	"testing"
)

func TestString_UInt64(t *testing.T) {
	var UInt64 String = "18446744073709551615"

	if v, err := UInt64.UInt64(); nil != err || v != math.MaxUint64 {
		t.Errorf("error:%s value:%d", err, v)
	}
}

func TestString_Int64(t *testing.T) {
	var UInt64 String = "9223372036854775807"

	if v, err := UInt64.Int64(); nil != err || v != math.MaxInt64 {
		t.Errorf("error:%s value:%d", err, v)
	}
}

func TestString_Bool(t *testing.T) {
	var TBool String = "true"

	if v, err := TBool.Bool(); nil != err || v != true {
		t.Errorf("error:%s value:%T", err, v)
	}

	var FBool String = "false"

	if v, err := FBool.Bool(); nil != err || v != false {
		t.Errorf("error:%s value:%T", err, v)
	}

	var EBool String = "error value"

	if v, err := EBool.Bool(); nil != err || v != false {
		t.Errorf("error:`%s` type:%T value:%v", err, EBool, EBool)
	}
}
