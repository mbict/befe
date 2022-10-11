package http

import (
	"context"
	. "github.com/mbict/befe/expr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOffsetPaginatedResults_multiple_calls(t *testing.T) {
	callstack := []string{}

	mockedApiCall := NewPromise(func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error) {
		callstack = append(callstack, "conditional")
		return success(rw, r)
	}).OnSuccess(MockedAction(true, nil, func() { callstack = append(callstack, "onSuccess") })).
		OnFailure(MockedAction(true, nil, func() { callstack = append(callstack, "onFailure") }))

	h := OffsetPaginatedResults(
		10,
		func(_ *http.Request) interface{} {
			callstack = append(callstack, "maxResultsFetcher")
			return 22
		}, mockedApiCall).
		BuildHandler(context.Background(), MockedHandler(true, nil, func() {
			callstack = append(callstack, "next")
		}))

	r := httptest.NewRequest("GET", "/foo", nil)
	c, err := h(nil, r)

	assert.True(t, c)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"conditional",
		"onSuccess",
		"maxResultsFetcher",
		"conditional",
		"onSuccess",
		"maxResultsFetcher",
		"conditional",
		"onSuccess",
		"maxResultsFetcher",
		"next",
	}, callstack)

}

func TestOffsetPaginatedResults_calls_with_failure(t *testing.T) {
	callstack := []string{}

	h := OffsetPaginatedResults(
		10,
		func(_ *http.Request) interface{} {
			callstack = append(callstack, "maxResultsFetcher")
			return 22
		},
		NewPromise(func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error) {
			callstack = append(callstack, "conditional")
			return failure(rw, r)
		}).
			OnSuccess(MockedAction(true, nil, func() { callstack = append(callstack, "onSuccess") })).
			OnFailure(MockedAction(true, nil, func() { callstack = append(callstack, "onFailure") }))).
		BuildHandler(context.Background(), MockedHandler(true, nil, func() {
			callstack = append(callstack, "next")
		}))

	r := httptest.NewRequest("GET", "/foo", nil)
	c, err := h(nil, r)

	assert.True(t, c)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"conditional",
		"onFailure",
		"next",
	}, callstack)

}
