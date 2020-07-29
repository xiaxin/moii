package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}
