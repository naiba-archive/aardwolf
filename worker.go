package aardwolf

import (
	"log"
	"sync/atomic"
)

// Worker worker
type Worker struct {
	pool *Pool
	args chan interface{}
}

func (w *Worker) start() {
	go func() {
		for arg := range w.args {
			atomic.AddUint64(&w.pool.runningNum, 1)
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
			atomic.AddUint64(&w.pool.runningNum, ^uint64(1-1))
			w.pool.luckWorkers.Lock()
			w.pool.idleWorkers = append(w.pool.idleWorkers, w)
			w.pool.luckWorkers.Unlock()
		}
	}()
}
