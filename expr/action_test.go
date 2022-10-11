package expr

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestActionsBuildHandler_stack_is_beeing_called_in_order(t *testing.T) {
	result := ""
	actions := Actions{ActionFunc(func(_ http.ResponseWriter, r *http.Request) (bool, error) {
		result += "1"
		return true, nil
	}), ActionFunc(func(_ http.ResponseWriter, r *http.Request) (bool, error) {
		result += "2"
		return true, nil
	}), ActionFunc(func(_ http.ResponseWriter, r *http.Request) (bool, error) {
		result += "3"
		return true, nil
	})}

	handler := actions.BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) (bool, error) {
		result += "4"
		return true, nil
	})

	handler(nil, nil)

	assert.Equal(t, "1234", result)
}
