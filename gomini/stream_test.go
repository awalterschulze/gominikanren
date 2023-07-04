package gomini

import (
	"context"
	"testing"
)

func TestStreamSingle(t *testing.T) {
	s := NewSingletonStream(context.Background(), NewEmptyState())
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
