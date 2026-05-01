// Package core is the single canonical middleware engine for v11.
//
// Rules:
//   - retry, timeout, circuit breaker, load shedding, and rate limit logic lives
//     in resilience/*.
//   - logging, metrics, tracing, and correlation logic lives in observability/*.
//   - middleware/core only orders, deduplicates, and adapts wrappers.
//   - transport packages must not create independent chain engines.
package core
