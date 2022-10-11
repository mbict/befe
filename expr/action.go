package expr

import (
	"context"
	"net/http"
)

type Action interface {
	BuildHandler(ctx context.Context, next Handler) Handler
}

type ActionFunc Handler

func (a ActionFunc) BuildHandler(_ context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if c, err := a(rw, r); c == false || err != nil {
			return c, err
		}
		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}
