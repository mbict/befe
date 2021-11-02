package dsl

import (
	"net/http"
)

func ResponseCode(code int) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(code)
			next(rw, r)
		}
	})
}

func Deny() Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusForbidden)
			next(rw, r)
		}
	})
}

func Unauthorized() Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusUnauthorized)
			next(rw, r)
		}
	})
}

func NotFound() Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			next(rw, r)
		}
	})
}

func InternalServerError() Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
			next(rw, r)
		}
	})
}
