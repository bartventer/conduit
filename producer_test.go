package conduit

import (
	"context"
	"slices"
	"testing"
)

func TestProducer(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctx context.Context) <-chan int
		want  []int
	}{
		{
			name: "From",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, From(ctx, 1, 2, 3), 3)
			},
			want: []int{1, 2, 3},
		},
		{
			name: "FromSeq",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, FromSeq(ctx, slices.Values([]int{1, 2, 3})), 3)
			},
			want: []int{1, 2, 3},
		},
		{
			name: "FromSeq2",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, FromSeq2(ctx, slices.All([]int{1, 2, 3})), 3)
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Repeat",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, Repeat(ctx, func(ctx context.Context) int { return 42 }), 2)
			},
			want: []int{42, 42},
		},
		{
			name: "ChanChan",
			setup: func(ctx context.Context) <-chan int {
				chStream := ChanChan(ctx, 3, func(ctx context.Context, index uint) int {
					return int(index * 10)
				})
				return Bridge(ctx, chStream)
			},
			want: []int{0, 10, 20},
		},
		{
			name: "FanOut",
			setup: func(ctx context.Context) <-chan int {
				streams := FanOut(ctx, 2, func(c context.Context, index uint) <-chan int {
					return Take(c, Repeat(c, func(ctx context.Context) int { return int(index * 10) }), 3)
				})
				return FanIn(ctx, streams...)
			},
			want: []int{
				0, 0, 0,
				10, 10, 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" ok", func(t *testing.T) {
			t.Parallel()
			runStreamTest(t, tt.setup, tt.want)
		})
		t.Run(tt.name+" cancelled", func(t *testing.T) {
			t.Parallel()
			runCancelledStreamTest(t, tt.setup)
		})
	}
}
