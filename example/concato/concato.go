package concato

import (
	"strings"

	"github.com/awalterschulze/gominikanren/comicro"
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

func Prepend(val *string, list *Node) *Node {
	return &Node{val, list}
}

// ConcatO is a goal that appends two lists into the third list.
// ConcatO(xs, ys, zs) ⇔
// (xs = ∅ ∧ ys = zs) ∨
// (∃ head, xtail, ztail: xs = Prepend(head, xtail) ∧ zs = Prepend(head, ztail) ∧ ConcatO(xtail, ys, ztail))
func ConcatO(xs, ys, zs *Node) comicro.Goal {
	return comicro.Disj(
		comicro.Conj(
			comicro.EqualO(xs, nil),
			comicro.EqualO(ys, zs),
		),
		comicro.Exists(func(head *string) comicro.Goal {
			return comicro.Exists(func(xtail *Node) comicro.Goal {
				return comicro.Exists(func(ztail *Node) comicro.Goal {
					return comicro.Conjs(
						comicro.EqualO(xs, Prepend(head, xtail)),
						comicro.EqualO(Prepend(head, ztail), zs),
						ConcatO(xtail, ys, ztail),
					)
				})
			})
		}),
	)
}

func PrependO(xhead *string, xtail *Node, xs *Node) comicro.Goal {
	return comicro.EqualO(Prepend(xhead, xtail), xs)
}
