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
	s := Mplus(context.Background(), s1, s2)
	count := 0
	for a := range s {
		if a == nil {
			t.Fatalf("expected non nil state")
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 states")
	}
}

func TestMplusSuspend(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s1 := suspend(ctx, single(ctx, s_x1()))
	s2 := cons(ctx, s_xy_y1(), suspend(ctx, single(ctx, nil)))
	s := Mplus(context.Background(), s1, s2)
	count := 0
	for a := range s {
		if a != nil {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 states")
	}
}

func TestMplusNeverO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s1 := NeverO(ctx, EmptyState())
	s2 := single(ctx, s_x1())
	s := Mplus(ctx, s1, s2)
	count := 0
	for i := 0; i < 5; i++ {
		a, ok := <-s
		if !ok {
			t.Fatalf("expected never ending stream")
		}
		if a != nil {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 non nil state")
	}
}
