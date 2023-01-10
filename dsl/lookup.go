package dsl

import (
	"context"
	"github.com/mbict/befe/expr"
	"net/http"
)

type lookup struct {
	promise    expr.Promise
	conditions expr.Conditions
	onFail     expr.Actions
}

func (l *lookup) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {

	conditions := l.conditions.BuildCondition(ctx)
	failureHandler := l.onFail.BuildHandler(ctx, nil)

	l.promise.OnSuccess(expr.ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if conditions(r) == false {
			return failureHandler(rw, r)
		}
		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}))
	l.promise.OnFailure(l.onFail...)

	return l.promise.BuildHandler(ctx, nil)
}

func (l *lookup) Must(condition ...expr.Condition) expr.ConditionMiddleware {
	l.conditions = append(l.conditions, condition...)
	return l
}

func (l *lookup) OnFailure(action ...expr.Action) expr.ConditionMiddleware {
	l.onFail = append(l.onFail, action...)
	return l
}

func Lookup(promise expr.Promise) expr.ConditionMiddleware {
	//try to clone the promise
	if cloner, ok := promise.(expr.Cloner[expr.Promise]); ok {
		promise = cloner.Clone()
	}

	return &lookup{
		promise:    promise,
		conditions: nil,
		onFail:     nil,
	}
}
