package concato

import (
	"strings"

	. "github.com/awalterschulze/gominikanren/gomini"
)

type Node struct {
	Val  *string
	Next *Node
}

func NewNode(v string, n *Node) *Node {
	return &Node{&v, n}
}

func (n *Node) String() string {
	if n == nil {
		return "[]"
	}
	res := "<nil>"
	if n.Val != nil {
		res = *n.Val
	}
	if n.Next != nil {
		next := n.Next.String()
		next = strings.TrimSuffix(strings.TrimPrefix(next, "["), "]")
		res += "," + next
	}
	res = "[" + res + "]"
	return res
}

// ConcatO is a goal that appends two lists into the third list.
// ConcatO(xs, ys, zs) ⇔
// (xs = ∅ ∧ ys = zs) ∨
// (∃ head, xtail, ztail: xs = Prepend(head, xtail) ∧ zs = Prepend(head, ztail) ∧ ConcatO(xtail, ys, ztail))
func ConcatO(xs, ys, zs *Node) Goal {
	return Disj(
		Conj(
			EqualO(xs, nil),
			EqualO(ys, zs),
		),
		Exists(func(head *string) Goal {
			return Exists(func(xtail *Node) Goal {
				return Exists(func(ztail *Node) Goal {
					return Conjs(
						PrependO(head, xtail, xs),
						PrependO(head, ztail, zs),
						ConcatO(xtail, ys, ztail),
					)
				})
			})
		}),
	)
}

func PrependO(head *string, tail *Node, list *Node) Goal {
	return EqualO(list, &Node{head, tail})
}
