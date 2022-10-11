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
			scenario:     "no token no access",
			method:       "GET",
			path:         "/foo",
			expectedCode: 401,
			expectedJson: JSON{"code": "access_token_required", "error": "access token required"},
		},
		{
			scenario: "successful request with the token injected in the reversed proxy",
			method:   "GET",
			path:     "/foo",
			headers:  headers{"Authorization": jwtTokenGenerator.GenerateBearer()},
			mockServer: httpmock.New(jwksMockedServer, func(s *httpmock.Server) {
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

				//proxy call
				s.ExpectGet("/foo").
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
			t.Setenv("JWK_URI", srv.URL()+"/.well-known/jwks.json")
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
