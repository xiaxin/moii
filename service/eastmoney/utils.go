package eastmoney

import (
	"time"
)

var monthDays = [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// 计算定投日
func IsPeriodDay(periodType PeriodType, periodValue PeriodValue, day, from, to time.Time) bool {

	//  TODO 周六周日 无法定投
	if periodType == Week {
		if time.Weekday(periodValue) == time.Sunday || time.Weekday(periodValue) == time.Saturday {
			return false
		}

		if day.Weekday() == time.Weekday(periodValue) {
			return true
		}
		return false
	}
	if periodType == Week2 { // 每两周定投一次
		if day.Weekday() == from.Weekday() && int(day.Sub(from).Hours()/24/7)%2 == 0 {
			return true
		}
		return false
	}
	if periodType == Month { // 每月定投
		if day.Day() == from.Day() {
			return true
		}
		dayMonthDayCnt := GetMonthDayCount(day.Year(), int(day.Month())) // from那月最多有多少天
		if dayMonthDayCnt < from.Day() && day.Day() == dayMonthDayCnt {
			// 如果day 为本月最后一天, 且起投日对应的那在 在本月无对应日,则以本月最后一天为定投日
			// 比如,定投日为31号, 而2月4月等根本没有31号,则以当月最后一天为定投日
			return true
		}

		return false
	}
	return false
}

func GetMonthDayCount(year, month int) int { // month[1~12]
	//  判断 2 月份是否是闰年
	if month == 2 && IsLeapYear(year) {
		return 29
	}
	return monthDays[month]
}

func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}
