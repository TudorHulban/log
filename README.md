# LogInfo

Do not use. Race conditions currently.

Simple logger.  
Levels:  
0 - none,  
1 - info,  
2 - warn,  
3 - debug.

## How to use

See external testing file.

## Profiling

```sh
go test -bench=. -run=^$ . -cpuprofile profile.out
go test -bench=Benchmark_Info_Logger -run=^$ . -cpuprofile profile.out
go test -bench=Benchmark_Local_Print_Logger -run=^$ . -cpuprofile profile.out
go test -bench=. -benchmem -cpuprofile profile.out
go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out

go tool pprof profile.out
# with option top, web or pdf
```

## Roadmap (Writer)

### a. New log file at day start

### b. New log file when file reaches x Mb

### c. Old log file is zipped

### d. satisfy Fiber interface below

```go
// WithLogger is a logger interface that output logs with a message and key-value pairs.
type WithLogger interface {
	Tracew(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
	Panicw(msg string, keysAndValues ...any)
}
```

## Resources

```html
https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern

https://blog.mike.norgate.xyz/unlocking-go-slice-performance-navigating-sync-pool-for-enhanced-efficiency-7cb63b0b453e

https://unskilled.blog/posts/lets-dive-a-tour-of-sync.pool-internals/

https://medium.com/@felipedutratine/profile-your-benchmark-with-pprof-fb7070ee1a94
```
