package conduit

import (
	"context"
	"sync"
)

// FanIn returns a channel that emits values from multiple input channels. The
// order of values is not preserved.
// See [FanOut] for creating multiple input channels.
// To preserve order, use [Bridge] with [ChanChan].
func FanIn[T any](ctx context.Context, streams ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(streams))
	for _, stream := range streams {
		go func(s <-chan T) {
			defer wg.Done()
			for v := range s {
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}(stream)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

// Bridge returns a channel that emits values from a stream of inner channels
// in order.
// See [ChanChan] for creating a channel of channels.
func Bridge[T any](ctx context.Context, chStream <-chan <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			var innerCh <-chan T
			select {
			case maybeInnerCh, ok := <-chStream:
				if !ok {
					return
				}
				innerCh = maybeInnerCh
			case <-ctx.Done():
				return
			}
			for v := range OrDone(ctx, innerCh) {
				select {
				case out <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return out
}

// OrDone returns a channel that emits values from the input stream until it is
// closed or the context is canceled. Useful for preventing goroutine leaks
// when consuming from possibly-blocking or shared channels.
func OrDone[T any](ctx context.Context, stream <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return out
}

// Tee returns two channels that each emit the same values as the input stream.
func Tee[T any](ctx context.Context, stream <-chan T) (_, _ <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		for v := range OrDone(ctx, stream) {
			var out1, out2 = out1, out2
			for range 2 {
				select {
				case <-ctx.Done():
					return
				case out1 <- v:
					out1 = nil
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}
