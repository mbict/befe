package expr

import (
	"context"
	"google.golang.org/appengine/log"
	"net/http"
)

type Handler func(rw http.ResponseWriter, r *http.Request) (bool, error)

type HandleBuilder func(next Handler) Handler

//var emptyHandler = Handler(func(_ http.ResponseWriter, _ *http.Request) (bool, error) {
//	return true, nil
//})

// AnyHandler makes sure that the provided handler is not nil
//func AnyHandler(handler Handler) Handler {
//	if handler == nil {
//		return func(_ http.ResponseWriter, _ *http.Request) (bool, error) {
//			return true, nil
//		}
//	}
//	return handler
//}

func (b HandleBuilder) BuildHandler(ctx context.Context, next Handler) Handler {
	return b(next)
}

func WrapHttpHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if _, err := h(rw, r); err != nil {
			log.Errorf(r.Context(), err.Error(), err)
		}
	})
}
