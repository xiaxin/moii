package web

import (
	"github.com/gocolly/colly/v2"
	"testing"
)

func TestClient(t *testing.T) {

	client := NewClient()

	client.OnResponse(func(r *colly.Response) {
	})

	err := client.Visit("http://www.baidu.com")

	if nil != err {
		t.Error(err)
	}
}
