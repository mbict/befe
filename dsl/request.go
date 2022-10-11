package dsl

import (
	. "github.com/mbict/befe/expr"
	"net/http"
	"path"
	"strings"
)

func AppendPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		appendPath := builder(r).(string)
		r.URL.Path = path.Join(r.URL.Path, appendPath)
		r.RequestURI = r.URL.String()
		return true, nil
	})
}

func PrependPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		prependPath := builder(r).(string)
		r.URL.Path = path.Join(prependPath, r.URL.Path)
		r.RequestURI = r.URL.String()
		return true, nil
	})
}

func SetPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		r.URL.Path = builder(r).(string)
		r.RequestURI = r.URL.String()
		return true, nil
	})
}

func StripPrefix(prefix string) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)
		return true, nil
	})
}

func RemoveAuthorizationHeader() Action {
	return RemoveHeader("Authorization")
}

func RemoveHeader(name string) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		r.Header.Del(name)
		return true, nil
	})
}

func SetQuery(name string, value interface{}) Action {
	return ActionFunc(
		func(rw http.ResponseWriter, r *http.Request) (bool, error) {
			q := r.URL.Query()
		convert:
			switch v := value.(type) {
			case Valuer:
				value = v(r)
				goto convert
			case string:
				q.Set(name, v)
			case []string:
				q[name] = v
			case []interface{}:
				q[name] = make([]string, len(v))
				for i, param := range v {
					q[name][i] = param.(string)
				}
			default:
				panic("no idea how to convert this")
			}
			r.URL.RawQuery = q.Encode()
			r.RequestURI = r.URL.String()

			return true, nil
		})
}
