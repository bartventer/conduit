package conduit

import (
	"context"
)

// First returns a channel that emits the first value from the input stream.
// If the stream is empty or the context is canceled, the channel will be closed
// without emitting any values.
func First[T any](ctx context.Context, stream <-chan T) <-chan T {
	out := make(chan T, 1)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return
		case v, ok := <-stream:
			if !ok {
				return
			}
			out <- v
		}
	}()
	return out
}

// Skip returns a channel that emits values from the input stream for which
// eq(ctx, v) returns false. If eq is nil, the input stream is returned
// unchanged.
func Skip[T any](ctx context.Context, stream <-chan T, eq func(context.Context, T) bool) <-chan T {
	if eq == nil {
		return stream
	}
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range stream {
			if eq(ctx, v) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// SkipN returns a channel that skips the first n values from the input stream,
// then emits the rest. If n is 0, the input stream is returned unchanged.
func SkipN[T any](ctx context.Context, stream <-chan T, n uint) <-chan T {
	if n == 0 {
		return stream
	}
	out := make(chan T)
	go func() {
		defer close(out)
		for range n {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-stream:
				if !ok {
					return
				}
			}
		}
		for v := range stream {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// Take returns a channel that emits the first n values from the input stream,
// or fewer if the stream closes or the context is canceled.
func Take[T any](ctx context.Context, stream <-chan T, n uint) <-chan T {
	switch n {
	case 0:
		// Similar convention to calling time.After(0) which returns a closed channel.
		out := make(chan T)
		close(out)
		return out
	case 1:
		return First(ctx, stream)
	}
	out := make(chan T)
	go func() {
		defer close(out)
		for range n {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
