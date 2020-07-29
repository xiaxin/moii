package eastmoney

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
)

type StockPriceDataResponse struct {
	Data StockPriceData `json:"data"`
}

type StockPriceData struct {
	Total int             `json:"total"`
	Diff  []StockDataDiff `json:"diff"`
}

type StockDataDiff struct {
	//f1,
	//f2,
	//f3,
	//f4,
	//f5,
	//f6,
	//f7,
	//f8,
	// 市盈率 Price Earnings Ratio
	PERatio interface{} `json:"f9"`
	//f10,
	//  编码
	Code string `json:"f12"`
	//f13,
	Name string `json:"f14"`
	//f15,
	//f16,
	//f17,
	//f18,
	//f20,
	//f21,
	// 市净率 Price-to-Book Ratio
	PBRatio interface{} `json:"f23"`
	//f24,
	//f25,
	//f26,
	//f22,
	//f33,
	//f11,
	//f62,
	//f128,
	//f136,
	//f115,
	//f152
}

// 股票指数
func (em *Eastmoney) GetStockListData(stockPriceType StockType, page int, callback func(page int, stockPriceType StockType, res StockPriceDataResponse, err error)) error {

	var (
		err error
		tmp = "http://1.push2.eastmoney.com/api/qt/clist/get?cb=%s&pn=%d&pz=%d&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=%s&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152&_=1583668617647"
	)

	client := em.HttpClient()

	client.OnRequest(func(r *colly.Request) {
	})

	client.OnResponse(func(r *colly.Response) {
		body := r.Body[len(StockRequestCallback)+1 : len(r.Body)-2]

		log.Info(string(body))

		var res StockPriceDataResponse

		err := json.Unmarshal(body, &res)

		callback(page, stockPriceType, res, err)

		//  TODO 成功解析数据， 进行下一页采集。
		//if nil == err {
		//	total := res.Data.Total
		//
		//	count := math.Ceil(float64(total) / StockPricePageSize)
		//
		//	if float64(page) < count {
		//		//em.GetStockIndexData(stockPriceType, page+1, GetStockCallback(stockType))
		//	}
		//}
	})
	// TODO OnError
	url := fmt.Sprintf(tmp, StockRequestCallback, page, StockRequestPageSize, stockPriceType)
	err = client.Request("GET", url, nil, nil, nil)

	return err
}
