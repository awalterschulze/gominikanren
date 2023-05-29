package comicro

import "testing"

func TestStreamSingle(t *testing.T) {
	s := NewSingletonStream(EmptyState())
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
	s := ConsStream(EmptyState(), func() StreamOfStates { return NewSingletonStream(EmptyState()) })
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
