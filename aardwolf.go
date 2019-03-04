package aardwolf

import (
	"sync"
	"time"
)

// Pool aardwolf goruntine pool
type Pool struct {
	Alloc    time.Duration
	Func     func(interface{})
	Recovery func(interface{})

	capNum      int64
	runningNum  int64
	luckCounter sync.Mutex
	idleWorkers []*Worker
	luckWorkers sync.Mutex
}

// New new pool
func New(capNum int64, alloc time.Duration, r, f func(interface{})) *Pool {
	return &Pool{
		capNum:   capNum,
		Alloc:    alloc,
		Func:     f,
		Recovery: r,
	}
}

// Push 向池中添加任务
func (p *Pool) Push(x interface{}) error {
	p.luckCounter.Lock()
	defer p.luckCounter.Unlock()
	if p.capNum <= p.runningNum {
		return ErrNoFreeWorker
	}
	// 取空闲 Worker
	var w *Worker
	p.luckWorkers.Lock()
	if len(p.idleWorkers) > 0 {
		w = p.idleWorkers[len(p.idleWorkers)-1]
		p.idleWorkers = p.idleWorkers[:len(p.idleWorkers)-1]
	} else {
		w = &Worker{
			pool: p,
		}
	}
	p.luckWorkers.Unlock()
	w.start()
	p.runningNum++
	w.args <- x
	return nil
}
