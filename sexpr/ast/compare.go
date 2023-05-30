package ast

import (
	"strings"
)

func (this *SExpr) Compare(that *SExpr) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	if c := this.Pair.Compare(that.Pair); c != 0 {
		return c
	}
	if c := this.Atom.Compare(that.Atom); c != 0 {
		return c
	}
	return 0
}

func (this *Pair) Compare(that *Pair) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	if c := this.Car.Compare(that.Car); c != 0 {
		return c
	}
	if c := this.Cdr.Compare(that.Cdr); c != 0 {
		return c
	}
	return 0
}

func (this *Atom) Compare(that *Atom) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	if c := compareStringPtr(this.Str, that.Str); c != 0 {
		return c
	}
	if c := compareStringPtr(this.Symbol, that.Symbol); c != 0 {
		return c
	}
	if c := compareFloatPtr(this.Float, that.Float); c != 0 {
		return c
	}
	if c := compareIntPtr(this.Int, that.Int); c != 0 {
		return c
	}
	if c := this.Var.Compare(that.Var); c != 0 {
		return c
	}
	return 0
}

func (this *Variable) Compare(that *Variable) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	if c := strings.Compare(this.Name, that.Name); c != 0 {
		return c
	}
	if c := compareUint(this.Index, that.Index); c != 0 {
		return c
	}
	return 0
}

// compareStringPtr returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareStringPtr(this, that *string) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	return strings.Compare(*this, *that)
}

// compareFloatPtr returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareFloatPtr(this, that *float64) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	return compareFloat(*this, *that)
}

// compareFloat returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareFloat(this, that float64) int {
	if this != that {
		if this < that {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

// compareIntPtr returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareIntPtr(this, that *int64) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	return compareInt(*this, *that)
}

// compareInt returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareInt(this, that int64) int {
	if this != that {
		if this < that {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

// compareUint returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func compareUint(this, that uint64) int {
	if this != that {
		if this < that {
			return -1
		} else {
			return 1
		}
	}
	return 0
}
