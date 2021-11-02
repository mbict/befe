package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TransformTestSuite struct {
	dsltest.Suite
}

func TestTransformTestSuite(t *testing.T) {
	suite.Run(t, &TransformTestSuite{dsltest.NewSuite(Program)})
}

func (suite *TransformTestSuite) TestResponseIsFiltered() {
	httpmock.RegisterResponder("GET", "http://localhost/filtered",
		httpmock.NewStringResponder(200, `
		{"data":[
			{"foo":"first", "bar": 1, "baz": ["foo", "bar", "baz"]},
			{"foo":"second", "bar": 2, "baz": ["foo"]},
			{"foo":"last", "bar": 3, "baz": []}
		]}`))

	req := httptest.NewRequest("GET", "http://localhost/filtered", nil)
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.JSONEq(`
		{"data":[
			{"foo":"first", "baz": ["foo", "bar", "baz"]},
			{"foo":"second", "baz": ["foo"]},
			{"foo":"last", "baz": []}
		]}`, res.Body.String())
}

func (suite *TransformTestSuite) TestAlternativeResponseIsFilter() {
	httpmock.RegisterResponder("GET", "http://localhost/filtered-alternative",
		httpmock.NewStringResponder(200, `
		{"data":[
			{"foo":"first", "bar": 1, "baz": ["foo", "bar", "baz"]},
			{"foo":"second", "bar": 2, "baz": ["foo"]},
			{"foo":"last", "bar": 3, "baz": []}
		]}`))

	req := httptest.NewRequest("GET", "http://localhost/filtered-alternative", nil)
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.JSONEq(`
		{"data":[
			{"baz": ["foo", "bar", "baz"]},
			{"baz": ["foo"]},
			{"baz": []}
		]}`, res.Body.String())
}

func (suite *TransformTestSuite) TestPassthroughUnModifiedResponse() {
	httpmock.RegisterResponder("GET", "http://localhost/",
		httpmock.NewStringResponder(200, `
		{"data":[
			{"foo":"first", "bar": 1, "baz": ["foo", "bar", "baz"]},
			{"foo":"second", "bar": 2, "baz": ["foo"]},
			{"foo":"last", "bar": 3, "baz": []}
		]}`))

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.JSONEq(`
		{"data":[
			{"foo":"first", "bar": 1, "baz": ["foo", "bar", "baz"]},
			{"foo":"second", "bar": 2, "baz": ["foo"]},
			{"foo":"last", "bar": 3, "baz": []}
		]}`, res.Body.String())
}

func (suite *TransformTestSuite) TestReturns404OnUnknownPath() {
	req := httptest.NewRequest("GET", "http://localhost/not-found", nil)
	res := suite.Request(req)

	suite.Equal(0, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusNotFound, res.Code)
	suite.Equal(``, res.Body.String())
}
