package column

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	xtime "github.com/xiaxin/moii/time"
)

// Date Mysql Column date
type Date struct {
	time.Time
}

// NewDate 创建Date
func NewDate(val string) (Date, error) {
	date, err := time.Parse(xtime.FormatYmd, val)

	if nil != err {
		return Date{}, err
	}

	return Date{date}, nil
}

func (t Date) String() string {
	return t.Format(xtime.FormatYmd)
}

// UnmarshalJSON  []byte to struct val
func (t *Date) UnmarshalJSON(data []byte) error {
	value := string(data)
	float, err := strconv.ParseInt(value, 10, 32)

	if nil != err {
		return err
	}

	time := time.Unix(float, 0)
	return t.Scan(time)
}

// MarshalJSON  json {data} to []byte
func (t Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, t.String())
	return []byte(formatted), nil
}

// Value return interface{}
func (t Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Format(xtime.FormatYmd), nil
}

// Scan data to mysql
func (t *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Date{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Now 设置当前时间
func (t *Date) Now() {
	t.Time = time.Now()
}

// SetString 设置字符串
func (t *Date) SetString(v string) error {
	date, err := time.Parse(xtime.FormatYmd, v)

	if nil != err {
		return err
	}

	t.Time = date
	return nil
}
