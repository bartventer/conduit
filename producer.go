package conduit

import (
	"context"
	"iter"
)

// From returns a channel that emits the provided values.
func From[T any](ctx context.Context, values ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// FromSeq returns a channel that emits values from the provided sequence.
func FromSeq[V any](ctx context.Context, seq iter.Seq[V]) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for v := range seq {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// FromSeq2 returns a channel that emits values from the key-value pairs in the
// provided sequence.
func FromSeq2[K, V any](ctx context.Context, seq iter.Seq2[K, V]) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for _, v := range seq {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// Repeat returns a channel that emits values by repeatedly calling fn(ctx)
// until the context is canceled. Useful for generating infinite streams.
func Repeat[T any](ctx context.Context, fn func(context.Context) T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			val := fn(ctx)
			select {
			case <-ctx.Done():
				return
			case out <- val:
			}
		}
	}()
	return out
}

// ChanChan returns a channel that emits n inner channels, each receiving a
// value from fn(ctx, index) where index ranges from 0 to n-1.
// See [Bridge] for consuming a channel of channels.
func ChanChan[T any](ctx context.Context, n uint, fn func(ctx context.Context, index uint) T) <-chan <-chan T {
	out := make(chan (<-chan T), n)
	go func() {
		defer close(out)
		for i := range n {
			stream := make(chan T, 1)
			go func(index uint) {
				defer close(stream)
				val := fn(ctx, index)
				select {
				case <-ctx.Done():
					return
				case stream <- val:
				}
			}(i)
			select {
			case out <- stream:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// FanOut returns a slice of channels where each channel is produced by
// calling fn(ctx, index) for each index from 0 to n-1.
// See [FanIn] for consuming a slice of channels.
func FanOut[T any](ctx context.Context, n uint, fn func(c context.Context, index uint) <-chan T) []<-chan T {
	streams := make([]<-chan T, n)
	for i := range n {
		streams[i] = fn(ctx, i)
	}
	return streams
}
