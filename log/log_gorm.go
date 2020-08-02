package log

import (
	"database/sql/driver"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"fmt"
	"time"
	"unicode"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type GormLogger struct {
	*Logger
}

func (l *Logger) NewGorm() *GormLogger {
	logger := l.Named("gorm").Clone(zap.AddCallerSkip(8))

	gorm := &GormLogger{
		Logger: logger,
	}

	return gorm
}

func (l *GormLogger) Print(values ...interface{}) {
	l.Println(values)
}

func (l *GormLogger) Println(values []interface{}) {
	gorm := createGormLog(values)

	fields := gorm.toZapFields()
	l.Debugw("_gorm_query", fields...)
}

type gormLog struct {
	occurredAt time.Time
	source     string
	duration   string
	sql        string
	values     []string
	other      []string
}

func (l *gormLog) GetSql() string {
	return l.sql
}

func (l *gormLog) GetDuration() string {
	return l.duration
}

func (l *gormLog) toZapFields() []interface{} {
	if len(l.other) == 0 {
		return []interface{}{
			//zap.Any("source", l.source),
			//  SQL
			//  SQL
			zap.Any("sql", l.sql),
			//  执行时间
			zap.Any("duration", l.duration),
			//  SQL 中的值
			//zap.Any("values", l.values),
			//  其他
			//zap.Any("other", l.other),
		}
	} else {
		return []interface{}{
			//zap.Any("source", l.source),
			zap.Any("other", l.other),
		}
	}
}

func createGormLog(values []interface{}) *gormLog {
	ret := &gormLog{}
	ret.occurredAt = gorm.NowFunc()

	if len(values) > 1 {
		var level = values[0]
		ret.source = getSource(values)

		if level == "sql" {
			//  SQL 运行时间
			ret.duration = getDuration(values)
			ret.values = getFormattedValues(values)

			//  处理SQL

			ret.sql = getSql(values)
		} else {
			ret.other = append(ret.other, fmt.Sprint(values[2:]))
		}
	}

	return ret
}

func getSource(values []interface{}) string {
	return fmt.Sprint(values[1])
}

func getDuration(values []interface{}) string {
	return fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)
}

func getFormattedValues(values []interface{}) []string {
	rawValues := values[4].([]interface{})
	formattedValues := make([]string, 0, len(rawValues))
	for _, value := range rawValues {
		switch v := value.(type) {
		case time.Time:
			formattedValues = append(formattedValues, fmt.Sprint(v))
		case []byte:
			if str := string(v); isPrintable(str) {
				formattedValues = append(formattedValues, fmt.Sprint(str))
			} else {
				formattedValues = append(formattedValues, "<binary>")
			}
		default:
			str := "NULL"
			if v != nil {
				str = fmt.Sprint(v)
			}
			formattedValues = append(formattedValues, str)
		}
	}
	return formattedValues
}

func getSql(values []interface{}) string {
	var (
		sql             string
		formattedValues []string
	)
	for _, value := range values[4].([]interface{}) {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				switch value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
					formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
				default:
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
		} else {
			formattedValues = append(formattedValues, "NULL")
		}
	}

	if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
		sql = values[3].(string)
		for index, value := range formattedValues {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
			sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(formattedValues)
		for index, value := range sqlRegexp.Split(values[3].(string), -1) {
			sql += value
			if index < formattedValuesLength {
				sql += formattedValues[index]
			}
		}
	}

	return sql
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
