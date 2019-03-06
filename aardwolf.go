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
	countLocker sync.Mutex

	runningNum uint64

	idleWorkers []*Worker
	lockWorkers sync.Mutex
}

// New Create a new pool
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

// Release pool
func (p *Pool) Release() {
	p.countLocker.Lock()
	p.capNum = 0
	p.workerNum = 0
	p.countLocker.Unlock()
	p.lockWorkers.Lock()
	for i := 0; i < len(p.idleWorkers); i++ {
		p.idleWorkers[i].release(false)
	}
	p.idleWorkers = nil
	p.lockWorkers.Unlock()
	p.Recover = nil
	p.Func = nil
}

// Push work to pool
func (p *Pool) Push(x interface{}) error {
	p.countLocker.Lock()
	p.lockWorkers.Lock()
	// 取空闲 Worker
	var w *Worker
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
		p.countLocker.Unlock()
		return ErrNoFreeWorker
	}
	p.lockWorkers.Unlock()
	p.countLocker.Unlock()
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
						p.idleWorkers[i].release(true)
						p.idleWorkers = append(p.idleWorkers[:i], p.idleWorkers[i+1:]...)
					}
				}
				p.lockWorkers.Unlock()
			}
		}
	}()
}
