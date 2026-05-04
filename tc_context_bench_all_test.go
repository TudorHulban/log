package log

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_OneField/Phuslu_OneField/G1-16      	 8853261	       135.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G2-16      	 8776542	       135.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G3-16      	 8791803	       136.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G4-16      	 8733136	       136.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G1-16     	 7888216	       151.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G2-16     	 7875188	       152.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G3-16     	 7851661	       152.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G4-16     	 7870828	       152.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G1-16    	18455356	        64.79 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G2-16    	11042335	       108.0 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G3-16    	10664072	       110.3 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G4-16    	10783002	       109.1 ns/op	       5 B/op	       0 allocs/op

func BenchmarkAll_OneField(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_OneField)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=1-16         	 4706324	       255.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=2-16         	 8641028	       138.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=3-16         	12812220	        94.44 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=4-16         	16818895	        72.68 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=1-16        	 4494318	       267.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=2-16        	 8168961	       145.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=3-16        	12094114	        99.91 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=4-16        	15629527	        77.51 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=1-16       	15442278	        76.86 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=2-16       	16374577	        71.50 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=3-16       	14127042	        88.08 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=4-16       	13617920	        86.94 ns/op	       0 B/op	       0 allocs/op

func BenchmarkAll_SeveralFields(b *testing.B) {
	b.Run("Phuslu_SeveralFields", BenchmarkPhuslu_WithFields_Parallel)
	b.Run("Zerolog_SeveralFields", BenchmarkZerolog_WithFields_Parallel)
	b.Run("Arenalog_SeveralFields", BenchmarkArenalog_MultipleFields_Parallel)
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
Arenalog achieves exceptional single-core performance and remains faster than alternatives even with cross-core ingestion,
at the cost of a fixed inter-core communication overhead.
*/

// NOTE ON INTERPRETING RESULTS
//
// These benchmarks measure a constrained CPU-bound ingestion model:
//   - many concurrent producers (RunParallel)
//   - controlled GOMAXPROCS
//   - in-memory ingestion pipeline
//   - no real I/O (disk, network, stdout sinks are excluded)
//
// This isolates the cost of the ingestion path and cross-goroutine coordination,
// but it does NOT represent full logging cost in production systems.
//
// Key limitations / divergence points vs real-world usage:
//
// 1. Burst behavior is not modeled
//    Real systems exhibit traffic spikes, retries, and uneven load.
//    Under bursts, contention on ingestion structures and cache lines can increase
//    tail latency significantly compared to steady-state benchmarks.
//
// 2. Backpressure effects are not fully represented
//    In production, ingestion slowdown propagates back to producers,
//    potentially coupling unrelated goroutines and affecting request latency.
//
// 3. I/O and sink costs are excluded
//    File writes, stdout, network transport, batching, compression, and formatting
//    typically dominate end-to-end logging cost in real deployments.
//
// 4. GC behavior may differ in production workloads
//    Even when hot paths show 0 allocs/op, surrounding runtime activity,
//    transient allocations, and escape analysis differences can introduce
//    GC pressure and tail latency spikes not visible here.
//
// 5. Cache-coherence and multi-core scaling effects are partially modeled
//    GOMAXPROCS variations and RunParallel introduce contention, but real systems
//    include NUMA effects, OS scheduling noise, and heterogeneous workloads.
//
// 6. Queueing model sensitivity
//    The ingestion layer behavior under saturation (bounded/unbounded queues,
//    drop policies, or blocking behavior) strongly affects tail latency and
//    is not fully exercised by this benchmark suite.
//
// CONCLUSION:
// Results should be interpreted as:
//   "single-core ingestion efficiency under controlled concurrent load"
// rather than:
//   "end-to-end logging performance in production"
//
// This benchmark is useful for comparing ingestion-path efficiency between
// implementations under identical synthetic conditions, but not for absolute
// system-level performance claims.

// INGESTOR DESIGN NOTES
//
// This ingestor implements a high-throughput in-memory event pipeline
// optimized for concurrent producers and low-latency ingestion.
//
// Core design assumptions:
//   - Many goroutines act as producers (Write calls)
//   - A single coordinated ingestion path processes events
//   - Communication cost between cores (cache coherence) is a primary factor
//   - The system prioritizes predictable low-latency ingestion over complex fan-out
//
// Performance characteristics measured in benchmarks reflect:
//
//   1. Ingestion cost under contention
//      Multiple goroutines concurrently calling Write() stress the ingestion
//      structure (queue/ring buffer/atomic coordination), exposing cross-core
//      synchronization overhead.
//
//   2. Cache-coherence behavior
//      Performance is strongly influenced by how frequently shared state
//      (counters, head/tail indexes, buffers) moves between CPU cores.
//
//   3. Allocation-free hot path (when observed)
//      Zero allocations in benchmarks indicate a stable fast path under test
//      conditions, but does not guarantee absence of allocations in all
//      execution paths or configurations.
//
// Important architectural implications:
//
//   - This is an ingestion system, not a full logging pipeline.
//     Formatting, I/O, batching, and persistence are outside its scope.
//
//   - Single-ingestor designs trade off parallel CPU execution for reduced
//     coordination overhead and more predictable latency.
//
//   - Performance gains primarily come from minimizing cross-core contention,
//     not from reducing instruction count alone.
//
//   - Under real-world workloads, behavior depends heavily on burst patterns,
//     queue saturation, and downstream consumer speed.
//
// In summary:
//
//   The ingestor is optimized for low-latency, high-concurrency ingestion
//   where cross-core communication cost dominates CPU execution cost.
