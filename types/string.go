package types

import (
	"strconv"
	"time"
)

type String string

func (s String) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s String) Int64() (int64, error) {
	return strconv.ParseInt(s.String(), 10, 64)
}

func (s String) UInt64() (uint64, error) {
	return strconv.ParseUint(s.String(), 10, 64)
}

func (s String) Float64() (float64, error) {
	return strconv.ParseFloat(s.String(), 64)
}

func (s String) Bool() (bool, error) {
	return strconv.ParseBool(s.String())
}

func (s String) String() string {
	return string(s)
}

func (s String) Time(format string) (time.Time, error) {
	tm, err := time.Parse(format, s.String())
	if nil != err {
		return time.Time{}, err
	}
	return tm, nil
}
