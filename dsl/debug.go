package dsl

import (
	"fmt"
	. "github.com/mbict/befe/expr"
	"net/http"
)

func Debug(message string, params ...interface{}) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		//execute all the valuers to their resolved value
		parsedParams := make([]interface{}, len(params))
		for i, param := range params {
			switch v := param.(type) {
			case Valuer:
				parsedParams[i] = v(r)
			case ConditionFunc:
				parsedParams[i] = v(r)
			case Condition:
				parsedParams[i] = v.BuildCondition(r.Context())(r)
			default:
				parsedParams[i] = param
			}
		}

		fmt.Println(append([]interface{}{"[DEBUG]", message}, parsedParams...)...)
		return true, nil
	})
}
