package aardwolf

import (
	"sync"
)

// Pool aardwolf goruntine pool
type Pool struct {
	Func func(interface{})

	capNum      uint64
	workerNum   uint64
	runningNum  uint64
	idleWorkers []*Worker
	luckCounter sync.Mutex
	luckWorkers sync.Mutex
}

// New new pool
func New(capNum uint64, f func(interface{})) *Pool {
	return &Pool{
		capNum: capNum,
		Func:   f,
	}
}

// Push 向池中添加任务
func (p *Pool) Push(x interface{}) error {

	// 取空闲 Worker
	var w *Worker
	p.luckCounter.Lock()
	defer p.luckCounter.Unlock()
	p.luckWorkers.Lock()
	if len(p.idleWorkers) > 0 {
		w = p.idleWorkers[len(p.idleWorkers)-1]
		p.idleWorkers = p.idleWorkers[:len(p.idleWorkers)-1]
	} else if p.workerNum < p.capNum {
		w = &Worker{
			pool: p,
			args: make(chan interface{}),
		}
		p.workerNum++
		w.start()
	} else {
		p.luckWorkers.Unlock()
		return ErrNoFreeWorker
	}
	p.luckWorkers.Unlock()
	w.args <- x
	return nil
}
