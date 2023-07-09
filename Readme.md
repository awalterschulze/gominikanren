# gominikanren

[![GoDoc](https://godoc.org/github.com/awalterschulze/gominikanren?status.svg)](https://godoc.org/github.com/awalterschulze/gominikanren)

[![Build Status](https://github.com/awalterschulze/gominikanren/workflows/Go/badge.svg)](https://github.com/awalterschulze/gominikanren/actions)

gominikarnen is an implementation of [miniKanren](http://minikanren.org/) in Go.
miniKanren is an embedded Domain Specific Language for logic programming.

## Implementations

Currently there are 3 implementations included in this repository:

1. gomini - the concurrent version that works on any Go data type that is a pointer
2. micro - the minimal implementation based on the [paper](http://webyrd.net/scheme-2013/papers/HemannMuKanren2013.pdf) (only works on ast.SExpr)
3. mini - extends micro with operators from the book [The Reasoned Schemer](https://mitpress.mit.edu/9780262535519/the-reasoned-schemer/) (only works on ast.SExpr)

## Installation

First [install Go](https://golang.org/doc/install)

And then run on the command line

```
$ go get github.com/awalterschulze/gominikanren
```

## Example: ConcatO using gomini

ConcatO is a goal that concats the first two input arguments into the third input argument.
In this example we use ConcatO to get all the combinations that can produce the linked list `[a, b, c, d]`.

```go
package main

import (
    "context"

    . "github.com/awalterschulze/gominikanren/gomini"
    . "github.com/awalterschulze/gominikanren/gomini/concato"
)

type Pair struct {
    Left  *Node
    Right *Node
}

func main() {
    pairs := toStrings(RunTake(context.Background(), -1, NewState(), func(q *Pair) Goal {
        return ExistO(func(x *Node) Goal {
            return ExistO(func(y *Node) Goal {
                return ConjO(
                    EqualO(&Pair{x, y}, q),
                    ConcatO(
                        x,
                        y,
                        NewNode("a", NewNode("b", NewNode("c", NewNode("d", nil)))),
                    ),
                )
            })
        })
    }))
    fmt.Println(pairs)
}

// Output:
// [
//   {[], [a,b,c,d]},
//   {[a,b,c,d], []},
//   {[a,b,c], [d]},
//   {[a,b], [c,d]},
//   {[a], [b,c,d]},
// ]
```

## Example: Translating Math to miniKanren

gominikanren includes the four logic operators:

 - $=$ => EqualO
 - $\exists$ => ExistO
 - $\land$ (and) => ConjO
 - $\lor$ (or) => DisjO

The next two operators are implicit most programming languages:

 - $\forall$ => a function parameter
 - $\in$ => providing a type

ConcatO was created using the mathematical formula:

$$
\forall\ \Phi\ \Psi\ \Omega, \\
\oplus\ \Phi\ \Psi\ \Omega \equiv \\
    (\Phi = \emptyset \land \Omega = \Psi) \lor \\
    (\exists\ \alpha\ \phi\ \omega,\
    \Phi = \alpha \dblcolon \phi
    \land \Omega = \alpha \dblcolon \omega
    \land \oplus\ \phi\ \Psi\ \omega
    )
$$

This looks rather intimidating, but we can give these variables some better names:

$$
\begin{align*}
&\ \forall\ \ xs\ ys\ zs &&\in        [string],\ \\
&\oplus\ xs\ ys\ zs      &&\equiv (xs = [] \land zs = ys) \\
&                        && \ \lor (\exists \ head \in string, \\
&                        && \ \ \ \ \ \ \ \exists \ xtail\ ztail \in [string], \\
&                        && \ \ \ \ \ \ \ \ \ \ \ \ \ xs = [head, xtail \ldots] \\
&                        && \ \ \ \ \ \ \ \ \ \land zs = [head, ztail \ldots] \\
&                        && \ \ \ \ \ \ \ \ \ \land \oplus\ xtail\ ys\ ztail \ )
\end{align*}
$$

Then we can use this to translate it to gominikanren:

```go
type Node struct {
    Value  *string
    Next   *Node
}

func ConcatO(xs, ys, zs *Node) Goal {
    return DisjO(
        ConjO(
            EqualO(xs, nil),
            EqualO(ys, zs),
        ),
        ExistO(func(head *string) Goal {
            return ExistO(func(xtail *Node) Goal {
                return ExistO(func(ztail *Node) Goal {
                    return ConjO(
                        EqualO(xs, &Node{head, xtail}),
                        EqualO(zs, &Node{head, ztail}),
                        ConcatO(xtail, ys, ztail),
                    )
                })
            })
        }),
    )
}
```


## Example: Original AppendO using micro, mini and ast.SExpr

AppendO is a goal that appends the first two input arguments into the third input argument.
In this example we use AppendO to get all the combinations that can produce the list `(cake & ice d t)`.

```go
package main

import (
    "github.com/awalterschulze/gominikanren/sexpr/ast"
    "github.com/awalterschulze/gominikanren/micro"
    "github.com/awalterschulze/gominikanren/mini"
)

func main() {
    states := micro.RunGoal(
        -1,
        micro.CallFresh(func(x *ast.SExpr) micro.Goal {
            return micro.CallFresh(func(y *ast.SExpr) micro.Goal {
                return micro.ConjunctionO(
                    // (== ,q (cons ,x ,y))
                    micro.EqualO(
                        ast.Cons(x, ast.Cons(y, nil)),
                        ast.NewVariable("q"),
                    ),
                    // (appendo ,x ,y (cake & ice d t))
                    mini.AppendO(
                        x,
                        y,
                        ast.NewList(
                            ast.NewSymbol("cake"),
                            ast.NewSymbol("&"),
                            ast.NewSymbol("ice"),
                            ast.NewSymbol("d"),
                            ast.NewSymbol("t"),
                        ),
                    ),
                )
            })
        }),
    )
    sexprs := micro.Reify("q", states)
    fmt.Println(ast.NewList(sexprs...).String())
}
//Output:
//(
//  (() (cake & ice d t))
//  ((cake) (& ice d t))
//  ((cake &) (ice d t))
//  ((cake & ice) (d t))
//  ((cake & ice d) (t))
//  ((cake & ice d t) ())
//)
```

## Learn More about miniKanren

If you are unfamiliar with miniKanren here is a great introduction video by Bodil Stokke:

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/2e8VFSSNORg/0.jpg)](https://www.youtube.com/watch?v=2e8VFSSNORg)

If you like that, then the book, [The Reasoned Schemer](https://mitpress.mit.edu/9780262535519/the-reasoned-schemer/), explains logical programming in miniKanren by example.

If you want to delve even deeper, then this implementation is based on a very readable paper, [µKanren: A Minimal Functional Core for Relational Programming](http://webyrd.net/scheme-2013/papers/HemannMuKanren2013.pdf), that explains the core algorithm of miniKanren.
