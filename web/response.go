package web

import (
	"github.com/gin-gonic/gin"
)

func ValidateFailed(ctx *gin.Context, err interface{}) {
	ctx.JSON(500, gin.H{
		"errno": 1,
		"error": err,
	})
}

func Json200(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, data)
}
