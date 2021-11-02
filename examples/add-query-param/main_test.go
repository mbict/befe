package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AddParamTestSuite struct {
	dsltest.Suite
}

func TestAddQueryParamSuite(t *testing.T) {
	suite.Run(t, &AddParamTestSuite{dsltest.NewSuite(Program)})
}

func (suite *AddParamTestSuite) TestSingleParamAdded() {
	httpmock.RegisterResponder("GET", "http://localhost/test?foo=bar&foo=baz",
		httpmock.NewStringResponder(200, `hi`))

	req := httptest.NewRequest("GET", "http://localhost/test", nil)
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.Equal(`hi`, res.Body.String())
}
