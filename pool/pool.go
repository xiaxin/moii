package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"github.com/xiaxin/moii/log"
)

var (
	// ErrInvalidPoolCap return if pool size <= 0
	ErrInvalidPoolCap = errors.New("invalid pool cap")
	// ErrPoolAlreadyClosed put task but pool already closed
	ErrPoolAlreadyClosed = errors.New("pool already closed")
)

const (
	// 状态: 运行
	StateRunning = 1
	// 状态: 停止
	StateStoped = 0
)

// Task task to-do
type Task struct {
	Handler func(v ...interface{})
	// 参数
	Params  []interface{}
}

// Pool task pool
type Pool struct {
	//  容量
	capacity       uint64
	//  运行数量
	runningWorkers uint64
	//  状态
	state          int64
	//  任务链
	taskChan       chan *Task
	PanicHandler   func(interface{})
	sync.Mutex
	wg sync.WaitGroup
}

// 新建
func New(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		state:    StateRunning,
		// 多1倍空间，用于缓存
		taskChan:  make(chan *Task, capacity * 2),
	}, nil
}

// 获取 容量
func (p *Pool) GetCapacity() uint64 {
	return p.capacity
}

// 获取 运行数量
func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

// 增加运行数量
func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
	p.wg.Add(1)
}

// 减少运行数量
func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
	p.wg.Done()
}

//  添加任务到运行池
func (p *Pool) Process(task *Task) error {

	//  判断状态
	if p.getState() == StateStoped {
		return ErrPoolAlreadyClosed
	}

	// 判断是否需要增加一个Worker
	p.Lock()
	if p.GetRunningWorkers() < p.GetCapacity() {

		p.incRunning()
		go func() {
			defer func() {
				p.decRunning()
				if r := recover(); r != nil {
					if p.PanicHandler != nil {
						p.PanicHandler(r)
					} else {
						log.Infof("Worker panic: %s\n", r)
					}
				}
			}()

			for {
				select {
				case task, ok := <-p.taskChan:
					if !ok {
						return
					}
					task.Handler(task.Params...)
				}
			}
		}()
	}
	p.Unlock()

	// 发送一个任务
	p.Lock()
	if p.state == StateRunning {
		p.taskChan <- task
	}
	p.Unlock()

	return nil
}

func (p *Pool) getState() int64 {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Pool) setState(state int64) {
	p.Lock()
	defer p.Unlock()

	p.state = state
}

// close safe
func (p *Pool) close() {
	p.Lock()
	defer p.Unlock()

	close(p.taskChan)
}

func (p *Pool) Close() {

	if p.getState() == StateStoped {
		return
	}

	p.setState(StateStoped) // stop put task

	for len(p.taskChan) > 0 { // wait all task be consumed
		time.Sleep(1e6) // reduce CPU load
	}

	p.close()
}

/**
	任务一加入列队，等待处理完关闭协程池。结束任务。
 */
func (p *Pool) WaitAndClose() {
	if p.getState() == StateStoped {
		return
	}

	p.setState(StateStoped) // stop put task

	for len(p.taskChan) > 0 { // wait all task be consumed
		time.Sleep(1e6) // reduce CPU load
	}

	p.close()

	p.wg.Wait()
}
