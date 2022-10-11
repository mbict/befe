package main

import (
	"context"
	"encoding/json"
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type headers map[string]string

func TestProgram(t *testing.T) {
	testcases := []struct {
		scenario string
		env      map[string]string

		expectedCode     int
		expectedResponse []byte
		expectedJson     JSON
	}{
		{
			scenario:         "should return printed default env value",
			expectedCode:     200,
			expectedResponse: []byte(`TEST_VAR = Hello`),
		},
		{
			scenario:         "should return printed set env value",
			env:              map[string]string{"TEST_VAR": "foobar"},
			expectedCode:     200,
			expectedResponse: []byte(`TEST_VAR = foobar`),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			//set env to point to the test server
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			expr.WrapHttpHandler(
				Program().BuildHandler(context.Background(), nil),
			).ServeHTTP(rw, r)

			assert.Equal(t, tc.expectedCode, rw.Code)

			if tc.expectedResponse != nil {
				assert.Equal(t, string(tc.expectedResponse), rw.Body.String())
			}

			if tc.expectedJson != nil {
				body := JSON{}

				err := json.Unmarshal(rw.Body.Bytes(), &body)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedJson, body)
			}

		})
	}
}
