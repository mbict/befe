package expr

import (
	"context"
	"net/http"
)

type BuildConditionFunc func(c context.Context) ConditionFunc

func (c BuildConditionFunc) BuildCondition(ctx context.Context) ConditionFunc {
	return c(ctx)
}

type ConditionFunc func(r *http.Request) bool

func (c ConditionFunc) BuildCondition(_ context.Context) ConditionFunc {
	return c
}

type Condition interface {
	BuildCondition(context.Context) ConditionFunc
}

type ConditionHandler interface {
	Action
	Condition
	ElseCondition

	Then(...Action) ElseCondition
}

type ElseCondition interface {
	Action

	Else(...Action) Action
}

type conditionHandler struct {
	Condition

	actions     Actions
	elseActions Actions
}

// ConditionHandler executes code in the next stack when a acondition is met.
func NewConditionHandler(condition Condition) ConditionHandler {
	return &conditionHandler{
		Condition:   condition,
		actions:     Actions{},
		elseActions: Actions{},
	}
}

func (c *conditionHandler) BuildHandler(ctx context.Context, next Handler) Handler {
	cond := c.Condition.BuildCondition(ctx)
	handler := c.actions.BuildHandler(ctx, next)         //<- on purpose that next is chained
	elseHandler := c.elseActions.BuildHandler(ctx, next) //<- on purpose that next is chained

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if cond(r) == true {
			return handler(rw, r)
		}
		return elseHandler(rw, r)
	}
}

func (c *conditionHandler) Else(action ...Action) Action {
	c.elseActions = append(c.elseActions, action...)
	return c
}

func (c *conditionHandler) Then(action ...Action) ElseCondition {
	c.actions = append(c.actions, action...)
	return c
}
