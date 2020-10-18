package micro

/*
(define- syntax conj+
(syntax- rules ()
    ((_ g) (Zzz g))
    ((_ g0 g . . . ) (conj (Zzz g0) (conj+ g . . . )))))
*/

func ConjPlus(gs ...Goal) Goal {
	if len(gs) == 0 {
		return SuccessO()
	}
	if len(gs) == 1 {
		return Zzz(gs[0])
	}
	g1 := Zzz(gs[0])
	g2 := ConjPlus(gs[1:]...)
	return func(s *State) StreamOfStates {
		g1s := g1(s)
		return Bind(g1s, g2)
	}
}

/*
(define- syntax disj+
(syntax-rules ()
    ((_ g) (Zzz g))
    ((_ g0 g . . . ) (disj (Zzz g0) (disj+ g . . . )))))
*/
func DisjPlus(gs ...Goal) Goal {
	if len(gs) == 0 {
		return FailureO()
	}
	if len(gs) == 1 {
		return Zzz(gs[0])
	}
	g1 := Zzz(gs[0])
	g2 := DisjPlus(gs[1:]...)
	return func(s *State) StreamOfStates {
		g1s := g1(s)
		g2s := g2(s)
		return mplus(g1s, g2s)
	}
}

/*
(define- syntax conde
(syntax- rules ()
    ((_ (g0 g . . . ) . . . ) (disj+ (conj+ g0 g . . . ) . . . ))))
*/
func Conde(gs ...[]Goal) Goal {
    conj := make([]Goal, len(gs))
    for i, v := range gs {
        conj[i] = ConjPlus(v...)
    }
    return DisjPlus(conj...)
}

// TODO keep only one, most likely Zzz2. 
// will be a heavier change to make though
/*
(define- syntax Zzz
(syntax- rules ()
    ((_ g) (λg (s/c) (λ$ () (g s/c))))))
*/
func Zzz(g Goal) Goal {
	return func(s *State) StreamOfStates {
		return Suspension(func() StreamOfStates {
            return g(s)
        })
    }
}

// Zzz2 is zzz 'fixed': it also delays function call itself
func Zzz2(g func()Goal) Goal {
	return func(s *State) StreamOfStates {
		return Suspension(func() StreamOfStates {
            return g()(s)
        })
    }
}
