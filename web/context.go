package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/xiaxin/moii/dig"
)

const (
	KeyDig   = "middleware#di"
	KeyOAuth = "middleware#oauth"
	KeyUser  = "middleware#user"
)

func SetOAuthToken(ctx *gin.Context, token oauth2.TokenInfo) {
	ctx.Set(KeyOAuth, token)
}

func GetOAuthToken(ctx *gin.Context) (oauth2.TokenInfo, bool) {
	ti, _ := ctx.Get(KeyOAuth)
	token, ok := ti.(oauth2.TokenInfo)
	exists := token!= nil && ok
	return token, exists
}



func SetUser(ctx *gin.Context, user User) {
	ctx.Set(KeyUser, user)
}

//  从上下文获取中间件
func GetUser(ctx *gin.Context) (User, bool) {
	diI, _ := ctx.Get(KeyUser)
	di, ok := diI.(User)
	exists := di != nil && ok
	return di, exists
}



//  从上下文获取中间件
func GetDig(ctx *gin.Context) (*dig.Dig, bool) {
	diI, _ := ctx.Get(KeyDig)
	di, ok := diI.(*dig.Dig)
	exists := di != nil && ok
	return di, exists
}
