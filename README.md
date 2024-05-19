# LogInfo

Simple logger.  
Levels:  
0 - none,  
1 - info,  
2 - warn,  
3 - debug.

## How to use

See external testing file.

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
```
