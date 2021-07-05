package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Json(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, data)
}

func JsonOK(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "ok"})
}

func JsonData(ctx *gin.Context, data interface{}) {
	Json(ctx, gin.H{
		"data": data,
	})
}

func Json200(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, data)
}

func HTML(ctx *gin.Context, format string, values ...interface{}) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, format, values...)
}

func ValidateFailed(ctx *gin.Context, err interface{}) {
	Json400(ctx, gin.H{
		"errno": 2,
		"error": err,
	})
}

func Json500(ctx *gin.Context, data interface{}) {
	ctx.JSON(500, data)
}

func JsonError(ctx *gin.Context, err error) {
	Json400(ctx, gin.H{
		"errno": 1,
		"error": err.Error(),
	})
}

func Json400(ctx *gin.Context, data interface{}) {
	ctx.JSON(400, data)
}

func Json403(ctx *gin.Context, data interface{}) {
	ctx.JSON(403, data)
}

func JsonNoLogin(ctx *gin.Context) {
	Json403(ctx, gin.H{
		"errno": 1,
		"error": "no login",
	})
}

func Redirect(ctx *gin.Context, url string) {
	ctx.Redirect(http.StatusMovedPermanently, url)
}
