package dsl

import (
	"context"
	. "github.com/mbict/befe/expr"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"strings"
)

// ConditionCallback is a wrapper f plain go callbacks that should return a true on a match and false when not
func ConditionCallback(fn func(req *http.Request) bool) ConditionHandler {
	return NewConditionHandler(ConditionFunc(fn))
}

// JsonPathHas is a condition to check based on a jsonPath pattern if the received or fetched data (http lookup/ http reverseProxy) has a specific value that matches the pattern
func JsonPathHas(pattern string) ConditionHandler {
	jq, err := jp.ParseString(pattern)
	if err != nil {
		panic("cannot compile JsonPath pattern `" + pattern + "`: " + err.Error())
	}

	return NewConditionHandler(ConditionFunc(func(r *http.Request) bool {
		if res := GetResultBucket(r.Context()); res != nil {
			return jq.Has(res.Data)
		}
		return false
	}))
}

// JsonPathHas is a condition to check based on a jsonPath pattern if the received or fetched data (http lookup/ http reverseProxy) has a specific value that matches the pattern
func JsonPathHasValue(pattern string, valuer Valuer) ConditionHandler {
	jq, err := jp.ParseString(pattern)
	if err != nil {
		panic("cannot compile JsonPath pattern `" + pattern + "`: " + err.Error())
	}

	return NewConditionHandler(ConditionFunc(func(r *http.Request) bool {
		if res := GetResultBucket(r.Context()); res != nil {
			compareValue := valuer(r)

			matches := jq.Get(res.Data)
			for _, match := range matches {
				if match == compareValue {
					return true
				}
			}
		}
		return false
	}))
}

// StatusCodeIs is a condition to check if the fetched data (http lookup/ http reverseProxy) returned a specific http statuscode
func StatusCodeIs(code int) ConditionHandler {
	return NewConditionHandler(ConditionFunc(func(r *http.Request) bool {
		if result := GetResultBucket(r.Context()); result != nil {
			return result.Response.StatusCode == code
		}
		return false
	}))
}

func IsMethod(method ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		for _, m := range method {
			if r.Method == m {
				return true
			}
		}
		return false
	})
}

func PathEquals(path ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		for _, p := range path {
			if r.URL.Path == p {
				return true
			}
		}
		return false
	})
}

func PathStartWith(path ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		for _, p := range path {
			if strings.HasPrefix(r.URL.Path, p) {
				return true
			}
		}
		return false
	})
}

func PathEndsWith(path ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		for _, p := range path {
			if strings.HasSuffix(r.URL.Path, p) {
				return true
			}
		}
		return false
	})
}

func HasCookie(name string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		_, err := r.Cookie(name)
		return err == nil
	})
}

func QueryEquals(name string, values ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		requestValues, ok := r.URL.Query()[name]
		if !ok {
			return false
		}

		for _, requestValue := range requestValues {
			for _, value := range values {
				if value == requestValue {
					return true
				}
			}
		}
		return false
	})
}

func HasQuery(name string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		if _, ok := r.URL.Query()[name]; ok {
			return true
		}
		return false
	})
}

// Any of the conditions test true, then the condition is met
// Works like an OR gate
func Any(conditions ...Condition) Condition {
	return BuildConditionFunc(func(ctx context.Context) ConditionFunc {
		next := func(r *http.Request) bool {
			return false
		}
		for _, condition := range conditions {
			fn := condition.BuildCondition(ctx)
			c := next
			next = func(r *http.Request) bool {
				if fn(r) == true {
					return true
				}
				return c(r)
			}
		}
		return next
	})
}

// Or is an alias for Any
func Or(conditions ...Condition) Condition {
	return Any(conditions...)
}

// All the conditions must test true for this condition to be true
// Works like an AND gate
func All(conditions ...Condition) Condition {
	return BuildConditionFunc(func(ctx context.Context) ConditionFunc {
		next := func(r *http.Request) bool {
			return true
		}
		for _, condition := range conditions {
			fn := condition.BuildCondition(ctx)
			c := next
			next = func(r *http.Request) bool {
				if fn(r) == false {
					return false
				}
				return c(r)
			}
		}
		return next
	})
}

// And is an alias for All
func And(conditions ...Condition) Condition {
	return All(conditions...)
}

// Not negates a condition
func Not(condition Condition) Condition {
	return BuildConditionFunc(func(ctx context.Context) ConditionFunc {
		cond := condition.BuildCondition(ctx)
		return func(r *http.Request) bool {
			return !cond(r)
		}
	})
}
