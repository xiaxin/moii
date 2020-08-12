package eastmoney

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/xiaxin/moii/time"
	"github.com/xiaxin/moii/web"
)

const (
	// 每年投入次数 245次
	InputCountByYear = 245
)

type Eastmoney struct {
}

func New() *Eastmoney {
	return &Eastmoney{}
}

func (em *Eastmoney) HttpClient() *colly.Collector {
	return web.NewClient()
}

//  获取开放基金净值
func GetAllOpenFundValue(callback func(r *colly.Response)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(callback)

	client.OnRequest(func(r *colly.Request) {
		//  请求添加 Referer 突破请求验证
		r.Headers.Set("Referer", "http://fund.eastmoney.com/data/fundranking.html")
	})

	date := time.GetTodayYMD()

	url := fmt.Sprintf("http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=all&rs=&gs=0&sc=zzf&st=desc&sd=%s&ed=%s&qdii=&tabSubtype=,,,,,&pi=1&pn=7000&dx=1&v=0.3971094993995896", date, date)

	err = client.Request("GET", url, nil, nil, nil)

	return err
}

//  获取货币基金净值 TODO
func GetAllCurrencyFundValue(callback func(r *colly.Response)) error {
	return nil
}

//  获取基金公司
func GetAllFundCompany(callback func(r *colly.Response)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(callback)

	err = client.Request("GET", "http://fund.eastmoney.com/Data/FundRankScale.aspx", nil, nil, nil)

	return err
}

//  获取基金列表
func GetAllFundList(callback func(r *colly.Response)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(callback)

	err = client.Request("GET", "http://fund.eastmoney.com/js/fundcode_search.js", nil, nil, nil)

	return err
}

//  获取基金公司 基金列表 -
func GetFundCompanyFundList(code string, callback func(r *colly.Response)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(callback)

	url := fmt.Sprintf("http://fund.eastmoney.com/Company/home/KFSFundNet?gsid=%s&fundType=", code)

	err = client.Request("GET", url, nil, nil, nil)

	return err
}

func GetFundCompanyFundListHtml(code string, callback func(r *colly.Response)) error {
	var err error

	client := web.NewClient()

	client.OnResponse(callback)

	url := fmt.Sprintf("http://fund.eastmoney.com/Company/%s.html", code)

	err = client.Request("GET", url, nil, nil, nil)

	return err
}

//  TODO 请求callback字符串 动态修改
// // Deprecated: 方法废弃
func GetFundValueHistory(code string, page int, callback func(r *colly.Response)) error {

	var err error

	client := web.NewClient()

	client.OnRequest(func(r *colly.Request) {
		//  请求添加 Referer 突破请求验证
		r.Headers.Set("Referer", fmt.Sprintf("http://fund.eastmoney.com/f10/jjjz_%s.html", code))
	})

	client.OnResponse(callback)

	url := fmt.Sprintf("http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery18307012244616400747_1582721929389&fundCode=%s&pageIndex=%d&pageSize=20&startDate=&endDate=&_=1582722120962", code, page)

	err = client.Request("GET", url, nil, nil, nil)

	return err
}

func GetFundDetail(code string, callback func(detail FundDetail)) error {
	var err error

	client := web.NewClient()

	client.OnHTML(".info", func(e *colly.HTMLElement) {
		var text string
		fundDetail := FundDetail{
			Code: code,
		}

		//  元素
		tr := e.DOM.Find("tr")

		// 基金类型
		fundDetail.TypeName = tr.Eq(1).Find("td").Eq(1).Text()
		// TODO fundDetail.Type

		// 资产规模
		text = tr.Eq(3).Find("td").Eq(0).Text()
		fmt.Sscanf(text, "%f亿元（", &fundDetail.TotalMoney)

		//  基金托管人
		fundDetail.Bank = tr.Eq(4).Find("td").Eq(1).Text()

		callback(fundDetail)
	})

	tmp := "http://fundf10.eastmoney.com/jbgk_%s.html"

	err = client.Visit(fmt.Sprintf(tmp, code))
	return err
}

