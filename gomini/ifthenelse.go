package gomini

import (
	"context"
)

// IfThenElseO is a goal that evaluates the second goal if the first goal is successful,
// otherwise it evaluates the third goal.
func IfThenElseO(cond, thn, els Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		conds := NewStreamForGoal(ctx, cond, s)
		ifThenElseO(ctx, conds, thn, els, s, ss)
	}
}

func ifThenElseO(ctx context.Context, conds StreamOfStates, thn, els Goal, s *State, res StreamOfStates) {
	headState, ok := conds.Read(ctx)
	if !ok {
		els(ctx, s, res)
		return
	}
	heads := NewStreamForGoal(ctx, thn, headState)
	rests := NewEmptyStream()
	Go(ctx, nil, func() {
		defer close(rests)
		Bind(ctx, conds, thn, rests)
	})
	Mplus(ctx, heads, rests, res)
}
