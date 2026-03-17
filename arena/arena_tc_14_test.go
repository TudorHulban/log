package arena

// Test Case 14: Backpressure policy

// Test: Writer stops indefinetly after n writes.
// Verifies logger enters full mode
// with silently drop (common for high-perf logging).
