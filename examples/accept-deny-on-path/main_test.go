package main

import (
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DecisionTestSuite struct {
	dsltest.Suite
}

func TestDecisionSuite(t *testing.T) {
	suite.Run(t, &DecisionTestSuite{dsltest.NewSuite(Program)})
}

func (suite *DecisionTestSuite) TestFooPathAccepted() {
	req := httptest.NewRequest("GET", "http://localhost/foo", nil)
	res := suite.Request(req)

	suite.Equal(http.StatusOK, res.Code)
	suite.Equal(`hey you!`, res.Body.String())
}

func (suite *DecisionTestSuite) TestDefaultPathDenied() {
	req := httptest.NewRequest("GET", "http://localhost", nil)
	res := suite.Request(req)

	suite.Equal(http.StatusForbidden, res.Code)
	suite.Equal(`nope, denied!`, res.Body.String())
}

func (suite *DecisionTestSuite) TestAnyPathDenied() {
	req := httptest.NewRequest("GET", "http://localhost/test", nil)
	res := suite.Request(req)

	suite.Equal(http.StatusForbidden, res.Code)
	suite.Equal(`nope, denied!`, res.Body.String())
}
