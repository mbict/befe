package main

import (
	"context"
	"github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestProgram(t *testing.T) {
	testcases := []struct {
		scenario string
		method   string
		path     string

		expectedCode     int
		expectedResponse []byte
	}{
		{
			scenario:     "get the index",
			method:       "GET",
			path:         "/",
			expectedCode: 200,
		},
		{
			scenario:     "get the about page",
			method:       "GET",
			path:         "/about",
			expectedCode: 200,
		},
		{
			scenario:     "get test index page",
			method:       "GET",
			path:         "/test",
			expectedCode: 200,
		},
		{
			scenario:     "get test wildcard page",
			method:       "GET",
			path:         "/test/foo/bar",
			expectedCode: 200,
		},
		{
			scenario:         "get test wildcard page",
			method:           "GET",
			path:             "/foo",
			expectedCode:     404,
			expectedResponse: []byte(`the page does not exist 404 html here`),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)

			expr.WrapHttpHandler(
				Program().BuildHandler(context.Background(), nil),
			).ServeHTTP(rw, r)

			if tc.expectedResponse != nil {
				assert.Equal(t, string(tc.expectedResponse), rw.Body.String())
			}

			assert.Equal(t, tc.expectedCode, rw.Code)

		})
	}
}
