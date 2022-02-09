package dsl

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseBody []byte

func (rb responseBody) BuildHandler(ctx context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(rb)

		next(rw, r)
	}
}

func WriteResponseBody(body []byte) Action {
	return responseBody(body)
}

func WriteJson(data interface{}) Action {
	body, _ := json.Marshal(data)
	return Actions{
		SetHeader("Content-Type", "application/json"),
		responseBody(body),
	}
}

type responseHeader struct {
	name  string
	value string
}

func (rh *responseHeader) BuildHandler(ctx context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set(rh.name, rh.value)
		next(rw, r)
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
