package dsl

import (
	"github.com/mbict/befe/expr"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"regexp"
)

func String(value string) expr.Valuer {
	return func(r *http.Request) interface{} {
		return value
	}
}

func Host() expr.Valuer {
	return func(r *http.Request) interface{} {
		return r.Host
	}
}

func ValueFromEnv(name string) expr.Valuer {
	return func(r *http.Request) interface{} {
		return FromEnv(name)
	}
}

func ValueFromEnvWithDefault(name string, defaultValue string) expr.Valuer {
	return func(r *http.Request) interface{} {
		return FromEnvWithDefault(name, defaultValue)
	}
}

func ValueFromJsonPath(pattern string) expr.Valuer {
	jpq, err := jp.ParseString(pattern)
	if err != nil {
		panic(err)
	}

	return func(r *http.Request) interface{} {
		bucket := expr.GetResultBucket(r.Context())
		if bucket != nil {
			return jpq.Get(bucket.Data)
		}
		return nil
	}
}

func ValueFromQuery(name string) expr.Valuer {
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

func ValueFromPattern(regex string, valuer expr.Valuer) expr.Valuer {
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

func DefaultValue(valuer expr.Valuer, defaultValue expr.Valuer) expr.Valuer {
	return func(r *http.Request) interface{} {
		return onEmptyDefaultValue(valuer(r), defaultValue, r)
	}
}

func onEmptyDefaultValue(in interface{}, defaultValuer expr.Valuer, r *http.Request) interface{} {
	switch v := in.(type) {
	case string:
		if v != "" {
			return v
		}
	case int:
		if v != 0 {
			return v
		}
	default:
		if in != nil {
			return in
		}
	}
	return defaultValuer(r)
}
