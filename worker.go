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
	go func() {
		for arg := range w.args {
			atomic.AddInt64(&w.pool.runningNum, 1)
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
			w.pool.luckCounter.Lock()
			w.pool.runningNum--
			w.pool.luckWorkers.Lock()
			w.pool.idleWorkers = append(w.pool.idleWorkers, w)
			w.pool.luckWorkers.Unlock()
			w.pool.luckCounter.Unlock()
		}
	}()
}
