package pagination

import (
	"github.com/jinzhu/gorm"
	"math"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

// 分页参数
type Param struct {
	DB      *gorm.DB
	Page    int      `form:"page"  json:"page"`
	Limit   int      `form:"limit" json:"limit"`
	OrderBy []string `form:"order" json:"order"`
	ShowSQL bool
}

// 分页结果
type Result struct {
	//  数据
	Data interface{} `json:"data"`
	//  分页数据
	Meta *ResultMeta `json:"meta"`
}

func (r *Result) GetData() interface{} {
	return r.Data
}

func (r *Result) GetMeta() *ResultMeta {
	return r.Meta
}

type ResultMeta struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	CurrentPage int `json:"current_page"`
	PrevPage    int `json:"prev_page"`
	NextPage    int `json:"next_page"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// Paging 分页
func Apply(p *Param, out interface{}) (res *Result, err error) {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}

	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.Limit == 0 {
		p.Limit = DefaultLimit
	}

	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)

	var count int
	var offset int

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	//  异步获取 Count 逻辑
	go getCount(db, out, done, &count)

	err = db.Limit(p.Limit).Offset(offset).Find(out).Error

	// goroutine done
	<-done

	//  结果设置
	result := &Result{
		Data: out,
		Meta: &ResultMeta{
			Total:       int(math.Ceil(float64(count) / float64(p.Limit))),
			Count:       count,
			CurrentPage: p.Page,
			Offset:      offset,
			Limit:       p.Limit,

			//  设置默认值
			PrevPage: p.Page,
			NextPage: p.Page + 1,
		},
	}

	//  结果数据校准
	if p.Page > 1 {
		result.Meta.PrevPage = p.Page - 1
	}

	if p.Page == result.Meta.Total {
		result.Meta.NextPage = p.Page
	}
	return result, err
}

func getCount(db *gorm.DB, out interface{}, done chan bool, count *int) {
	db.Model(out).Count(count)
	done <- true
}
