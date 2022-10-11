package dsl

import (
	"github.com/mbict/befe/expr"
	"net/http"
	"strings"
)

// UrlBuilder creates a new url based on the templated and dynamically fills the Params (named Valuers)
// Usage: UrlBuilder("http://test.com/{name}?test={query}")
func UrlBuilder(urlTemplate string, params ...expr.Param) expr.Valuer {
	return func(r *http.Request) interface{} {
		//replace the named params
		lookupUrl := urlTemplate
		for _, v := range params {
			name, value := v(r)
			lookupUrl = strings.ReplaceAll(lookupUrl, "{"+name+"}", expr.ValueToString(value))
		}
		return lookupUrl
	}
}
