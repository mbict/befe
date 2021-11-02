package dsl

import (
	"context"
	"fmt"
	"net/http"
)

type Actions []Action

func (a Actions) BuildHandler(ctx context.Context, next Handler) Handler {
	for i := len(a); i > 0; i-- {
		next = a[i-1].BuildHandler(ctx, next)
	}
	return next
}

type ActionFunc http.HandlerFunc

func (a ActionFunc) BuildHandler(_ context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) {
		a(rw, r)
		if next != nil {
			next(rw, r)
		}
	}
}

type Action interface {
	BuildHandler(ctx context.Context, next Handler) Handler
}

type ActionBuilder func(ctx context.Context, next Handler) Handler

func (a ActionBuilder) BuildHandler(ctx context.Context, next Handler) Handler {
	return (ActionBuilder(a))(ctx, next)
}

//--- utility helpers

//With wraps the actions in a chain, useful for middleware chaining
func With(actions ...Action) Action {
	return Actions(actions)
}

func Print(message string) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			fmt.Println(message)
			next(rw, r)
		}
	})
}
