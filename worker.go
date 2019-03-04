package aardwolf

import (
	"log"
	"sync/atomic"
)

// Worker worker
type Worker struct {
	pool    *Pool
	args    chan interface{}
	release chan struct{}
}

func (w *Worker) start() {
	for arg := range w.args {
		atomic.AddInt64(&w.pool.Running, 1)
		if w.pool.Func != nil {
			w.pool.Func(arg)
		} else {
			fn, ok := arg.(func())
			if ok {
				fn()
			} else {
				log.Println("Aardwolf:", "work is invalid")
			}
		}
		atomic.AddInt64(&w.pool.Running, -1)
	}
}
