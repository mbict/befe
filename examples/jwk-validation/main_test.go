package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JwkValidationTestTestSuite struct {
	dsltest.Suite
}

func TestJwkValidationTestSuite(t *testing.T) {
	suite.Run(t, &JwkValidationTestTestSuite{dsltest.NewSuite(Program)})
}

func (suite *JwkValidationTestTestSuite) TestUnauthorized_when_no_token_provided() {
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	res := suite.Request(req)

	suite.Equal(0, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusUnauthorized, res.Code)
	suite.Equal(`no token`, res.Body.String())
}

func (suite *JwkValidationTestTestSuite) TestAuthorized_with_valid_token() {
	httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
		httpmock.NewBytesResponder(200, dsltest.JwkKeySet())) //we return a mocked JWK key set with public keys
	httpmock.RegisterResponder("GET", "http://localhost/",
		httpmock.NewStringResponder(200, `hi`)) //this should be the service response in the reverse proxy

	//generate a request with a valid jwt token, that was signed with a key from the public JWK
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set("Authorization", "Bearer "+dsltest.ValidJwtToken("t.tester@localhost"))
	res := suite.Request(req)

	suite.Equal(2, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.Equal(`hi`, res.Body.String())
}

func (suite *JwkValidationTestTestSuite) TestDenied_with_malformed_token() {
	httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
		httpmock.NewBytesResponder(200, dsltest.JwkKeySet())) //we return a mocked JWK key set with public keys

	//generate a request with a valid jwt token, that was signed with a key from the public JWK
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set("Authorization", "Bearer foo.bar.baz")
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusForbidden, res.Code)
	suite.Equal(`invalid token`, res.Body.String())
}

func (suite *JwkValidationTestTestSuite) TestDenied_with_invalid_signed_token() {
	httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
		httpmock.NewBytesResponder(200, dsltest.JwkKeySet())) //we return a mocked JWK key set with public keys

	//generate a request with a valid jwt token, that was signed with a key from the public JWK
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set("Authorization", "Bearer "+dsltest.InvalidSignedJwtToken("t.tester@localhost"))
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusForbidden, res.Code)
	suite.Equal(`invalid token`, res.Body.String())
}

func (suite *JwkValidationTestTestSuite) TestDenied_with_expired_token() {
	httpmock.RegisterResponder("GET", "http://localhost/.well-known/jwks.json",
		httpmock.NewBytesResponder(200, dsltest.JwkKeySet())) //we return a mocked JWK key set with public keys

	//generate a request with a valid jwt token, that was signed with a key from the public JWK
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set("Authorization", "Bearer "+dsltest.ExpiredJwtToken("t.tester@localhost"))
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusForbidden, res.Code)
	suite.Equal(`expired token`, res.Body.String())
}
