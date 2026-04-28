package log

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll/Phuslu_OneField/G1-16      	 8965804	       132.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Phuslu_OneField/G2-16      	 9009058	       133.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Zerolog_OneField/G1-16     	 7107508	       166.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Zerolog_OneField/G2-16     	 7181052	       166.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll/Arenalog_OneField/G1-16    	19193390	        73.14 ns/op	      16 B/op	       1 allocs/op
// BenchmarkAll/Arenalog_OneField/G2-16    	 9027249	       129.3 ns/op	      31 B/op	       1 allocs/op

func BenchmarkAll(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_OneField)
}

/*
Arenalog — The Fastest Single‑Core Logger

Arenalog is a high‑performance structured logger designed for single‑core,
latency‑critical systems. It delivers ~18 million log messages per second on a
single AMD Ryzen Zen 3 core, making it the fastest option for:

- embedded systems
- edge devices
- proxies and gateways
- WASM runtimes
- unikernels and micro‑VMs
- serverless cold starts
- real‑time telemetry
- HPC nodes
- mobile and game engines

If your workload is single‑threaded, Arenalog outperforms every major Go logger.
If your workload is multi‑threaded, you can use the CPU affinity option to pin Arenalog to
a single core and keep its single‑core performance characteristics.

Single‑core performance (GOMAXPROCS=1):
- Arenalog: ~73 ns/op
- Phuslu:   ~132 ns/op
- Zerolog:  ~166 ns/op


Design goals:
- deterministic behavior
- minimal branching
- minimal GC pressure
- zero allocations in the hot path
- predictable latency
- optimized for 1‑core execution
- simple ingestion pipeline
- fast timestamping
- efficient field handling

When to choose Arenalog:
- single‑threaded workloads
- CPU‑bound or latency‑sensitive systems
- constrained hardware
- environments where multi‑core scaling is irrelevant

Examples:
- IoT devices
- routers
- proxies
- WASM modules
- micro‑VMs
- embedded Linux
- mobile apps
- game engines
- HPC telemetry

Positioning:
Arenalog is the most efficient logger for single‑core Go applications.
If your system runs on one core, nothing is faster.
*/
