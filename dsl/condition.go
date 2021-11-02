package dsl

import "net/http"

type Conditions []Condition

func (c Conditions) BuildCondition() ConditionHandler {
	h := func(r *http.Request) bool {
		return true
	}

	for i := len(c); i > 0; i-- {
		cond := c[i-1].BuildCondition()
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

type ConditionHandler func(r *http.Request) bool

func (c ConditionHandler) BuildCondition() ConditionHandler {
	return c
}

type Condition interface {
	BuildCondition() ConditionHandler
}

//func (c ConditionHandler) BuildHandler(next Handler) Handler {
//	return func(r *http.Request) {
//		//Handler(c)(r)
//		next(r)
//	}
//}
