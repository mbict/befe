package dsl

import (
	"context"
	"encoding/json"
	. "github.com/mbict/befe/expr"
	"net/http"
	"time"
)

func ResponseCode(code int) Action {
	return ActionFunc(func(rw http.ResponseWriter, req *http.Request) (bool, error) {
		rw.WriteHeader(code)
		return true, nil
	})
}

func Ok() Action {
	return ResponseCode(http.StatusOK)
}

func Created() Action {
	return ResponseCode(http.StatusCreated)
}

func Deny() Action {
	return ResponseCode(http.StatusForbidden)
}

func Unauthorized() Action {
	return ResponseCode(http.StatusUnauthorized)
}

func NotFound() Action {
	return ResponseCode(http.StatusNotFound)
}

func InternalServerError() Action {
	return ResponseCode(http.StatusInternalServerError)
}

func Delay(duration time.Duration) Action {
	return ActionFunc(func(_ http.ResponseWriter, _ *http.Request) (bool, error) {
		time.Sleep(duration)
		return true, nil
	})
}

func Redirect(url string, status int) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		http.RedirectHandler(url, status).ServeHTTP(rw, r)
		return false, nil
	})
}

func TemporaryRedirect(url string) Action {
	return Redirect(url, http.StatusTemporaryRedirect)
}

func PermanentRedirect(url string, status int) Action {
	return Redirect(url, http.StatusPermanentRedirect)
}

// write response headers
type JSON = map[string]interface{}

type responseHeader struct {
	name  string
	value string
}

func (rh *responseHeader) BuildHandler(ctx context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		rw.Header().Set(rh.name, rh.value)

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func SetHeader(name string, value string) Action {
	return &responseHeader{
		name:  name,
		value: value,
	}
}

func JsonContentType() Action {
	return SetHeader("Content-Type", "application/json; charset=utf-8")
}

// write response body
type responseBody []byte

func (rb responseBody) BuildHandler(ctx context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if _, err := rw.Write(rb); err != nil {
			return false, err
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func WriteJson(data interface{}) Action {
	var writeAction Action
	if v, ok := data.(Valuer); ok {
		writeAction = ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
			body, err := json.Marshal(v(r))
			if err != nil {
				return false, err
			}

			if _, err = rw.Write(body); err != nil {
				return false, err
			}

			return true, nil
		})
	} else {
		body, _ := json.Marshal(data)
		writeAction = responseBody(body)
	}

	return Actions{
		JsonContentType(),
		writeAction,
	}
}

func WriteResponse(data []byte) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		rw.Write(data)
		return true, nil
	})
}
