package concato

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
)

type Node struct {
	Val  *string
	Next *Node
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	res := "<nil>"
	if n.Val != nil {
		res = *n.Val
	}
	if n.Next != nil {
		res += n.Next.String()
	}
	return res
}

func Prepend(val *string, list *Node) *Node {
	return &Node{val, list}
}

// Concat is a goal that appends two lists into the third list.
func ConcatO(xs, ys, out *Node) comicro.Goal {
	return comicro.Disj(
		comicro.Conj(
			comicro.EqualO(xs, nil),
			comicro.EqualO(ys, out),
		),
		comicro.Exists(func(xhead *string) comicro.Goal {
			return comicro.Exists(func(xtail *Node) comicro.Goal {
				return comicro.Exists(func(res *Node) comicro.Goal {
					return comini.Conjs(
						comicro.EqualO(xs, Prepend(xhead, xtail)),
						comicro.EqualO(Prepend(xhead, res), out),
						ConcatO(xtail, ys, res),
					)
				})
			})
		}),
	)
}

func PrependO(xhead *string, xtail *Node, xs *Node) comicro.Goal {
	return comicro.EqualO(Prepend(xhead, xtail), xs)
}
