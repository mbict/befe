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
			scenario: "obtaining token when oauth server is down should result in 500",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				//pre fetch (make sure we do not expose any info)
				s.ExpectPost("/oauth/token").
					WithBody("grant_type=client_credentials").
					ReturnCode(200)

				//actual call
				s.ExpectPost("/oauth/token").
					WithBody("client_id=client_id&client_secret=client_secret&grant_type=client_credentials").
					ReturnCode(500)
			}),
			expectedCode: 500,
			expectedJson: JSON{"code": "internal_error", "error": "internal error, obtaining access token"},
		},

		{
			scenario: "invalid credentials oauth",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				//pre fetch (make sure we do not expose any info)
				s.ExpectPost("/oauth/token").
					WithBody("grant_type=client_credentials").
					ReturnCode(200)

				//actual call
				s.ExpectPost("/oauth/token").
					WithBody("client_id=client_id&client_secret=client_secret&grant_type=client_credentials").
					ReturnCode(401)
			}),
			expectedCode: 500,
			expectedJson: JSON{"code": "internal_error", "error": "internal error, obtaining access token"},
		},

		{
			scenario: "bad token response from api call",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				//pre fetch oauth (make sure we do not expose any info)
				s.ExpectPost("/oauth/token").
					WithBody("grant_type=client_credentials").
					ReturnCode(200)

				//actual oauth call
				s.ExpectPost("/oauth/token").
					WithBody("client_id=client_id&client_secret=client_secret&grant_type=client_credentials").
					ReturnCode(200).
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"access_token":  "foobar",
						"token_type":    "bearer",
						"refresh_token": "goo_bar_baz",
						"expires_in":    86400,
					})

				//api call
				s.ExpectGet("/test").
					WithHeader("Authorization", "Bearer foobar"). //injected token
					ReturnCode(401)
			}),
			expectedCode: 500,
			expectedJson: JSON{"error": "fetch"},
		},

		{
			scenario: "successful response, with injected oauth token",
			method:   "GET",
			path:     "/foo",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				//pre fetch oauth (make sure we do not expose any info)
				s.ExpectPost("/oauth/token").
					WithBody("grant_type=client_credentials").
					ReturnCode(200)

				//actual oauth call
				s.ExpectPost("/oauth/token").
					WithBody("client_id=client_id&client_secret=client_secret&grant_type=client_credentials").
					ReturnCode(200).
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"access_token":  "foobar",
						"token_type":    "bearer",
						"refresh_token": "goo_bar_baz",
						"expires_in":    86400,
					})

				//api call
				s.ExpectGet("/test").
					WithHeader("Authorization", "Bearer foobar"). //injected token
					ReturnCode(200).
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{"foo": "bar"})
			}),
			expectedCode: 200,
			expectedJson: JSON{"foo": "bar"},
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
			t.Setenv("OAUTH2_TOKEN_URL", srv.URL()+"/oauth/token")
			t.Setenv("OAUTH2_CLIENT_ID", "client_id")
			t.Setenv("OAUTH2_CLIENT_SECRET", "client_secret")

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
