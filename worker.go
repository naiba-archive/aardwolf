package aardwolf

import (
	"log"
	"sync/atomic"
	"time"
)

// Worker worker
type Worker struct {
	pool     *Pool
	args     chan interface{}
	lastWork time.Time
}

func (w *Worker) start() {
	go func() {
		for arg := range w.args {
			// panic handler
			defer func() {
				if r := recover(); r != nil {
					if w.pool.Recover != nil {
						w.pool.Recover(r)
					} else {
						log.Println("Aardwolf: panic", r)
					}
					w.release()
				}
			}()
			atomic.AddUint64(&w.pool.runningNum, 1)
			w.lastWork = time.Now()
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
			w.free()
		}
	}()
}

func (w *Worker) release() {
	atomic.AddUint64(&w.pool.workerNum, ^uint64(1-1))
	close(w.args)
	w.args = nil
	w = nil
}

func (w *Worker) free() {
	atomic.AddUint64(&w.pool.runningNum, ^uint64(1-1))
	w.pool.lockWorkers.Lock()
	w.pool.idleWorkers = append(w.pool.idleWorkers, w)
	w.pool.lockWorkers.Unlock()
}
