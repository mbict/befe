package dsltest

import (
	"context"
	"github.com/mbict/befe/dsl"
	"net/http"
	"net/http/httptest"
)

func RunRequest(program dsl.Action, req *http.Request) *httptest.ResponseRecorder {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	h := program.BuildHandler(ctx, func(_ http.ResponseWriter, _ *http.Request) {})

	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	return rw
}
