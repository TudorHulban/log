package arena

// Test Case 13: Hammer arena with huge messages

// Test: Try to hammer arena with write request larger than arena size.
// Say 90% - 100% of requests are greater than arena size.
// Verifies:
// 1. The 10% of valid writes are correctly written under multiple rotations.
// 2. Cursor works correctly.
