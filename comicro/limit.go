package comicro

import "context"

func SetMaxRoutines(ctx context.Context, max int) context.Context {
	limitChan := make(chan struct{}, max)
	return context.WithValue(ctx, "limit", limitChan)
}

func WaitForRoutine(ctx context.Context) {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if ok {
		limitChan <- struct{}{}
	}
}

func ReleaseRoutine(ctx context.Context) {
	limitChan, ok := ctx.Value("limit").(chan struct{})
	if ok {
		<-limitChan
	}
}
