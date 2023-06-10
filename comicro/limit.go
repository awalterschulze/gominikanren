package comicro

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func Go(ctx context.Context, w *sync.WaitGroup, f func()) {
	WaitForRoutine(ctx)
	if w != nil {
		w.Add(1)
	}
	go func() {
		if w != nil {
			defer w.Done()
		}
		defer ReleaseRoutine(ctx)
		f()
	}()
}

func SetMaxRoutines(ctx context.Context, max int) context.Context {
	var count atomic.Int32
	ctx = context.WithValue(ctx, "count", &count)
	limitChan := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		limitChan <- struct{}{}
	}
	return context.WithValue(ctx, "limit", limitChan)
}

// Slows down the creation of go routines, by waiting for a routine to finish if the max has been reached.
// It will only wait for 1 second before giving up.
func WaitForRoutine(ctx context.Context) bool {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if !ok {
		return true
	}
	select {
	case <-time.After(1 * time.Second):
	case <-limitChan:
	case <-ctx.Done():
		return false
	}
	MaybePrintRoutineCount(ctx)
	count, ok := ctx.Value("count").(*atomic.Int32)
	if ok {
		count.Add(1)
	}
	return true
}

func ReleaseRoutine(ctx context.Context) bool {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if !ok {
		return true
	}
	select {
	case <-time.After(1 * time.Second):
	case limitChan <- struct{}{}:
	case <-ctx.Done():
		return false
	}
	MaybePrintRoutineCount(ctx)
	count, ok := ctx.Value("count").(*atomic.Int32)
	if ok {
		count.Add(-1)
	}
	return true
}

func MaybePrintRoutineCount(ctx context.Context) {
	if rand.Intn(100) == 0 {
		count, ok := ctx.Value("count").(*atomic.Int32)
		if ok {
			fmt.Printf("Go Routines: %v\n", count.Load())
		}
	}
}
