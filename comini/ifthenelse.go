package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
)

// IfThenElseO is a goal that evaluates the second goal if the first goal is successful,
// otherwise it evaluates the third goal.
func IfThenElseO(cond, thn, els comicro.Goal) comicro.Goal {
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		conds := comicro.NewStreamForGoal(ctx, cond, s)
		ifThenElseO(ctx, conds, thn, els, s, ss)
	}
}

func ifThenElseO(ctx context.Context, conds comicro.StreamOfStates, thn, els comicro.Goal, s *comicro.State, res comicro.StreamOfStates) {
	headState, ok := conds.ReadNonNil(ctx)
	if !ok {
		els(ctx, s, res)
		return
	}
	heads := comicro.NewStreamForGoal(ctx, thn, headState)
	rests := comicro.NewEmptyStream()
	comicro.Go(ctx, nil, func() {
		defer close(rests)
		comicro.Bind(ctx, conds, thn, rests)
	})
	comicro.Mplus(ctx, heads, rests, res)
}
