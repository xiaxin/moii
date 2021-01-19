package column

import (
	"database/sql/driver"
	"fmt"
	"time"

	xtime "github.com/xiaxin/moii/time"
)

// Date Mysql Column date
type Date struct {
	time.Time
}

// NewDate 创建Date
func NewDate(val string) (*Date, error) {
	date, err := time.Parse(xtime.FormatYmd, val)

	if nil != err {
		return nil, err
	}

	return &Date{date}, nil
}

// Set 设置字符串
func (t *Date) Set(val string) error {
	date, err := NewDate(val)

	if nil != err {
		return err
	}

	t.Time = date.Time

	return nil
}

// Get 获取字符串
func (t *Date) Get() string {
	return t.String()
}

func (t *Date) String() string {
	return t.Format(xtime.FormatYmd)
}

// MarshalJSON  json {data} to []byte
func (t Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(xtime.FormatYmd))
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
