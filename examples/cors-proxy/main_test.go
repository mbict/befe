package main

import (
	"context"
	"encoding/json"
	. "github.com/mbict/befe/dsl"
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
		headers    map[string]string
		mockServer httpmock.Mocker

		expectedCode     int
		expectedResponse []byte
		expectedJson     JSON
	}{
		{
			scenario: "when upstream request fails with internal error",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/foo").ReturnCode(500)
			}),
			expectedCode: 500,
		},
		{
			scenario: "when upstream request fails with not found",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/foo").ReturnCode(404)
			}),
			expectedCode: 404,
		},
		{
			scenario: "when upstream returns a valid json response, we filter the response and only return abc",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/foo").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{"hello": "world", "abc": "foo"})
			}),
			expectedCode: 200,
			expectedJson: JSON{"abc": "foo"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)
			for h, v := range tc.headers {
				r.Header.Add(h, v)
			}

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
