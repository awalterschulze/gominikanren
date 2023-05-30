package comicro

import (
	"context"
	"testing"
)

func TestStreamSingle(t *testing.T) {
	s := NewSingletonStream(context.Background(), EmptyState())
	if s == nil {
		t.Fatalf("expected non nil stream")
	}
	a, ok := <-s
	if !ok {
		t.Fatalf("expected unclosed stream")
	}
	if a == nil {
		t.Fatalf("expected non nil state")
	}
	a, ok = <-s
	if ok {
		t.Fatalf("expected closed stream")
	}
}

func TestStreamCons(t *testing.T) {
	s := ConsStream(context.Background(), EmptyState(), func() StreamOfStates { return NewSingletonStream(context.Background(), EmptyState()) })
	if s == nil {
		t.Fatalf("expected non nil stream")
	}
	a, ok := <-s
	if !ok {
		t.Fatalf("expected unclosed stream")
	}
	if a == nil {
		t.Fatalf("expected non nil state")
	}
	a, ok = <-s
	if !ok {
		t.Fatalf("expected still unclosed stream")
	}
	if a == nil {
		t.Fatalf("expected still non nil state")
	}
	a, ok = <-s
	if ok {
		t.Fatalf("expected closed stream")
	}
}
