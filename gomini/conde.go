package gomini

// Conde is a disjunction of conjunctions
func Conde(gs ...[]Goal) Goal {
	conj := make([]Goal, len(gs))
	for i, v := range gs {
		conj[i] = Conjs(v...)
	}
	return Disjs(conj...)
}
