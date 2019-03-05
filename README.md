# Aardwolf

[![Go Report Card](https://goreportcard.com/badge/github.com/naiba/aardwolf)](https://goreportcard.com/report/github.com/naiba/aardwolf)  [![build status](https://travis-ci.com/naiba/aardwolf.svg?branch=master)](https://travis-ci.com/naiba/aardwolf)

:wolf: A high-performance goroutine pool for Go, inspired by ants.

## Benchmark

```shell
# 10k work, 2k worker
goos: darwin
goarch: amd64
pkg: github.com/naiba/aardwolf
BenchmarkSingleFunc-4   	       1	1292621237 ns/op	79728136 B/op	 1101245 allocs/op
--- BENCH: BenchmarkSingleFunc-4
    aardwolf_test.go:72: memory usage:76 MB
BenchmarkMultiFunc-4    	       1	1348455223 ns/op	79673160 B/op	 1099530 allocs/op
--- BENCH: BenchmarkMultiFunc-4
    aardwolf_test.go:92: memory usage:76 MB
PASS
coverage: 92.2% of statements
ok  	github.com/naiba/aardwolf	2.845s
```