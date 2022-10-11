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
		mockServer httpmock.Mocker

		expectedCode int
		expectedJson JSON
	}{
		{
			scenario: "success should return the first active customer from the response",
			method:   "GET",
			path:     "/first-active-customer",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/customers").
					ReturnHeader("Content-Type", "application/json").
					ReturnCode(http.StatusOK).ReturnJSON(JSON{
					"data": []interface{}{
						JSON{"id": "1", "customer_number": "customer-1", "firstname": "Joe", "lastname": "Tester", "active": false, "secret": "super-secret-do-not-share"},
						JSON{"id": "2", "customer_number": "customer-2", "firstname": "Tinus", "lastname": "Tester", "active": true, "secret": "super-secret-do-not-share"},
					},
				})
			}),
			expectedCode: http.StatusOK, //200
			expectedJson: JSON{"id": "2", "customer_number": "customer-2", "firstname": "Tinus", "lastname": "Tester"},
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

			if tc.expectedJson != nil {
				body := JSON{}

				err := json.Unmarshal(rw.Body.Bytes(), &body)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedJson, body)
			}

		})
	}
}
