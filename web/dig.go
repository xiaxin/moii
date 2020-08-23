package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaxin/moii/dig"
)

func NewDig(ctx *gin.Context) *dig.Dig {
	dig := dig.New(nil)
	return dig
}

// 中间件
func NewDigMiddleware(dig *dig.Dig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO
		ctx.Set(KeyDig, dig)
		ctx.Next()
	}
}
