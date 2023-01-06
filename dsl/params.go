package dsl

import (
	"github.com/mbict/befe/expr"
	"net/http"
)

func WithParam(name string, valuer expr.Valuer) expr.Param {
	return func(r *http.Request) (string, interface{}) {
		return name, valuer(r)
	}
}

func ParamValue(name string, value interface{}) expr.Param {
	return func(r *http.Request) (string, interface{}) {
		return name, value
	}
}

func ParamString(name string, value string) expr.Param {
	return func(r *http.Request) (string, interface{}) {
		return name, value
	}
}

func ParamInt(name string, value int) expr.Param {
	return func(r *http.Request) (string, interface{}) {
		return name, value
	}
}

func ParamFromQuery(paramName string, queryParamName string) expr.Param {
	return WithParam(paramName, ValueFromQuery(queryParamName))
}

func ParamFromJsonPath(paramName string, pattern string) expr.Param {
	return WithParam(paramName, ValueFromJsonPath(pattern))
}

func DefaultParam(param expr.Param, defaultValuer expr.Valuer) expr.Param {
	return func(r *http.Request) (string, interface{}) {
		name, value := param(r)
		return name, onEmptyDefaultValue(value, defaultValuer, r)
	}
}
