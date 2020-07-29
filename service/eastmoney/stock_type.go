package eastmoney

type (
	StockType      string
	StockModelType int
)

// 上海 m:1
// 上证A  m:1+t:2,m:1+t:23

// 深圳 m:0
// 深圳A  m:0+t:6,m:0+t:13,m:0+t:80
// 创业板 m:0+t:80
// 中小板 m:0 t:13

// B股 fs: m:0 t:7,m:1 t:3

const (
	StockTypeA           StockType = "m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23"
	StockTypeIndexSH     StockType = "m:1+s:2"
	StockTypeIndexSZ     StockType = "m:0+t:5"
	StockRequestCallback           = "jQuery11240510512757979825_1583668617565"
	StockRequestPageSize           = 20

	//  未知
	StockModelTypeD StockModelType = 0
	//  股票
	StockModelTypeS StockModelType = 1
	//  基金
	StockModelTypeI StockModelType = 2
)

var StockTypeDict = map[StockType]StockModelType{
	StockTypeA:       StockModelTypeS,
	StockTypeIndexSH: StockModelTypeI,
	StockTypeIndexSZ: StockModelTypeI,
}

func (st StockType) ToModelType() StockModelType {
	if v, ok := StockTypeDict[st]; ok {
		return v
	}
	return StockModelTypeD
}
