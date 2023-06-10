package comini

import (
	"context"
	"sync"

	"github.com/awalterschulze/gominikanren/comicro"
)

/*
Disjs is a macro that extends disjunction to arbitrary arguments

(define- syntax disj+
(syntax-rules ()

	((_ g) (Zzz g))
	((_ g0 g . . . ) (disj (Zzz g0) (disj+ g . . . )))))
*/
func Disjs(gs ...comicro.Goal) comicro.Goal {
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		wait := sync.WaitGroup{}
		for i := range gs {
			f := func(i int) func() { return func() { gs[i](ctx, s, ss) } }(i)
			comicro.Go(ctx, &wait, f)
		}
		wait.Wait()
	}
}
