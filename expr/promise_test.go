package expr

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPromise_OnFailure(t *testing.T) {
	callstack := []string{}

	c := NewPromise(
		func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error) {
			callstack = append(callstack, "handler")
			return failure(rw, r)
		}).
		OnFailure(MockedAction(true, nil, func() { callstack = append(callstack, "failure") })).
		OnSuccess(MockedAction(true, nil, func() { callstack = append(callstack, "success") }))

	handler := c.BuildHandler(nil, MockedHandler(true, nil, func() { callstack = append(callstack, "next") }))

	cb, err := handler(nil, nil)

	assert.True(t, cb)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"handler",
		"failure",
		"next",
	}, callstack)
}

func TestPromise_OnSuccess(t *testing.T) {
	callstack := []string{}

	c := NewPromise(
		func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error) {
			callstack = append(callstack, "handler")
			success(rw, r)
			return true, nil
		}).
		OnFailure(MockedAction(true, nil, func() { callstack = append(callstack, "failure") })).
		OnSuccess(MockedAction(true, nil, func() { callstack = append(callstack, "success") }))

	handler := c.BuildHandler(nil, MockedHandler(true, nil, func() { callstack = append(callstack, "next") }))

	handler(nil, nil)

	assert.Equal(t, []string{
		"handler",
		"success",
		"next",
	}, callstack)
}
