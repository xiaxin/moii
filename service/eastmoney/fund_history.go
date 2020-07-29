package eastmoney

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
)

// 基金历史
type (
	FenHongType  int
	FenHongValue float64
)

/**
example:
  1. 每份基金份额折算1.012175663份
  2. 每份派现金0.2150元,
  3. 每份基金份额分拆1.162668813份 (拆分后净值一般会变成1,用户持有份额会相应增加)
*/
const (
	FenHongType1 int = 1
	FenHongType2 int = 2
	FenHongType3 int = 3
)

type HistoryResponse struct {
	Data       HistoryResponseList `json:"Data"`
	ErrCode    int                 `json:"ErrCode"`
	ErrMsg     string              `json:"ErrMsg"`
	TotalCount int                 `json:"TotalCount"`
	PageSize   int                 `json:"PageSize"`
	PageIndex  int                 `json:"PageIndex"`
}

type HistoryResponseList struct {
	LSJZList []HistoryResponseData `json:"LSJZList"`
	FundType string                `json:"FundType"`
	//SYType
	//IsNewType
	//Feature
}

type HistoryResponseData struct {
	//  发售日期
	FSRQ string `json:"FSRQ"`
	//  单位净值
	DWJZ string `json:"DWJZ"`
	//  累计净值
	LJJZ string `json:"LJJZ"`
	// SDATE
	// ACTUALSYI
	// NAVTYPE
	// 净值增长率
	JZZZL string `json:"JZZZL"`
	// 申购状态
	SGZT string `json:"SGZT"`
	// 赎回状态
	SHZT string `json:"SHZT"`
	// FHFCZ
	// FHFCBZ
	// DTYPE
	//  分红送配
	FHSP string `json:"FHSP"`
}

/**
@params
  - code string 基金编码
  - page string 分页
*/
func (em *Eastmoney) GetFundValueHistory(code string, page int, callback func(*HistoryResponse, error)) error {

	var err error

	client := em.HttpClient()

	client.OnRequest(func(r *colly.Request) {
		//  请求添加 Referer 突破请求验证
		r.Headers.Set("Referer", fmt.Sprintf("http://fund.eastmoney.com/f10/jjjz_%s.html", code))
	})

	client.OnResponse(func(r *colly.Response) {
		body := r.Body[41 : len(r.Body)-1]

		var data *HistoryResponse

		var err = json.Unmarshal(body, &data)

		callback(data, err)

	})

	url := fmt.Sprintf("http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery18307012244616400747_1582721929389&fundCode=%s&pageIndex=%d&pageSize=20&startDate=&endDate=&_=1582722120962", code, page)

	err = client.Request("GET", url, nil, nil, nil)

	return err
}

func (em *Eastmoney) ParseFHSP(fhsp string) (int, float64, error) {
	var (
		t   int
		v   float64
		err error
	)
	_, err = fmt.Sscanf(fhsp, "每份基金份额折算%f份", &v)

	t = FenHongType1
	if err != nil {
		_, err = fmt.Sscanf(fhsp, "每份派现金%f元", &v)

		t = FenHongType2
		if err != nil {
			t = FenHongType3

			log.Info("c")
			_, err = fmt.Sscanf(fhsp, "每份基金份额分拆%f份", &v)
			if err != nil {
				err = errors.New(fmt.Sprintf("解析基金分红 error %s %s", fhsp, err))
			}
		}
	}

	return t, v, err
}
