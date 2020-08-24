package eastmoney

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github.com/xiaxin/moii/web"
)

type FundDetail struct {
	Name       string
	Alias      string
	Code       string
	Type       int
	TypeName   string   //  类型名称
	TotalMoney float64  //  规模
	Bank       string   //  银行
	BuyState       string
	BuyStateRemark string //  buy_state 限额时 限额事项写在这里  buy_state 为 开放状态
	SellState  string
}

//func (em *Eastmoney) GetFundDetail(code string, callback func(detail FundDetail)) error {
//	var err error
//
//	client := web.NewClient()
//
//	client.OnHTML(".info", func(e *colly.HTMLElement) {
//		var text string
//		fundDetail := FundDetail{
//			Code: code,
//		}
//
//		//  TODO 补充交易状态
//
//		//  元素
//		tr := e.DOM.Find("tr")
//
//		// 基金类型
//		fundDetail.TypeName = tr.Eq(1).Find("td").Eq(1).Text()
//		// TODO fundDetail.Type
//
//		// 资产规模
//		text = tr.Eq(3).Find("td").Eq(0).Text()
//		fmt.Sscanf(text, "%f亿元（", &fundDetail.TotalMoney)
//
//		//  基金托管人
//		fundDetail.Bank = tr.Eq(4).Find("td").Eq(1).Text()
//
//		callback(fundDetail)
//	})
//
//	tmp := "http://fundf10.eastmoney.com/jbgk_%s.html"
//
//	err = client.Visit(fmt.Sprintf(tmp, code))
//	return err
//}

func (em *Eastmoney) GetFundList(callback func(data [][]string)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(func(r *colly.Response) {
		body := r.Body[10 : len(r.Body)-1]

		var data [][]string

		json.Unmarshal(body, &data)

		callback(data)
	})

	err = client.Request("GET", "http://fund.eastmoney.com/js/fundcode_search.js", nil, nil, nil)

	return err
}
