# gominikanren

[![GoDoc](https://godoc.org/github.com/awalterschulze/gominikanren?status.svg)](https://godoc.org/github.com/awalterschulze/gominikanren)

[![Build Status](https://travis-ci.org/awalterschulze/gominikanren.svg?branch=master)](https://travis-ci.org/awalterschulze/gominikanren)

gominikarnen is an implementation of [miniKanren](http://minikanren.org/) in Go.

## What is miniKanren

miniKanren is an embedded Domain Specific Language for logic programming.

If you are unfamiliar with miniKanren here is a great introduction video by Bodil Stokke:

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/2e8VFSSNORg/0.jpg)](https://www.youtube.com/watch?v=2e8VFSSNORg)

If you like that, then the book, [The Reasoned Schemer](https://mitpress.mit.edu/books/reasoned-schemer), explains logical programming in miniKanren by example.

If you want to delve even deeper, then this implementation is based on a very readable paper, [µKanren: A Minimal Functional Core for Relational Programming](http://webyrd.net/scheme-2013/papers/HemannMuKanren2013.pdf), that explains the core algorithm of miniKanren.

## Installation

First [install Go](https://golang.org/doc/install)

And then run on the command line

```
$ go get github.com/awalterschulze/gominikanren
```

## Example

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
	sexprs := micro.Reify(ast.NewVariable("q"), states)
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
