package time

import (
	"context"
	"database/sql/driver"
	"strconv"
	xtime "time"
)

// TODO 待增加
const (
	FormatYmd = "2006-01-02"
	FormatHis = "15:04:05"

	FormatYear      = "2006"
	FormatMonth     = "1"
	FormatZeroMonth = "01"
	FormatDay       = "2"
	FormatZeroDay   = "02"
	//  TODO 24小时制时 是否有0
	// Zero
	FormatHour       = "15"
	FormatHour12     = "3"
	FormatZeroHour12 = "03"
)

// Time be used to MySql timestamp converting.
type Time int64

// Scan scan time.
func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case xtime.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

// Value get time value.
func (jt Time) Value() (driver.Value, error) {
	return xtime.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt Time) Time() xtime.Time {
	return xtime.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration xtime.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := xtime.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := xtime.Until(deadline); ctimeout < xtime.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, xtime.Duration(d))
	return d, ctx, cancel
}

func GetTodayYMD() string {
	return xtime.Now().Format(FormatYmd)
}

func GetTodayHIS() string {
	return xtime.Now().Format(FormatHis)
}
func GetYestodayYMD() string {
	return xtime.Now().AddDate(0, 0, -1).Format(FormatYmd)
}

func GetYear() string {
	return xtime.Now().Format(FormatYear)
}

func GetMonth() string {
	return xtime.Now().Format(FormatMonth)
}

func GetDay() string {
	return xtime.Now().Format(FormatDay)
}

func GetHour() string {
	return xtime.Now().Format(FormatHour)
}

func Get2019Year() (xtime.Time, xtime.Time) {
	startDate, _ := xtime.Parse(FormatYmd, "2019-01-01")
	endDate, _ := xtime.Parse(FormatYmd, "2019-12-31")

	return startDate, endDate
}

func Get2018Year() (xtime.Time, xtime.Time) {
	startDate, _ := xtime.Parse(FormatYmd, "2018-01-01")
	endDate, _ := xtime.Parse(FormatYmd, "2018-12-31")

	return startDate, endDate
}
