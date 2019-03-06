package aardwolf

import (
	"log"
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
				}
			}()
			w.pool.numLocker.Lock()
			w.pool.runningNum++
			w.pool.numLocker.Unlock()
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

func (w *Worker) free() {
	w.pool.numLocker.Lock()
	w.pool.runningNum--
	w.pool.numLocker.Unlock()
	w.pool.workerPool.Put(w)
}
