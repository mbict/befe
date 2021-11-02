package dsl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConditionsBuildCondition_stack_is_beeing_called_in_order(t *testing.T) {
	result := ""
	conditions := Conditions{ConditionHandler(func(r *http.Request) bool {
		result += "1"
		return true
	}), ConditionHandler(func(r *http.Request) bool {
		result += "2"
		return true
	}), ConditionHandler(func(r *http.Request) bool {
		result += "3"
		return true
	})}

	handler := conditions.BuildCondition()

	ok := handler(nil)

	assert.True(t, ok)
	assert.Equal(t, "123", result)
}

func TestConditionsBuildCondition_first_fail(t *testing.T) {
	result := ""
	conditions := Conditions{ConditionHandler(func(r *http.Request) bool {
		result += "1"
		return false
	}), ConditionHandler(func(r *http.Request) bool {
		result += "2"
		return true
	}), ConditionHandler(func(r *http.Request) bool {
		result += "3"
		return true
	})}

	handler := conditions.BuildCondition()

	ok := handler(nil)

	assert.False(t, ok)
	assert.Equal(t, "1", result)
}

func TestConditionsBuildCondition_last_fail(t *testing.T) {
	result := ""
	conditions := Conditions{ConditionHandler(func(r *http.Request) bool {
		result += "1"
		return true
	}), ConditionHandler(func(r *http.Request) bool {
		result += "2"
		return true
	}), ConditionHandler(func(r *http.Request) bool {
		result += "3"
		return false
	})}

	handler := conditions.BuildCondition()

	ok := handler(nil)

	assert.False(t, ok)
	assert.Equal(t, "123", result)
}
