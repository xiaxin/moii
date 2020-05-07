package web

import "github.com/gocolly/colly/v2"

var (
	client = colly.NewCollector()
)

func init() {
	client.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
}

func NewClient() *colly.Collector {
	return client.Clone()
}
