package expr

import (
	"context"
	"net/http"
)

// Conditions is a convience function to perform a set of conditions
type Conditions []Condition

func (c Conditions) BuildHandler(ctx context.Context, next Handler) Handler {
	//TODO implement me
	panic("implement me condition build handler")
}

func (c Conditions) BuildCondition(ctx context.Context) ConditionFunc {
	h := func(r *http.Request) bool {
		return true
	}

	for i := len(c); i > 0; i-- {
		cond := c[i-1].BuildCondition(ctx)
		next := h
		h = func(r *http.Request) bool {
			if cond(r) == false {
				return false
			}
			return next(r)
		}
	}
	return h
}

func (c Conditions) Else(action ...Action) Action {
	//TODO implement me
	panic("implement me condition handler else")
}

func (c Conditions) Then(action ...Action) ElseCondition {
	//TODO implement me
	panic("implement me condition handler then")
}
