package dig

import (
	"github.com/xiaxin/moii/log"
	"go.uber.org/dig"
)

type Dig struct {
	container *dig.Container
	logger    *log.Logger
}

func New(logger *log.Logger) *Dig {
	return &Dig{
		container: dig.New(),
		logger: logger,
	}
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

