package comicro

import (
	"context"
	"testing"
)

func TestMplus1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s1 := single(ctx, s_x1())
	s2 := single(ctx, s_xy_y1())
	ss := NewEmptyStream()
	go func() {
		defer close(ss)
		Mplus(context.Background(), s1, s2, ss)
	}()
	count := 0
	for a := range ss {
		if a == nil {
			t.Fatalf("expected non nil state")
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 states")
	}
}

func TestMplusNeverO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s1 := NewStreamForGoal(ctx, NeverO, EmptyState())
	s2 := single(ctx, s_x1())
	ss := NewEmptyStream()
	go func() {
		defer close(ss)
		Mplus(ctx, s1, s2, ss)
	}()
	count := 0
	for i := 0; i < 100; i++ {
		a, ok := <-ss
		if !ok {
			t.Fatalf("expected never ending stream")
		}
		if a != nil {
			count++
		}
	}
	if count != 1 {
		if _, ok := ss.ReadNonNil(ctx); !ok {
			t.Fatalf("expected 1 non nil state")
		}
	}
}
