package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ConditionalRouterTestSuite struct {
	dsltest.Suite
}

func TestConditionalRouterSuite(t *testing.T) {
	suite.Run(t, &ConditionalRouterTestSuite{dsltest.NewSuite(Program)})
}

func (suite *ConditionalRouterTestSuite) TestUnauthorizedWhenNoJwtIsProvided() {

	//httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
	//	httpmock.NewStringResponder(200, `{}`))
	//
	//httpmock.RegisterResponder("GET", "http://localhost/", httpmock.NewStringResponder(200, ``))

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	res := suite.Request(req)

	//	suite.Equal( 1 ,httpmock.GetCallCountInfo()["http://localhost/.well-known/jwks.json"])
	//	suite.Equal( 0 ,httpmock.GetCallCountInfo()["http://localhost/"])
	suite.Equal(0, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusUnauthorized, res.Code)
	suite.Equal(``, res.Body.String())
}

func (suite *ConditionalRouterTestSuite) TestNotFoundOnInvalidRoute() {
	//fake the jwk endpoint
	httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
		httpmock.NewBytesResponder(200, dsltest.JwkKeySet()))

	req := httptest.NewRequest("GET", "http://localhost/non-existing-path", nil)
	req.Header.Set("Authorization", "Bearer "+dsltest.ValidJwtToken("t.tester@localhost"))
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusNotFound, res.Code)
	suite.Equal(``, res.Body.String())
}
