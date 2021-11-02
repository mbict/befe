package dsl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	httprouter "github.com/mbict/befe/dsl/router"
	"github.com/mbict/befe/utils/token"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Valuer func(r *http.Request) interface{}

func ValueFromEnv(name string) Valuer {
	return func(r *http.Request) interface{} {
		return FromEnv(name)
	}
}

func ValueFromEnvWithDefault(name string, defaultValue string) Valuer {
	return func(r *http.Request) interface{} {
		return FromEnvWithDefault(name, defaultValue)
	}
}

func ValueFromPath(name string) Valuer {
	return func(r *http.Request) interface{} {
		return httprouter.ParamsFromContext(r.Context()).ByName(name)
	}
}

func ParamFromPath(paramName string, pathParamName string) Param {
	return WithParam(paramName, ValueFromPath(pathParamName))
}

func ValueFromQuery(name string) Valuer {
	return func(r *http.Request) interface{} {
		v := r.URL.Query()[name]
		if len(v) == 0 {
			return ""
		}
		if len(v) == 1 {
			return v[0]
		}
		return v
	}
}

func ParamFromQuery(paramName string, queryParamName string) Param {
	return WithParam(paramName, ValueFromQuery(queryParamName))
}

func ValueFromJwtClaim(name string) Valuer {
	return func(r *http.Request) interface{} {
		t, ok := r.Context().Value(&jwtContextKey).(*token.JwtToken)
		if !ok || t == nil {
			return nil
		}

		v, ok := t.Get(name)
		if !ok {
			return nil
		}

		/*
			//string value
			if value, ok := v.(string); ok {
				return value
			}

			//slice value
			if values, ok := v.([]string); ok {
				return values
			}
		*/

		return v
	}
}

func ValueFromPattern(regex string, valuer Valuer) Valuer {
	rxp := regexp.MustCompile(regex)
	extract := func(in string) string {
		if res := rxp.FindStringSubmatch(in); res != nil && len(res) >= 2 {
			return res[1]
		}
		return ""
	}

	return func(r *http.Request) interface{} {
		switch v := valuer(r).(type) {
		case string:
			return extract(v)
		case []string:
			for _, sv := range v {
				if p := extract(sv); p != "" {
					return p
				}
			}
		}
		return ""
	}
}

func ParamString(name string, value string) Param {
	return func(r *http.Request) (string, interface{}) {
		return name, value
	}
}

func ParamFromJwtClaim(paramName string, claimName string) Param {
	return WithParam(paramName, ValueFromJwtClaim(claimName))
}

type Lookup interface {
	Action
	Condition

	Must(...ResponseCondition) Lookup
	OnFail(...Action) Lookup
}

type httpLookup struct {
	client       http.Client
	pathTemplate Valuer
	conditions   []ResponseCondition
	onFail       Actions
}

func (h *httpLookup) BuildCondition() ConditionHandler {
	return func(req *http.Request) bool {
		resp := h.doLookup(req)
		return h.checkConditions(resp, req)
	}

}

func (h *httpLookup) doLookup(req *http.Request) *http.Response {
	url := h.pathTemplate(req).(string)
	resp, err := h.client.Get(url)

	if err != nil {
		panic(err)
	}
	return resp
}

func (h *httpLookup) checkConditions(resp *http.Response, req *http.Request) bool {
	for _, cond := range h.conditions {
		if false == cond(resp, req) {
			return false
		}
	}
	return true
}

func (h *httpLookup) BuildHandler(ctx context.Context, next Handler) Handler {
	onFailHandler := h.onFail.BuildHandler(ctx, emptyHandler)
	return func(rw http.ResponseWriter, req *http.Request) {
		resp := h.doLookup(req)
		if h.checkConditions(resp, req) == false {
			onFailHandler(rw, req)
			return
		}
		next(rw, req)
	}
}

func (h *httpLookup) Must(condition ...ResponseCondition) Lookup {
	h.conditions = append(h.conditions, condition...)
	return h
}

func (h *httpLookup) OnFail(actions ...Action) Lookup {
	h.onFail = append(h.onFail, actions...)
	return h
}

func HttpLookup(urlTemplate string, params ...Param) Lookup {
	return &httpLookup{
		client:       http.Client{},
		pathTemplate: UrlBuilder(urlTemplate, params...),
		conditions:   nil,
		onFail:       nil,
	}
}

//UrlBuilder creates a new url based on the templated and dynamically filled from Params (named Valuers)
//Usage: UrlBuilder("http://test.com/{name}?test={query}")
func UrlBuilder(urlTemplate string, params ...Param) Valuer {
	return func(r *http.Request) interface{} {
		//replace the named params
		lookupUrl := urlTemplate
		for _, v := range params {
			name, value := v(r)
			lookupUrl = strings.ReplaceAll(lookupUrl, "{"+name+"}", valueToString(value))
		}
		return lookupUrl
	}
}

func valueToString(in interface{}) string {
	if in == nil {
		return ""
	}

	switch v := in.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, `,`)
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v == true {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%s", v)
	}
}

type Param func(r *http.Request) (string, interface{})

func WithParam(name string, valuer Valuer) Param {
	return func(r *http.Request) (string, interface{}) {
		return name, valuer(r)
	}
}

type ResponseCondition func(*http.Response, *http.Request) bool

func IsSuccessfulResponse() ResponseCondition {
	return func(response *http.Response, r *http.Request) bool {
		return response.StatusCode >= 200 && response.StatusCode <= 299
	}
}

func HaveResponseCode(code int) ResponseCondition {
	return func(response *http.Response, _ *http.Request) bool {
		return response.StatusCode == code
	}
}

func BodyContainsElements(path string) ResponseCondition {
	return func(response *http.Response, req *http.Request) bool {
		return true
	}
}

type JsonResponseCondition func(interface{}, *http.Response, *http.Request) bool

func JsonResponse(conditions ...JsonResponseCondition) ResponseCondition {
	return func(response *http.Response, request *http.Request) bool {
		decoder := json.NewDecoder(response.Body)
		document := make(map[string]interface{})
		if err := decoder.Decode(&document); err != nil {
			panic(fmt.Errorf("unable to decode json body: %w", err))
		}

		for _, condition := range conditions {
			if condition(document, response, request) == false {
				return false
			}
		}
		return true
	}
}

func JsonHasValue(jpath string, valuer Valuer) JsonResponseCondition {
	//https://pkg.go.dev/github.com/PaesslerAG/jsonpath#example-Get
	path, err := jsonpath.New(jpath)
	if err != nil {
		panic(fmt.Errorf("unable to compile jsonpath expression: %w", err))
	}

	return func(document interface{}, response *http.Response, req *http.Request) bool {
		value, err := path(req.Context(), document)
		if err != nil {
			return false
		}
		return valueToString(value) == valueToString(valuer(req))
	}
}
