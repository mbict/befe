package main

import (
	"context"
	"fmt"
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
		mockServer httpmock.Mocker

		expectedCode int
		expectedJson JSON
	}{
		{
			scenario:     "no existing path",
			method:       "GET",
			path:         "/foo",
			expectedCode: 404,
		},
		{
			scenario: "external call to /accounts fails with internal server error",
			method:   "GET",
			path:     "/locations",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/accounts").
					ReturnCode(500)
			}),
			expectedCode: 500,
		},
		{
			scenario: "external call to sub page fails with internal server error",
			method:   "GET",
			path:     "/locations",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/accounts").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "123"},
							JSON{"id": "345"},
							JSON{"id": "678"},
						},
					})
				s.ExpectGet("/bff/product-plans?account_id=123,345,678&offset=0&size=10").
					ReturnCode(500)
			}),
			expectedCode: 500,
		},
		{
			scenario: "successfully call and one page result (no pagination needed)",
			method:   "GET",
			path:     "/locations",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/accounts").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "123"},
							JSON{"id": "345"},
							JSON{"id": "678"},
						},
						"metadata": JSON{"total": 3},
					})

				s.ExpectGet("/bff/product-plans?account_id=123,345,678&offset=0&size=10").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "101", "account_id": "678"},
							JSON{"id": "102", "account_id": "123"},
							JSON{"id": "103", "account_id": "678"},
						},
						"metadata": JSON{"total": 3},
					})
			}),
			expectedCode: 200,
		},
		{
			scenario: "successfully call with 3 sub pages on price_list results (3 pages needed to be paginated)",
			method:   "GET",
			path:     "/locations",
			mockServer: httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/accounts").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "123"},
							JSON{"id": "345"},
							JSON{"id": "678"},
						},
					})

				s.ExpectGet("/bff/product-plans?account_id=123,345,678&offset=0&size=10").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "101", "account_id": "123"},
							JSON{"id": "102", "account_id": "123"},
							JSON{"id": "103", "account_id": "678"},
							JSON{"id": "104", "account_id": "123"},
							JSON{"id": "105", "account_id": "678"},
							JSON{"id": "106", "account_id": "678"},
							JSON{"id": "107", "account_id": "678"},
							JSON{"id": "108", "account_id": "678"},
							JSON{"id": "109", "account_id": "678"},
							JSON{"id": "110", "account_id": "678"},
						},
						"metadata": JSON{"total": 23},
					})

				s.ExpectGet("/bff/product-plans?account_id=123,345,678&offset=10&size=10").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "201", "account_id": "123"},
							JSON{"id": "202", "account_id": "123"},
							JSON{"id": "203", "account_id": "678"},
							JSON{"id": "204", "account_id": "123"},
							JSON{"id": "205", "account_id": "678"},
							JSON{"id": "206", "account_id": "678"},
							JSON{"id": "207", "account_id": "678"},
							JSON{"id": "208", "account_id": "678"},
							JSON{"id": "209", "account_id": "678"},
							JSON{"id": "210", "account_id": "678"},
						},
						"metadata": JSON{"total": 23},
					})

				s.ExpectGet("/bff/product-plans?account_id=123,345,678&offset=20&size=10").
					ReturnHeader("Content-Type", "application/json").
					ReturnJSON(JSON{
						"data": []interface{}{
							JSON{"id": "301", "account_id": "123"},
							JSON{"id": "302", "account_id": "123"},
						},
						"metadata": JSON{"total": 23},
					})
			}),
			expectedCode: 200,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)

			//run mock server if config present
			if tc.mockServer != nil {
				srv := tc.mockServer(t)
				defer srv.Close()
				t.Setenv("API_URI", srv.URL())
			}
			//
			//Program().
			//	BuildHandler(context.Background(), nil).
			//	ServeHTTP(rw, r)

			expr.WrapHttpHandler(
				Program().BuildHandler(context.Background(), nil),
			).ServeHTTP(rw, r)

			fmt.Println("body ->", rw.Body.String())

			assert.Equal(t, tc.expectedCode, rw.Code)

		})
	}
}
