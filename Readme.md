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
