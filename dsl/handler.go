package dsl

import (
	"context"
	"net/http"
)

type Handler = http.HandlerFunc

type HandleBuilder func(next Handler) Handler

var emptyHandler = func(_ http.ResponseWriter, _ *http.Request) {}

func (b HandleBuilder) BuildHandler(ctx context.Context, next Handler) Handler {
	return b(next)
}

//
//func wrapActionHandlers(next Handler, actions ...Action) Handler {
//	for i := len(actions); i > 0; i-- {
//		next = actions[i-1].BuildHandler(next)
//	}
//	return next
//}
//
//func wrapConditionHandlers(next Handler, conditions ...Condition) Handler {
//	for i := len(conditions); i > 0; i-- {
//		next = conditions[i-1].BuildHandler(next)
//	}
//	return next
//}
