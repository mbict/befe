package main

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsltest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ModifyPathTestSuite struct {
	dsltest.Suite
}

func TestModifyPathSuite(t *testing.T) {
	suite.Run(t, &ModifyPathTestSuite{dsltest.NewSuite(Program)})
}

func (suite *ModifyPathTestSuite) TestPathChanged() {
	httpmock.RegisterResponder("GET", "http://localhost/foo/bar/baz",
		httpmock.NewStringResponder(200, `hi`))

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	res := suite.Request(req)

	suite.Equal(1, httpmock.GetTotalCallCount())
	suite.Equal(http.StatusOK, res.Code)
	suite.Equal(`hi`, res.Body.String())
}
