package conduit

import (
	"context"
)

// Map returns a channel that emits the results of applying fn to each value
// from the input stream.
func Map[T, U any](ctx context.Context, stream <-chan T, fn func(context.Context, T) U) <-chan U {
	out := make(chan U)
	go func() {
		defer close(out)
		for v := range stream {
			val := fn(ctx, v)
			select {
			case <-ctx.Done():
				return
			case out <- val:
			}
		}
	}()
	return out
}
