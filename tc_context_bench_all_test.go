package log

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll/Phuslu_OneField/G1-16      	 9060376	       131.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Phuslu_OneField/G2-16      	 9119248	       131.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Zerolog_OneField/G1-16     	 8059534	       149.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Zerolog_OneField/G2-16     	 7811768	       151.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Arenalog_OneField/G1-16    	29849972	        40.29 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Arenalog_OneField/G2-16    	12593814	        94.44 ns/op	      11 B/op	       0 allocs/op

func BenchmarkAll(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_OneField)
}

/*
Arenalog — The most efficient single‑core logger

Arenalog is a high‑performance structured logger engineered for workloads where
logging must remain predictable, low‑latency, and never become a burden inside
constrained environments.

Single‑core timings (G1):
- Arenalog: 40.29 ns/op  (~24.8M logs/sec)
- Phuslu:   131.4 ns/op  (~7.6M logs/sec)
- Zerolog:  149.5 ns/op  (~6.6M logs/sec)

Two‑goroutine timings (G2):
- Arenalog: 94.44 ns/op  (~10.5M logs/sec)
- Phuslu:   131.9 ns/op  (~7.5M logs/sec)
- Zerolog:  151.4 ns/op  (~6.6M logs/sec)

The G1 numbers show that normal applications can rely on Arenalog for all their
logging without dedicating more than one core. The G2 slowdown reflects
Arenalog’s design choice: it prioritizes single‑core determinism over
multi‑core scaling. Multi‑threaded applications can use CPU affinity to pin
Arenalog to a dedicated core and preserve its single‑core characteristics.

Design goals:
- deterministic behavior under load
- minimal branching
- zero allocations in the hot path
- minimal GC interaction
- predictable latency
- optimized for 1‑core execution
- simple ingestion pipeline
- efficient timestamping and field handling

When to choose Arenalog:
- two threads containers where one thread shares the business logic with arenalog
and the other thread can run business logic 100%.
- CPU‑bound or latency‑sensitive systems
- constrained hardware
- environments where multi‑core scaling is irrelevant or undesirable

Examples:
- embedded systems
- edge devices
- proxies and gateways
- WASM runtimes
- unikernels and micro‑VMs
- serverless cold starts
- real‑time telemetry
- HPC nodes
- mobile and game engines

Positioning:
If your application runs on bare resources, no logger comes close in terms of efficiency.
*/
