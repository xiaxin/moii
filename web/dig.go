package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaxin/moii/dig"
)

const (
	KeyDi = "middleware#di"
)

func NewDig(ctx *gin.Context) *dig.Dig {
	dig := dig.New(nil)
	return dig
}

// 中间件
func NewDigMiddleware(dig *dig.Dig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(KeyDi, dig)
		ctx.Next()
	}
}

//  从上下文获取中间件
func GetDig(ctx *gin.Context) (*dig.Dig, bool) {
	diI, _ := ctx.Get(KeyDi)
	di, ok := diI.(*dig.Dig)
	exists := di != nil && ok
	return di, exists
}
