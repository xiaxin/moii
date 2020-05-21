package types

import "strconv"

type String string

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
