package micro

/*
(define (disj2 g1 g2)
	(lambda (s)
		(appendinf (g1 s) (g2 s))
	)
)
*/
func DisjointO(gs ...Goal) Goal {
	if len(gs) == 0 {
		return FailureO()
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := DisjointO(gs[1:]...)
	return func(s Substitution) StreamOfSubstitutions {
		g1s := g1(s)
		g2s := g2(s)
		return appendInf(g1s, g2s)
	}
}
