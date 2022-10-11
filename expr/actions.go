package expr

import (
	"context"
)

type Actions []Action

func (a Actions) BuildHandler(ctx context.Context, next Handler) Handler {
	for i := len(a); i > 0; i-- {
		next = a[i-1].BuildHandler(ctx, next)
	}
	return next
}
