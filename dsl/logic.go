package dsl

import "net/http"

func Or(conditions ...Condition) Condition {
	h := func(r *http.Request) bool {
		return false
	}

	for i := len(conditions); i > 0; i-- {
		cond := conditions[i-1].BuildCondition()
		next := h
		h = func(r *http.Request) bool {
			if cond(r) == true {
				return true
			}
			return next(r)
		}
	}
	return ConditionHandler(h)
}

func And(conditions ...Condition) Condition {
	return Conditions(conditions)
}

func Not(conditions ...Condition) Condition {
	h := Conditions(conditions).BuildCondition()
	return ConditionHandler(func(r *http.Request) bool {
		return !h(r)
	})
}
