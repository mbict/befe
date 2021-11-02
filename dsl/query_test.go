package dsl

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSetQuery(t *testing.T) {
	r, _ := http.NewRequest("get", "/test", nil)
	a := SetQuery("test", "foo")
	h := a.BuildHandler(context.Background(), func(_ http.ResponseWriter, r *http.Request) {})

	h(nil, r)

	assert.Equal(t, "/test?test=foo", r.URL.String())
}

func TestSetQuery_string_slice(t *testing.T) {
	r, _ := http.NewRequest("get", "/test", nil)
	a := SetQuery("test", []string{"foo", "bar"})
	h := a.BuildHandler(context.Background(), func(_ http.ResponseWriter, r *http.Request) {})

	h(nil, r)

	assert.Equal(t, "/test?test=foo&test=bar", r.URL.String())
}

func TestSetQuery_string_slice_interface(t *testing.T) {
	r, _ := http.NewRequest("get", "/test", nil)
	a := SetQuery("test", []interface{}{"foo", "bar"})
	h := a.BuildHandler(context.Background(), func(_ http.ResponseWriter, r *http.Request) {})

	h(nil, r)

	assert.Equal(t, "/test?test=foo&test=bar", r.URL.String())
}

func TestSetQuery_with_valuer(t *testing.T) {
	r, _ := http.NewRequest("get", "/test", nil)
	a := SetQuery("test", Valuer(func(r *http.Request) interface{} {
		return "foo"
	}))
	h := a.BuildHandler(context.Background(), func(_ http.ResponseWriter, r *http.Request) {})

	h(nil, r)

	assert.Equal(t, "/test?test=foo", r.URL.String())
}
