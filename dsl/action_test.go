package dsl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestActionsBuildHandler_stack_is_beeing_called_in_order(t *testing.T) {
	result := ""
	actions := Actions{ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result += "1"
	}), ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result += "2"
	}), ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result += "3"
	})}

	handler := actions.BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) {
		result += "4"
	})

	handler(nil, nil)

	assert.Equal(t, "1234", result)
}

func TestWith_stack_is_beeing_called_in_order(t *testing.T) {
	result := ""
	handler := With(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result += "1"
	}), ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result += "2"
	})).BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) {
		result += "3"
	})

	handler(nil, nil)

	assert.Equal(t, "123", result)
}
