package eastmoney

import "time"

// 定投类型
type (
	PeriodType  int
	PeriodValue int
)

const (
	// 每日
	Day PeriodType = 4
	// 每周定投一次
	Week PeriodType = 1
	// 每2周定投一次
	Week2 PeriodType = 2
	// 每月
	Month PeriodType = 3
)

//  定投请求
type PeriodRequest struct {
	Code      string
	StartDate time.Time
	EndDate   time.Time
	// 定投周期
	PeriodType  PeriodType
	PeriodValue PeriodValue

	//  投入金额
	InputMoney float64
	InputCount float64
}
