package main

import (
	"context"
	"encoding/json"
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"go.nhat.io/httpmock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type headers map[string]string

func TestProgram(t *testing.T) {
	testcases := []struct {
		scenario   string
		method     string
		path       string
		cookies    []*http.Cookie
		headers    map[string]string
		mockServer httpmock.Mocker

		expectedCode     int
		expectedResponse []byte
		expectedJson     JSON
		expectedHeaders  headers
	}{
		{
			scenario: "main should go to frontend",
			method:   "GET",
			path:     "/",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/frontend/").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario: "login with cookie, should go to backend check session",
			method:   "GET",
			path:     "/login",
			cookies: []*http.Cookie{
				{
					Name:  "ssid",
					Value: "test",
					Path:  "/",
				},
			},
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/backend/check-session").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario: "post to login should go to backend login",
			method:   "POST",
			path:     "/login",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectPost("/backend/login").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario: "post to login should go to backend login",
			method:   "POST",
			path:     "/register",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectPost("/backend/register").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario: "get to login should go to frontend",
			method:   "GET",
			path:     "/register",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/frontend/register").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario: "get to consent should go to backend",
			method:   "GET",
			path:     "/consent",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/backend/consent").ReturnCode(200)
			}),
			expectedCode: 200,
		},
		{
			scenario:        "not-here endpoint should redirect to login",
			method:          "GET",
			path:            "/not-here",
			mockServer:      httpmock.New(func(s *httpmock.Server) {}),
			expectedHeaders: headers{"Location": "/login"},
			expectedCode:    307,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)

			for _, cookie := range tc.cookies {
				r.AddCookie(cookie)
			}

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
			t.Setenv("BACKEND_URI", srv.URL()+"/backend")
			t.Setenv("FRONTEND_URI", srv.URL()+"/frontend")

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

			for k, v := range tc.expectedHeaders {
				assert.Equal(t, v, rw.Header().Get(k))
			}

		})
	}
}
