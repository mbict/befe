package dsl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//--- ValueFromPatten
func TestValueFromPattern(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/1234?foo=abc-123-xyz", nil)
	p := ValueFromPattern(`abc-(.*)-xyz`, ValueFromQuery("foo"))

	value := p(r)

	assert.Equal(t, "123", value)
}

func TestValueFromPattern_no_match(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/1234?foo=foo", nil)
	p := ValueFromPattern(`abc-(.*)-xyz`, ValueFromQuery("foo"))

	value := p(r)

	assert.Equal(t, "", value)
}

func TestValueFromPattern_with_param_array_first_mactch(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/1234?foo=abc-123-xyz&foo=abc-456-xyz", nil)
	p := ValueFromPattern(`abc-(.*)-xyz`, ValueFromQuery("foo"))

	value := p(r)

	assert.Equal(t, "123", value)
}

func TestValueFromPattern_with_param_array(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/1234?foo=abc&foo=abc-456-xyz", nil)
	p := ValueFromPattern(`abc-(.*)-xyz`, ValueFromQuery("foo"))

	value := p(r)

	assert.Equal(t, "456", value)
}

func TestValueFromPattern_no_match_on_array(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/1234?foo=abc&foo=xyz", nil)
	p := ValueFromPattern(`abc-(.*)-xyz`, ValueFromQuery("foo"))

	value := p(r)

	assert.Equal(t, "", value)
}
