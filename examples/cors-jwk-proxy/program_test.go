package main

import (
	"context"
	"encoding/json"
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/jwt/jwtest"
	"github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"go.nhat.io/httpmock"
	"net/http/httptest"
	"testing"
)

type headers map[string]string

// for mocking JWT tokens and JWKS signing keys server
var jwksMockedServer = func(s *httpmock.Server) {
	s.ExpectGet("/.well-known/jwks.json").
		ReturnCode(200).
		ReturnJSON(jwtest.PublicJwkKeys())
}

// jwtTokenGenerator
var jwtTokenGenerator = jwtest.JwtGenerator().
	WithSubject("1234user").
	WithAudiences("12loc").
	WithIssuer("http://localhost")

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
			scenario:     "no token",
			method:       "GET",
			path:         "/foo",
			expectedCode: 401,
			expectedJson: JSON{"error": "missing_token"},
		},
		{
			scenario: "bad jwk keys server 500 error",
			method:   "GET",
			path:     "/foo",
			headers:  headers{"Authorization": "Bearer invali.jwt.here"},
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/.well-known/jwks.json").ReturnCode(500)
			}),
			expectedCode: 403, //forbidden
			expectedJson: JSON{"error": "invalid_token"},
		},
		{
			scenario:     "malformed jwt token",
			method:       "GET",
			path:         "/foo",
			headers:      headers{"Authorization": "Bearer invali.jwt.here"},
			mockServer:   httpmock.New(jwksMockedServer),
			expectedCode: 403, //forbidden
			expectedJson: JSON{"error": "invalid_token"},
		},
		{
			scenario:     "invalid jwt token, is no signed by jwk set",
			method:       "GET",
			path:         "/foo",
			headers:      headers{"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			mockServer:   httpmock.New(jwksMockedServer),
			expectedCode: 403, //forbidden
			expectedJson: JSON{"error": "invalid_token"},
		},
		{
			scenario:     "valid jwt token, but is expired",
			method:       "GET",
			path:         "/foo",
			headers:      headers{"Authorization": jwtTokenGenerator.IsExpired().GenerateBearer()},
			mockServer:   httpmock.New(jwksMockedServer),
			expectedCode: 403, //forbidden
			expectedJson: JSON{"error": "expired_token"},
		},
		{
			scenario: "success jwt, but upstream fails",
			method:   "GET",
			path:     "/foo",
			headers:  headers{"Authorization": jwtTokenGenerator.GenerateBearer()},
			mockServer: httpmock.New(jwksMockedServer, func(s *httpmock.Server) {
				s.ExpectGet("/foo").
					ReturnCode(500)
			}),
			expectedCode: 500,
		},
		{
			scenario: "success",
			method:   "GET",
			path:     "/foo",
			headers:  headers{"Authorization": jwtTokenGenerator.GenerateBearer()},
			mockServer: httpmock.New(jwksMockedServer, func(s *httpmock.Server) {
				s.ExpectGet("/foo").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{"hello": "world", "abc": "foo"})
			}),
			expectedCode: 200,
			expectedJson: JSON{"hello": "world", "abc": "foo"},
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
			t.Setenv("JWK_URI", srv.URL()+"/.well-known/jwks.json")

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
