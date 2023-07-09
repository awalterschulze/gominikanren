package gomini

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func Go(ctx context.Context, w *sync.WaitGroup, f func()) {
	waitForRoutine(ctx)
	if w != nil {
		w.Add(1)
	}
	go func() {
		if w != nil {
			defer w.Done()
		}
		defer releaseRoutine(ctx)
		f()
	}()
}

func SetMaxRoutines(ctx context.Context, max int) context.Context {
	var count atomic.Int32
	ctx = context.WithValue(ctx, "count", &count)
	limitChan := make(chan struct{}, max)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
				limitChan <- struct{}{}
				count, ok := getCount(ctx)
				if ok {
					if rand.Intn(100) == 0 {
						fmt.Printf("Go Routines: %v\n", count.Load())
					}
				}
			}
		}
	}()
	for i := 0; i < max; i++ {
		limitChan <- struct{}{}
	}
	return context.WithValue(ctx, "limit", limitChan)
}

// Slows down the creation of go routines, by waiting for a routine to finish if the max has been reached.
// It will only wait for 1 second before giving up.
func waitForRoutine(ctx context.Context) bool {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if !ok {
		return true
	}
	select {
	case <-limitChan:
	case <-ctx.Done():
		return false
	}
	count, ok := ctx.Value("count").(*atomic.Int32)
	if ok {
		count.Add(1)
	}
	return true
}

func releaseRoutine(ctx context.Context) bool {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if !ok {
		return true
	}
	select {
	case limitChan <- struct{}{}:
	case <-ctx.Done():
		return false
	}
	count, ok := getCount(ctx)
	if ok {
		count.Add(-1)
	}
	return true
}

func getCount(ctx context.Context) (*atomic.Int32, bool) {
	c, ok := ctx.Value("count").(*atomic.Int32)
	return c, ok
}
