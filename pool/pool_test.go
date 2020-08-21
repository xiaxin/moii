package pool

import (
	"github.com/xiaxin/moii/log"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	p, _ := New(5)

	for i:=0; i<100 ;i++  {
		p.Process(&Task{
			Handler: func(v ...interface{}) {
				d := v[0]

				if d == 5 {
					panic("test")
				}

				log.Info(d)
				time.Sleep(1 * time.Second)
			},
			Params:  []interface{}{i},
		})
	}

	p.WaitAndClose()
}


func TestRoutine(t *testing.T) {
	var (
		wg sync.WaitGroup
	)

	c := make(chan int)

	go func() {
		defer func() {
			wg.Done()
		}()
		v := <- c

		t.Log(v)
		wg.Add(1)
	}()

	c <- 10

	wg.Wait()
}
