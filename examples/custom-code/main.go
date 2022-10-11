package main

import (
	"context"
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/expr"
	"net/http"
)

func Program() Expr {

	//callback as middleware give you control when and with data the next chain will be executed, or even break the chain
	callbackAsMiddleware := MiddlewareHandlerCallback(func(next expr.Handler) expr.Handler {
		return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
			r = r.WithContext(context.WithValue(r.Context(), "test", "foobar"))

			return next(rw, r)
		}
	})

	//you only can perform and action, that is execute outside the loop.
	//So changing or adding context will not be shared with the next action in the chain
	//You can also use the HandlerCallback that gives you the oppertunity to break the execution chain by returning an error
	callbackAsAction := HttpHandlerCallback(func(rw http.ResponseWriter, r *http.Request) {
		valueFromContext := r.Context().Value("test").(string)

		rw.Write([]byte(`got value from middleware via context : ` + valueFromContext))
	})

	return With(callbackAsMiddleware, callbackAsAction)
}
