package web

import (
	"github.com/gin-gonic/gin"
)

const (
	KeyUser = "middleware#user"
)

type User interface {
	GetUid() uint64
	GetUserName() string
	GetToken() string
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
