package dsl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNot(t *testing.T) {
	h := Not(ConditionHandler(func(r *http.Request) bool {
		return false
	})).BuildCondition()

	assert.True(t, h(nil))
}

func TestNot_match(t *testing.T) {
	h := Not(ConditionHandler(func(r *http.Request) bool {
		return true
	})).BuildCondition()

	assert.False(t, h(nil))
}

func TestOr_firstfail(t *testing.T) {
	h := Or(ConditionHandler(func(r *http.Request) bool {
		return false
	}),
		ConditionHandler(func(r *http.Request) bool {
			return true
		})).BuildCondition()

	assert.True(t, h(nil))
}

func TestOr_2ndfail(t *testing.T) {
	h := Or(ConditionHandler(func(r *http.Request) bool {
		return true
	}),
		ConditionHandler(func(r *http.Request) bool {
			return false
		})).BuildCondition()

	assert.True(t, h(nil))
}

func TestOr_allfail(t *testing.T) {
	h := Or(ConditionHandler(func(r *http.Request) bool {
		return false
	}),
		ConditionHandler(func(r *http.Request) bool {
			return false
		})).BuildCondition()

	assert.False(t, h(nil))
}

func TestOr_allpass(t *testing.T) {
	h := Or(ConditionHandler(func(r *http.Request) bool {
		return true
	}),
		ConditionHandler(func(r *http.Request) bool {
			return true
		})).BuildCondition()

	assert.True(t, h(nil))
}

func TestAnd_firstfail(t *testing.T) {
	h := And(ConditionHandler(func(r *http.Request) bool {
		return false
	}),
		ConditionHandler(func(r *http.Request) bool {
			return true
		})).BuildCondition()

	assert.False(t, h(nil))
}

func TestAnd_2ndfail(t *testing.T) {
	h := And(ConditionHandler(func(r *http.Request) bool {
		return true
	}),
		ConditionHandler(func(r *http.Request) bool {
			return false
		})).BuildCondition()

	assert.False(t, h(nil))
}

func TestAnd_allfail(t *testing.T) {
	h := And(ConditionHandler(func(r *http.Request) bool {
		return false
	}),
		ConditionHandler(func(r *http.Request) bool {
			return false
		})).BuildCondition()

	assert.False(t, h(nil))
}

func TestAnd_allpass(t *testing.T) {
	h := And(ConditionHandler(func(r *http.Request) bool {
		return true
	}),
		ConditionHandler(func(r *http.Request) bool {
			return true
		})).BuildCondition()

	assert.True(t, h(nil))
}
