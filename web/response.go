package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Json(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, data)
}

func Json200(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, data)
}

func HTML(ctx *gin.Context, format string, values ...interface{}) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, format, values...)
}

func ValidateFailed(ctx *gin.Context, err interface{}) {
	Json500(ctx, gin.H{
		"errno": 2,
		"error": err,
	})
}

func Json500(ctx *gin.Context, data interface{}) {
	ctx.JSON(500, data)
}

func JsonError(ctx *gin.Context, err error) {
	Json500(ctx, gin.H{
		"errno": 1,
		"error": err.Error(),
	})
}

func Json400(ctx *gin.Context, data interface{}) {
	ctx.JSON(400, data)
}

func Redirect(ctx *gin.Context, url string) {
	ctx.Redirect(http.StatusMovedPermanently, url)
}
