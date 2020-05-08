package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaxin/moii/log"
	"go.uber.org/dig"
)

const (
	KeyDi = "middleware#di"
)

func NewDig(container *dig.Container, ctx *gin.Context) *Dig {
	return &Dig{
		container: container,
		logger:    nil,
		ctx:       ctx,
	}
}

// 中间件
func NewDigMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dig := NewDig(dig.New(), ctx)

		dig.Provide(func() *gin.Context {
			return ctx
		})

		ctx.Set(KeyDi, dig)
		ctx.Next()
	}
}

//  从上下文获取中间件
func GetDig(ctx *gin.Context) (*Dig, bool) {
	diI, _ := ctx.Get(KeyDi)
	di, ok := diI.(*Dig)
	exists := di != nil && ok
	return di, exists
}

type Dig struct {
	container *dig.Container
	logger    *log.Logger
	ctx       *gin.Context
}

func (d *Dig) Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return d.error(d.container.Provide(constructor, opts...))
}

func (d *Dig) Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return d.error(d.container.Invoke(function, opts...))
}

func (d *Dig) SetLogger(logger *log.Logger) {
	d.logger = logger
}

func (d *Dig) error(err error) error {
	if nil != err && nil != d.logger {
		d.logger.Error(err)
	}
	return err
}
