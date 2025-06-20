package conduit_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"time"

	"github.com/bartventer/conduit"
)

func ExampleFirst() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 10, 20, 30)
	out := conduit.First(ctx, stream)
	for v := range out {
		fmt.Println(v)
	}
	// Output: 10
}

func ExampleSkip() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3, 4, 5)
	out := conduit.Skip(ctx, stream, func(_ context.Context, v int) bool {
		return v%2 == 0
	})
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 3
	// 5
}

func ExampleSkipN() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3, 4, 5)
	out := conduit.SkipN(ctx, stream, 2)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 4
	// 5
}

func ExampleTake() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3, 4, 5)
	out := conduit.Take(ctx, stream, 3)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleMap() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3)
	out := conduit.Map(ctx, stream, func(_ context.Context, v int) string {
		return fmt.Sprintf("num:%d", v)
	})
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// num:1
	// num:2
	// num:3
}

func ExampleFanIn() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streams := conduit.FanOut(ctx, 3, func(c context.Context, index uint) <-chan int {
		return conduit.Take(c, conduit.Repeat(c, func(ctx context.Context) int {
			return int(index * 10)
		}), 3)
	})
	out := conduit.FanIn(ctx, streams...)
	for v := range conduit.Take(ctx, out, 9) {
		fmt.Println(v)
	}
	// Unordered output:
	// 0
	// 0
	// 0
	// 10
	// 10
	// 10
	// 20
	// 20
	// 20
}

func ExampleBridge() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chStream := conduit.ChanChan(ctx, 3, func(_ context.Context, i uint) int {
		return int(i + 1)
	})
	out := conduit.Bridge(ctx, chStream)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleOrDone() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3)
	out := conduit.OrDone(ctx, stream)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleTee() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := conduit.From(ctx, 1, 2, 3)
	out1, out2 := conduit.Tee(ctx, stream)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v1 := range out1 {
			fmt.Printf("out1: %d\n", v1)
		}
	}()

	go func() {
		defer wg.Done()
		for v2 := range out2 {
			fmt.Printf("out2: %d\n", v2)
		}
	}()
	wg.Wait()
	// Unordered output:
	// out1: 1
	// out1: 2
	// out1: 3
	// out2: 1
	// out2: 2
	// out2: 3
}

func ExampleFrom() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := conduit.From(ctx, 7, 8, 9)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 7
	// 8
	// 9
}

func ExampleRepeat() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := conduit.Take(ctx, conduit.Repeat(ctx, func(context.Context) int { return 42 }), 3)
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 42
	// 42
	// 42
}

func ExampleChanChan() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chStream := conduit.ChanChan(ctx, 2, func(_ context.Context, i uint) int {
		return int(i * 100)
	})
	for ch := range chStream {
		for v := range ch {
			fmt.Println(v)
		}
	}
	// Output:
	// 0
	// 100
}

func ExampleFanOut() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streams := conduit.FanOut(ctx, 2, func(_ context.Context, i uint) <-chan int {
		return conduit.From(ctx, int(i), int(i+10))
	})
	for i, ch := range streams {
		fmt.Printf("stream %d:\n", i)
		for v := range ch {
			fmt.Println(v)
		}
	}
	// Output:
	// stream 0:
	// 0
	// 10
	// stream 1:
	// 1
	// 11
}

func ExampleFromSeq() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := conduit.FromSeq(ctx, slices.Values([]int{4, 5, 6}))
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 4
	// 5
	// 6
}

func ExampleFromSeq2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := conduit.FromSeq2(ctx, slices.All([]int{1, 2, 3}))
	for v := range out {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

var randPCG = rand.New(rand.NewPCG(0xCAFEBABEDEADBEEF, 0xDEADBEEFCAFEBABE))

func genInt(_ context.Context) int {
	return randPCG.IntN(100)
}

func Example_pipeline() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Generate an infinite stream of random numbers
	randStream := conduit.Repeat(ctx, genInt)

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

	// Unordered output:
	// odd: 5
	// odd: 13
	// odd: 5
	// odd: 13
	// odd: 49
	// odd: 49
	// odd: 65
	// odd: 65
	// odd: 45
	// odd: 45
}
