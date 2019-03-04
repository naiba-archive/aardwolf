package aardwolf

import (
	"sync"
	"time"
)

// Pool aardwolf goruntine pool
type Pool struct {
	Cap      uint64
	Running  int64
	Alloc    time.Duration
	Func     func(interface{})
	Recovery func(interface{})

	FreeWorkers []Worker
	WorkPool    *sync.Pool
}

// New new pool
func New(cap uint64, alloc time.Duration, f, r func(interface{})) *Pool {
	p := &Pool{
		Cap:      cap,
		Alloc:    alloc,
		Func:     f,
		Recovery: r,
	}
	p.WorkPool = &sync.Pool{
		New: newPool(p),
	}
	return p
}

func newPool(p *Pool) func() interface{} {
	return func() interface{} {
		return &Worker{
			pool: p,
		}
	}
}
