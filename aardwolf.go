package aardwolf

import (
	"sync"
	"time"
)

// Pool aardwolf goruntine pool
type Pool struct {
	Func    func(interface{})
	Recover func(interface{})
	Recycle time.Duration

	capNum      uint64
	workerNum   uint64
	runningNum  uint64
	idleWorkers []*Worker
	lockCounter sync.Mutex
	lockWorkers sync.Mutex
}

// New new pool
func New(size uint64, wr time.Duration, f, r func(interface{})) *Pool {
	p := &Pool{
		capNum:  size,
		Recycle: wr,
		Func:    f,
		Recover: r,
	}
	p.recycling()
	return p
}

// Push 向池中添加任务
func (p *Pool) Push(x interface{}) error {
	// 取空闲 Worker
	var w *Worker
	p.lockCounter.Lock()
	defer p.lockCounter.Unlock()
	p.lockWorkers.Lock()
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
		p.lockWorkers.Unlock()
		return ErrNoFreeWorker
	}
	p.lockWorkers.Unlock()
	w.args <- x
	return nil
}

func (p *Pool) recycling() {
	go func() {
		t := time.NewTicker(p.Recycle)
		for {
			select {
			case <-t.C:
				p.lockWorkers.Lock()
				for i := 0; i < len(p.idleWorkers); i++ {
					if p.idleWorkers[i].lastWork.Add(p.Recycle).Before(time.Now()) {
						p.idleWorkers[i].release()
						p.idleWorkers = append(p.idleWorkers[:i], p.idleWorkers[i+1:]...)
					}
				}
				p.lockWorkers.Unlock()
			}
		}
	}()
}
