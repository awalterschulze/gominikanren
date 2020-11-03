package example

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/example/peano"
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (len x y)
    (conde
    [(== x '()) (== y 0)]
    [(fresh (a d z)
        (succ z y)
        (== x `(,a . ,d))
        (len d z))]
  ))
*/
func length(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, nil), micro.EqualO(y, peano.Zero)},
		[]micro.Goal{
			micro.CallFresh(func(a *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(z *ast.SExpr) micro.Goal {
						return mini.ConjPlus(
							peano.Succ(z, y),
							micro.EqualO(x, ast.Cons(a, d)),
							length(d, z),
						)
					})
				})
			}),
		},
	)
}

func TestLength(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return length(ast.NewList(peano.One, peano.Zero, peano.One, peano.Zero, peano.One), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := peano.Parsenat(sexprs[0]), 5
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return length(q, peano.Makenat(4))
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	gotS := sexprs[0]
	for i := 0; i < 4; i++ {
		car, cdr := gotS.Car(), gotS.Cdr()
		if !reflect.DeepEqual(car, ast.NewSymbol(fmt.Sprintf("_%d", i))) {
			t.Fatalf("expected variable, but got %v instead", car)
		}
		gotS = cdr
	}
	if gotS != nil {
		t.Fatalf("expected empty list, but got %v instead", gotS)
	}
}

/*
;; split a list x into a list y of len n and a list z
(define (split-at x n y z)
  (conde
    [(len y n)
    (appendo y z x)]
  ))
*/
func splitAt(x, n, y, z *ast.SExpr) micro.Goal {
	return mini.ConjPlus(
		// interestingly, will fail to bind properly if doing length first
		mini.AppendO(y, z, x),
		length(y, n),
	)
}

/*
(define (split-in-half x y z)
  (fresh (l h)
    (len x l)
    (half l h)
    (split-at x h y z)
  ))
*/
func splitInHalf(x, y, z *ast.SExpr) micro.Goal {
	return micro.CallFresh(func(l *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(h *ast.SExpr) micro.Goal {
			return mini.ConjPlus(
				length(x, l),
				peano.Half(l, h),
				splitAt(x, h, y, z),
			)
		})
	})
}

func TestSplitInHalf(t *testing.T) {
	list := peano.ToList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(x *ast.SExpr) micro.Goal {
			return micro.CallFresh(func(y *ast.SExpr) micro.Goal {
				return mini.ConjPlus(
					micro.EqualO(q, ast.NewList(x, y)),
					splitInHalf(list, x, y),
				)
			})
		})
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := sexprs[0].Car(), peano.ToList([]int{1, 2, 3, 4})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

/*
(define (merge a b c)
  (conde
    [(== a '()) (== b c)]
    [(fresh (aa da)
        (== a `(,aa . ,da))
        (conde
            [(== b '()) (== c a)]
            [(fresh (ab db res)
                (== b `(,ab . ,db))
                (conde
                    [(leq aa ab)
                     (== c `(,aa . ,res))
                     (merge da b res)]
                    [(leq ab aa)
                     (== c `(,ab . ,res))
                     (merge a db res)]
            ))]
        ))]
  ))
*/
func merge(a, b, c *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(a, nil), micro.EqualO(b, c)},
		[]micro.Goal{
			micro.CallFresh(func(aa *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(da *ast.SExpr) micro.Goal {
					return micro.ConjunctionO(
						micro.EqualO(a, ast.Cons(aa, da)),
						mini.Conde(
							[]micro.Goal{micro.EqualO(b, nil), micro.EqualO(c, a)},
							[]micro.Goal{
								micro.CallFresh(func(ab *ast.SExpr) micro.Goal {
									return micro.CallFresh(func(db *ast.SExpr) micro.Goal {
										return micro.CallFresh(func(res *ast.SExpr) micro.Goal {
											return micro.ConjunctionO(
												micro.EqualO(b, ast.Cons(ab, db)),
												mini.Conde(
													[]micro.Goal{
														peano.Leq(aa, ab),
														micro.EqualO(c, ast.Cons(aa, res)),
														merge(da, b, res),
													},
													[]micro.Goal{
														peano.Leq(ab, aa),
														micro.EqualO(c, ast.Cons(ab, res)),
														merge(a, db, res),
													},
												),
											)
										})
									})
								}),
							},
						),
					)
				})
			}),
		},
	)
}

func TestMerge(t *testing.T) {
	list1 := peano.ToList([]int{1, 2, 3, 7, 8, 9})
	list2 := peano.ToList([]int{4, 5, 6})
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return merge(list1, list2, q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := sexprs[0], peano.ToList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

/*
(define (merge-sort a b)
  (conde
    [(== a '()) (== b '())]
    [(fresh (x) (== a `(,x)) (== b `(,x)))]
    [(fresh (x y res a1 a2 s1 s2)
        (== a `(,x ,y . ,res))
        (split-in-half a a1 a2)
        ;; candidates for parallelisation
        (merge-sort a1 s1)
        (merge-sort a2 s2)
        (merge s1 s2 b)
    )]
  ))
*/
func mergeSort(a, b *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(a, nil), micro.EqualO(b, nil)},
		[]micro.Goal{
			micro.CallFresh(func(x *ast.SExpr) micro.Goal {
				return mini.ConjPlus(
					micro.EqualO(a, ast.Cons(x, nil)),
					micro.EqualO(b, ast.Cons(x, nil)),
				)
			}),
		},
		[]micro.Goal{
			micro.CallFresh(func(x *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(y *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(res *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(a1 *ast.SExpr) micro.Goal {
							return micro.CallFresh(func(a2 *ast.SExpr) micro.Goal {
								return micro.CallFresh(func(s1 *ast.SExpr) micro.Goal {
									return micro.CallFresh(func(s2 *ast.SExpr) micro.Goal {
										return mini.ConjPlus(
											micro.EqualO(a, ast.Cons(x, ast.Cons(y, res))),
											splitInHalf(a, a1, a2),
											mergeSort(a1, s1),
											mergeSort(a2, s2),
											merge(s1, s2, b),
										)
									})
								})
							})
						})
					})
				})
			}),
		},
	)
}

func TestMergeSort(t *testing.T) {
	list := peano.ToList([]int{6, 5, 1, 2, 7, 3, 8, 4, 9})
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return mergeSort(list, q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := sexprs[0], peano.ToList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
	// Wayyy too slow to run backwards with larger list.
	list = peano.ToList([]int{1})
	sexprs = micro.Run(1, func(q *ast.SExpr) micro.Goal {
		return mergeSort(q, list)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 6, len(sexprs))
	}
	gotS, wantS := sexprs, []*ast.SExpr{
		peano.ToList([]int{1}),
	}
	if !reflect.DeepEqual(gotS, wantS) {
		t.Fatalf("expected %v, but got %v instead", wantS, gotS)
	}
}
