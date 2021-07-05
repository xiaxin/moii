package column

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	xtime "github.com/xiaxin/moii/time"
)

// DateTime Mysql Column datetime
type DateTime struct {
	time.Time
}

func (t DateTime) String() string {
	return t.Format(xtime.FormatYmdHis)
}

// UnmarshalJSON  []byte to struct val
func (t *DateTime) UnmarshalJSON(data []byte) error {
	value := string(data)
	float, err := strconv.ParseInt(value, 10, 32)

	if nil != err {
		return err
	}

	time := time.Unix(float, 0)
	return t.Scan(time)
}

// MarshalJSON  json {data} to []byte
func (t DateTime) MarshalJSON() ([]byte, error) {
	//  返回的字符串需要加引号
	datetime := fmt.Sprintf(`"%s"`, t.String())
	return []byte(datetime), nil
}

// Value return interface{}
func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan data to mysql
func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Now 设置当前时间
func (t *DateTime) Now() {
	t.Time = time.Now()
}
