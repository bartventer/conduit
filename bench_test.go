package conduit

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"
)

var randPCG = rand.New(rand.NewPCG(0xCAFEBABEDEADBEEF, 0xDEADBEEFCAFEBABE))

func makeDelays(n int) []time.Duration {
	delays := make([]time.Duration, n)
	for i := range delays {
		delays[i] = time.Duration(randPCG.IntN(100)) * time.Millisecond
	}
	return delays
}

func BenchmarkConcurrency(b *testing.B) {
	collect := func(_ context.Context, ch <-chan uint, count int) []uint {
		results := make([]uint, 0, count)
		for v := range ch {
			results = append(results, v)
			if len(results) == count {
				break
			}
		}
		return results
	}
	const totalCount = 500000
	delays := makeDelays(totalCount)
	processFunc := func(ctx context.Context, index uint) uint {
		select {
		case <-ctx.Done():
			return 0
		case <-time.After(delays[index]):
			return index
		}
	}

	ctx := context.Background()
	b.Run("Fanout-Fanin", func(b *testing.B) {
		for b.Loop() {
			streams := FanOut(ctx, totalCount, func(c context.Context, index uint) <-chan uint {
				return Map(c, From(c, index), func(ctx context.Context, v uint) uint {
					return processFunc(ctx, v)
				})
			})
			results := collect(ctx, FanIn(ctx, streams...), totalCount)
			_ = results
		}
	})

	b.Run("Bridge w/ ChanStream", func(b *testing.B) {
		for b.Loop() {
			chStream := ChanChan(ctx, totalCount, processFunc)
			results := collect(ctx, Bridge(ctx, chStream), totalCount)
			_ = results
		}
	})
}
