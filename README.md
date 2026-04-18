# Arenalogger

## Benchmark Conditions

### Why runtime.GOMAXPROCS(1) and b.SetParallelism(16) Are Required

Benchmarking loggers in Go is subtle because testing.B.RunParallel does not measure the latency of a single log call. By default, it measures aggregate throughput across all CPU cores.

On a machine with 16 logical CPUs, RunParallel will spawn 16 workers and distribute the work across them.
This makes any logger appear 10–16× faster, even though the logger itself did not improve.

To measure true per‑operation latency, not throughput, the benchmark must remove CPU‑level parallelism and exercise the logger under realistic concurrency.

#### runtime.GOMAXPROCS(1)

This forces the Go scheduler to run the benchmark on exactly one logical CPU.

All goroutines created by RunParallel will execute on the same logical CPU.
This eliminates the throughput illusion caused by multiple cores dividing the work.

This setting ensures that the benchmark measures:

- the real cost of a log call in a concurrency scenario
- the real timestamp cost
- the real JSON cost
- the real writer cost
- the real branch and pipeline behavior

In other words, it reveals true latency.

#### b.SetParallelism(16)

This instructs the benchmark to spawn 16 worker goroutines, even though they all run on a single logical CPU.

This is important because it:

- keeps the CPU pipeline hot
- stabilizes branch prediction
- stabilizes timestamp generation
- stabilizes JSON formatting paths
- simulates realistic concurrent logging load

The result is a stable, low‑jitter measurement of the logger’s actual per‑operation cost.

#### Combined Effect

Using both settings:

```go
runtime.GOMAXPROCS(1)
b.SetParallelism(16)
```

produces the only benchmark configuration that:

- removes multi‑core throughput distortion
- preserves realistic concurrency
- reveals true per‑operation latency
- allows fair comparison between loggers

## Resources

```text
https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
```
