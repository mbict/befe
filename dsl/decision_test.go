package dsl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDecisionTree_BuildHandler(t *testing.T) {
	result := []string{}
	tree := Decisions()
	tree.When(ConditionHandler(func(r *http.Request) bool {
		result = append(result, "cond 1")
		return true
	})).Then(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "act 1")
	}))

	handler := tree.BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "next")
	})

	handler(nil, nil)

	assert.Equal(t, []string{"cond 1", "act 1", "next"}, result)
}

func TestDecisionTree_BuildHandler_DefaultPath(t *testing.T) {
	result := []string{}
	tree := Decisions()
	tree.When(ConditionHandler(func(r *http.Request) bool {
		result = append(result, "cond 1")
		return false
	})).Then(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "act 1")
	}))
	tree.Default(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "default")
	}))

	handler := tree.BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "next")
	})

	handler(nil, nil)

	assert.Equal(t, []string{"cond 1", "default", "next"}, result)
}

func TestDecisionTree_BuildHandler_2nd_Condition_Match(t *testing.T) {
	result := []string{}
	tree := Decisions()
	tree.When(ConditionHandler(func(r *http.Request) bool {
		result = append(result, "cond 1")
		return false
	})).Then(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "act 1")
	}))
	tree.When(ConditionHandler(func(r *http.Request) bool {
		result = append(result, "cond 2")
		return true
	})).Then(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "act 2")
	}))
	tree.Default(ActionFunc(func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "default")
	}))

	handler := tree.BuildHandler(nil, func(_ http.ResponseWriter, r *http.Request) {
		result = append(result, "next")
	})

	handler(nil, nil)

	assert.Equal(t, []string{"cond 1", "cond 2", "act 2", "next"}, result)
}
