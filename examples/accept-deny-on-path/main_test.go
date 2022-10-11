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
			scenario:         "path foo is accessible",
			method:           "GET",
			path:             "/foo",
			expectedCode:     200,
			expectedResponse: []byte(`hey you!`),
		},
		{
			scenario:         "any other path should be denied",
			method:           "GET",
			path:             "/test",
			expectedCode:     403,
			expectedResponse: []byte(`nope, denied!`),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)

			expr.WrapHttpHandler(
				Program().BuildHandler(context.Background(), nil),
			).ServeHTTP(rw, r)

			assert.Equal(t, tc.expectedCode, rw.Code)

			if tc.expectedResponse != nil {
				assert.Equal(t, string(tc.expectedResponse), rw.Body.String())
			}
		})
	}
}
