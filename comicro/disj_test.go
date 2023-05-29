package comicro

import "testing"

func TestMplus1(t *testing.T) {
	s1 := single(s_x1())
	s2 := single(s_xy_y1())
	s := Mplus(s1, s2)
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
	s1 := suspend(single(s_x1()))
	s2 := cons(s_xy_y1(), suspend(single(nil)))
	s := Mplus(s1, s2)
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
	s1 := NeverO(EmptyState())
	s2 := single(s_x1())
	s := Mplus(s1, s2)
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
