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
					w.release()
				}
			}()

			w.pool.runningNumL.Lock()
			w.pool.runningNum++
			w.pool.runningNumL.Unlock()

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
	w.pool.workerNumL.Lock()
	defer w.pool.workerNumL.Unlock()
	if w.pool.workerNum > 1 {
		w.pool.workerNum--
	}
	close(w.args)
	w.args = nil
}

func (w *Worker) free() {
	w.pool.runningNumL.Lock()
	defer w.pool.runningNumL.Unlock()
	if w.pool.runningNum > 1 {
		w.pool.runningNum--
	}

	w.pool.lockWorkers.Lock()
	defer w.pool.lockWorkers.Unlock()
	w.pool.idleWorkers = append(w.pool.idleWorkers, w)
}
