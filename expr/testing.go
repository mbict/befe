package expr

import (
	"context"
	"net/http"
)

func MockedCondition(returns bool, callbacks ...func()) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		for _, cf := range callbacks {
			cf()
		}
		return returns
	})
}

func MockedAction(continues bool, returnsError error, callbacks ...func()) Action {
	return ActionFunc(func(_ http.ResponseWriter, _ *http.Request) (bool, error) {
		for _, cf := range callbacks {
			cf()
		}
		return continues, returnsError
	})
}

func MockedHandler(continues bool, returnsError error, callbacks ...func()) Handler {
	return MockedAction(continues, returnsError, callbacks...).BuildHandler(context.Background(), nil)
}
