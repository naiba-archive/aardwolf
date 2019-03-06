package aardwolf

import (
	"sync"
)

// Pool aardwolf goruntine pool
type Pool struct {
	Func    func(interface{})
	Recover func(interface{})

	capNum     uint64
	runningNum uint64
	numLocker  sync.Mutex
	workerPool *sync.Pool
}

// New Create a new pool
func New(size uint64, f, r func(interface{})) *Pool {
	p := &Pool{
		capNum:  size,
		Func:    f,
		Recover: r,
	}
	p.workerPool = &sync.Pool{
		New: func() interface{} {
			w := &Worker{
				args: make(chan interface{}),
				pool: p,
			}
			w.start()
			return w
		},
	}
	return p
}

// Release pool
func (p *Pool) Release() {
	p.numLocker.Lock()
	defer p.numLocker.Unlock()
	p.capNum = 0
}

// Push work to pool
func (p *Pool) Push(x interface{}) error {
	p.numLocker.Lock()
	if p.runningNum > p.capNum {
		p.numLocker.Unlock()
		return ErrNoFreeWorker
	}
	p.numLocker.Unlock()

	w := p.workerPool.Get().(*Worker)
	w.args <- x
	return nil
}
