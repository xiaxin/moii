package column

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// Unix Mysql Column timestamp
type Unix struct {
	time.Time
}

// UnmarshalJSON  []byte to struct val
func (t *Unix) UnmarshalJSON(data []byte) error {
	value := string(data)
	float, err := strconv.ParseInt(value, 10, 32)

	if nil != err {
		return err
	}

	time := time.Unix(float, 0)
	return t.Scan(time)
}

// MarshalJSON  json {data} to []byte
func (t Unix) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, t.String())
	return []byte(formatted), nil
}

// Value return interface{}
func (t Unix) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Unix(), nil
}

// Scan data to mysql
func (t *Unix) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Unix{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
