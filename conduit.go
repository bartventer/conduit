// Package conduit provides utilities for building composable, channel-based
// data streams and pipelines in Go.
//
// Conduit offers a set of generic combinators and helpers for constructing,
// transforming, and consuming streams of data. The API is inspired by patterns
// from "Concurrency in Go", by Katherine Cox-Buday, and uses Go generics and
// context-based cancellation throughout.
//
// Features include:
//   - Creating streams from values, sequences, or generators ([From],
//     [FromSeq], [FromSeq2], [Repeat])
//   - Transforming and filtering streams ([Map], [Skip], [SkipN], [Take],
//     [First])
//   - Combining and splitting streams ([FanIn], [FanOut], [Tee], [Bridge],
//     [ChanChan])
//   - Safe consumption ([OrDone])
//
// All functions are context-aware and designed to prevent goroutine leaks.
package conduit
