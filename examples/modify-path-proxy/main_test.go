package main

import (
	"context"
	"github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"go.nhat.io/httpmock"
	"net/http/httptest"
	"testing"
)

func TestProgram(t *testing.T) {
	testcases := []struct {
		scenario   string
		method     string
		path       string
		mockServer httpmock.Mocker

		expectedCode int
	}{
		{
			scenario: "should change any request to the fixed path",
			method:   "GET",
			path:     "/test",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/foo/bar/baz").ReturnCode(200)
			}),
			expectedCode: 200,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)

			//create empty mock server if no mocks are present
			if tc.mockServer == nil {
				tc.mockServer = httpmock.New()
			}

			srv := tc.mockServer(t)
			defer srv.Close()

			//set env to point to the test server
			t.Setenv("API_URI", srv.URL())

			expr.WrapHttpHandler(
				Program().BuildHandler(context.Background(), nil),
			).ServeHTTP(rw, r)

			assert.Equal(t, tc.expectedCode, rw.Code)
		})
	}
}
