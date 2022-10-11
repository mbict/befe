package dsl

import (
	"github.com/mbict/befe/expr"
	"net/http"
)

func HttpHandlerCallback(handlerFunc http.HandlerFunc) Expr {
	return expr.ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		handlerFunc(rw, r)
		return true, nil
	})
}

func HandlerCallback(handlerFunc expr.Handler) Expr {
	return expr.ActionFunc(handlerFunc)
}

func MiddlewareHandlerCallback(handlerFunc expr.HandleBuilder) Expr {
	return handlerFunc
}
