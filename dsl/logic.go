package dsl

import (
	"context"
	"github.com/mbict/befe/expr"
	"net/http"
)

type breaker struct {
}

func (b breaker) BuildHandler(_ context.Context, _ expr.Handler) expr.Handler {
	return func(_ http.ResponseWriter, _ *http.Request) (bool, error) {
		return false, nil
	}
}

func Stop() expr.Action {
	return &breaker{}
}
