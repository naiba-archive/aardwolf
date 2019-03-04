package aardwolf

import (
	"errors"
)

var (
	// ErrNoFreeWorker 没有空闲Worker
	ErrNoFreeWorker = errors.New("Error no free worker")
)
