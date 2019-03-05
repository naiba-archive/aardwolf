# Aardwolf

[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/aardwolf)](https://goreportcard.com/report/github.com/naiba/aardwolf)  [![build status](https://travis-ci.com/naiba/aardwolf.svg?branch=master)](https://travis-ci.com/naiba/aardwolf)

:wolf: A high-performance goroutine pool for Go, inspired by ants.

## Benchmark

```shell
# 10k work, 2k worker
goos: darwin
goarch: amd64
pkg: github.com/naiba/aardwolf
BenchmarkSingleFunc-4   	       1	1335181825 ns/op	80490664 B/op	 1109447 allocs/op
--- BENCH: BenchmarkSingleFunc-4
    aardwolf_test.go:72: memory usage:76 MB
BenchmarkMultiFunc-4    	       1	1493482161 ns/op	78909176 B/op	 1100540 allocs/op
--- BENCH: BenchmarkMultiFunc-4
    aardwolf_test.go:92: memory usage:75 MB
PASS
coverage: 92.6% of statements
ok  	github.com/naiba/aardwolf	3.019s
```