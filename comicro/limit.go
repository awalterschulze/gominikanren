package comicro

import (
	"context"
	"sync"
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
	limitChan := make(chan struct{}, max)
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
		return true
	case limitChan <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func ReleaseRoutine(ctx context.Context) bool {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if !ok {
		return true
	}
	select {
	case <-time.After(1 * time.Second):
		return true
	case <-limitChan:
		return true
	case <-ctx.Done():
		return false
	}
}
