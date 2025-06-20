# conduit

**conduit** provides utilities for building composable, channel-based data streams and pipelines in Go.

Inspired by patterns from "Concurrency in Go", by Katherine Cox-Buday, conduit offers generic combinators and helpers for constructing, transforming, and consuming streams of data. All functions are context-aware and designed to prevent goroutine leaks.

## Features

- **Create streams:** `From`, `FromSeq`, `FromSeq2`, `Repeat`
- **Transform/filter:** `Map`, `Skip`, `SkipN`, `Take`, `First`
- **Combine/split:** `FanIn`, `FanOut`, `Tee`, `Bridge`, `ChanChan`
- **Safe consumption:** `OrDone`

## Example

```go
package main

import (
    "context"
    "fmt"
    "math/rand/v2"
    "time"

    "github.com/bartventer/conduit"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Generate an infinite stream of random numbers
    randStream := conduit.Repeat(ctx, func(context.Context) int {
        return rand.IntN(100)
    })

    // Filter out even numbers
    odds := conduit.Skip(ctx, randStream, func(_ context.Context, v int) bool {
        return v%2 == 0
    })

    // Map to string with formatting
    formatted := conduit.Map(ctx, odds, func(_ context.Context, v int) string {
        return fmt.Sprintf("odd: %d", v)
    })

    // Duplicate the stream
    out1, out2 := conduit.Tee(ctx, formatted)

    // Fan-in the two streams back together
    merged := conduit.FanIn(ctx, out1, out2)

    // Take the first 10 values from the merged stream
    for v := range conduit.Take(ctx, merged, 10) {
        fmt.Println(v)
    }
}
```

## Documentation

See [GoDoc](https://pkg.go.dev/github.com/bartventer/conduit) for full API documentation and examples.

## Installation

```bash
go get github.com/bartventer/conduit
```

## License

This project is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See the [LICENSE](LICENSE) file for details.

---

*Contributions and feedback welcome!*
