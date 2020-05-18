package web

import (
	"github.com/gocolly/colly/v2"
)

var (
	client = colly.NewCollector()
)

const (
	UserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"
	URLRevisit = true
	Async      = false
)

func init() {
	// 浏览器
	client.UserAgent = UserAgent
	// 允许重复访问（开启缓存可加速）
	client.AllowURLRevisit = URLRevisit
	// 同步
	client.Async = Async

	client.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
}

func NewClient() *colly.Collector {
	client := client.Clone()
	return client
}
