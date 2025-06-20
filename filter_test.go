package conduit

import (
	"context"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(ctx context.Context) <-chan int
		want       []int
		skipCancel bool
	}{
		{
			name: "First",
			setup: func(ctx context.Context) <-chan int {
				return First(ctx, From(ctx, 1, 2, 3))
			},
			want: []int{1},
		},
		{
			name: "Skip",
			setup: func(ctx context.Context) <-chan int {
				stream := Skip(ctx, From(ctx, 1, 2, 3), func(ctx context.Context, v int) bool {
					return v%2 == 0 // Skip even numbers
				})
				return OrDone(ctx, stream)
			},
			want: []int{1, 3},
		},
		{
			name: "Skip no eq",
			setup: func(ctx context.Context) <-chan int {
				stream := Skip(ctx, From(ctx, 1, 2, 3), nil)
				return OrDone(ctx, stream)
			},
			want:       []int{1, 2, 3},
			skipCancel: true,
		},
		{
			name: "SkipN",
			setup: func(ctx context.Context) <-chan int {
				stream := SkipN(ctx, From(ctx, 1, 2, 3, 4, 5), 2)
				return OrDone(ctx, stream)
			},
			want: []int{3, 4, 5},
		},
		{
			name: "Take",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, From(ctx, 1, 2, 3, 4, 5), 3)
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Take 1",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, From(ctx, 1, 2, 3, 4, 5), 1)
			},
			want: []int{1},
		},
		{
			name: "Take 0",
			setup: func(ctx context.Context) <-chan int {
				return Take(ctx, From(ctx, 1, 2, 3, 4, 5), 0)
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" ok", func(t *testing.T) {
			t.Parallel()
			runStreamTest(t, tt.setup, tt.want)
		})

		t.Run(tt.name+" cancelled", func(t *testing.T) {
			t.Parallel()
			if tt.skipCancel {
				t.Skip("Skipping cancellation test for this case")
			}
			runCancelledStreamTest(t, tt.setup)
		})
	}
}
