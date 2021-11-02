package dsltest

import (
	"github.com/jarcoal/httpmock"
	"github.com/mbict/befe/dsl"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

//Suite is a preconfigured suite for testify suite that setups and teardown httpmock and knows how to load the
// program and requests it
type Suite struct {
	suite.Suite
	program func() dsl.Action
}

//Request loads the program and runs the request
func (suite *Suite) Request(req *http.Request) *httptest.ResponseRecorder {
	return RunRequest(suite.program(), req)
}

func (suite *Suite) SetupSuite() {
	httpmock.Activate()
}

func (suite *Suite) SetupTest() {
	httpmock.Reset()
}

func (suite *Suite) TearDownSuite() {
	httpmock.Deactivate()
}

func NewSuite(program func() dsl.Action) Suite {
	return Suite{
		program: program,
	}
}
