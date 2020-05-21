package types

import (
	xtime "github.com/xiaxin/moii/time"
	"strconv"
	"time"
)

type Int64 int64

func (i Int64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int64) Datetime() string {
	return time.Unix(int64(i), 0).Format(xtime.FormatYmdHis)
}

type UInt64 uint64

func (i UInt64) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i UInt64) Datetime() string {
	return time.Unix(int64(i), 0).Format(xtime.FormatYmdHis)
}
