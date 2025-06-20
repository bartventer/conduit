package conduit

import (
	"context"
	"slices"
	"sync"
	"testing"
)

func checkStream(t *testing.T, got, want []int) {
	t.Helper()
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func runCancelledStreamTest(t *testing.T, setup func(ctx context.Context) <-chan int) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := setup(ctx)
	var got []int
	for v := range ch {
		got = append(got, v)
	}
	if len(got) > 0 {
		t.Errorf("expected no values after cancellation, got %v", got)
	}
}

func runStreamTest(t *testing.T, setup func(ctx context.Context) <-chan int, want []int) {
	var got []int
	for v := range setup(t.Context()) {
		got = append(got, v)
	}
	checkStream(t, got, want)
}

func TestTee(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctx context.Context) (<-chan int, <-chan int)
		want  []int
	}{
		{
			name: "Tee",
			setup: func(ctx context.Context) (<-chan int, <-chan int) {
				stream := From(ctx, 1, 2, 3)
				return Tee(ctx, stream)
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" ok", func(t *testing.T) {
			ch1, ch2 := tt.setup(t.Context())
			var got1, got2 []int
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for v := range ch1 {
					got1 = append(got1, v)
				}
			}()
			go func() {
				defer wg.Done()
				for v := range ch2 {
					got2 = append(got2, v)
				}
			}()
			wg.Wait()
			slices.Sort(got1)
			slices.Sort(got2)
			if !slices.Equal(got1, tt.want) || !slices.Equal(got2, tt.want) {
				t.Errorf("got %v and %v, want %v", got1, got2, tt.want)
			}
		})

		t.Run(tt.name+" cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			ch1, ch2 := tt.setup(ctx)
			var got1, got2 []int
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for v := range ch1 {
					got1 = append(got1, v)
				}
			}()
			go func() {
				defer wg.Done()
				for v := range ch2 {
					got2 = append(got2, v)
				}
			}()
			wg.Wait()
			if len(got1) > 0 || len(got2) > 0 {
				t.Errorf("expected no values after cancellation, got %v and %v", got1, got2)
			}
		})
	}
}

func TestFanIn(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctx context.Context) <-chan int
		want  []int
	}{
		{
			name: "FanIn",
			setup: func(ctx context.Context) <-chan int {
				streams := FanOut(ctx, 2, func(c context.Context, index uint) <-chan int {
					return From(ctx, 1, 2, 3)
				})
				return FanIn(ctx, streams...)
			},
			want: []int{1, 1, 2, 2, 3, 3},
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

func TestBridge(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctx context.Context) <-chan int
		want  []int
	}{
		{
			name: "Bridge",
			setup: func(ctx context.Context) <-chan int {
				chStream := ChanChan(ctx, 3, func(ctx context.Context, index uint) int {
					return int(index * 10)
				})
				return Bridge(ctx, chStream)
			},
			want: []int{0, 10, 20},
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
