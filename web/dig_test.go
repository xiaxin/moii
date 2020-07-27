package web

import (
	"github.com/xiaxin/moii/log"
	"testing"
)

type DigService struct {
}

func (ds *DigService) Say() {
	log.Info("hello world")
}

func TestDig(t *testing.T) {

	dig := NewDig(nil)

	dig.Provide(func() *DigService {
		return &DigService{}
	})
	dig.Invoke(func(ds *DigService) {
		ds.Say()
	})
}
