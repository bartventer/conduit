package conduit

import (
	"context"
	"testing"
)

func TestTransformer(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctx context.Context) <-chan int
		want  []int
	}{
		{
			name: "Map",
			setup: func(ctx context.Context) <-chan int {
				stream := Map(ctx, From(ctx, 1, 2, 3), func(ctx context.Context, v int) int {
					return v * 10 // Multiply each value by 10
				})
				return Take(ctx, stream, 3)
			},
			want: []int{10, 20, 30},
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
