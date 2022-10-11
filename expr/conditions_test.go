package expr

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionHandlerThen(t *testing.T) {
	var callOrder []string
	ch := NewConditionHandler(MockedCondition(true, func() { callOrder = append(callOrder, "condition") })).
		Then(MockedAction(true, nil, func() { callOrder = append(callOrder, "then") })).
		Else(MockedAction(true, nil, func() { callOrder = append(callOrder, "else") })).
		BuildHandler(context.Background(), MockedHandler(true, nil, func() { callOrder = append(callOrder, "next") }))

	ch(nil, nil)

	assert.Equal(t, []string{"condition", "then", "next"}, callOrder)
}

func TestConditionHandlerElse(t *testing.T) {
	var callOrder []string
	ch := NewConditionHandler(MockedCondition(false, func() { callOrder = append(callOrder, "condition") })).
		Then(MockedAction(true, nil, func() { callOrder = append(callOrder, "then") })).
		Else(MockedAction(true, nil, func() { callOrder = append(callOrder, "else") })).
		BuildHandler(context.Background(), MockedHandler(true, nil, func() { callOrder = append(callOrder, "next") }))

	ch(nil, nil)

	assert.Equal(t, []string{"condition", "else", "next"}, callOrder)
}
