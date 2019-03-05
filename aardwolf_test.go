package aardwolf

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

const (
	_  = 1 << (10 * iota)
	kb // 1024
	mb // 1048576
	gb // 1073741824
	tb // 1099511627776
	pb // 1125899906842624
	eb // 1152921504606846976
	zb // 1180591620717411303424
	yb // 1208925819614629174706176
)

const (
	poolSize uint64 = 200000
	testTime        = 1000000
)

var curMem uint64

func demoPoolFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func TestNoneRecover(t *testing.T) {
	p := New(2, func(i interface{}) {
		panic(i)
	}, nil)
	for i := 0; i < 10; i++ {
		p.Push(fmt.Sprintf("bingo-no_recover-%d", i))
	}
}

func TestHasRecover(t *testing.T) {
	p := New(2, func(i interface{}) {
		panic(i)
	}, func(e interface{}) {
		t.Log("recover", e)
	})
	for i := 0; i < 10; i++ {
		p.Push(fmt.Sprintf("bingo-recover-%d", i))
	}
}

func BenchmarkSingleFunc(b *testing.B) {
	p := New(poolSize, demoPoolFunc, nil)

	b.StartTimer()
	for j := 0; j < b.N; j++ {
		for i := 0; i < testTime; i++ {
			p.Push(10)
		}
	}
	b.StopTimer()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	tmpMem := (mem.TotalAlloc - curMem) / mb
	curMem = mem.TotalAlloc
	b.Logf("memory usage:%d MB", tmpMem)
}
