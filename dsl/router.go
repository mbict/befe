package dsl

import (
	"context"
	httprouter "github.com/mbict/befe/dsl/router"
	"net/http"
	"path"
)

type PathRouter interface {
	Get(path string) Endpoint
	Post(path string) Endpoint
	Put(path string) Endpoint
	Patch(path string) Endpoint
	Delete(path string) Endpoint
	Connect(path string) Endpoint
	Options(path string) Endpoint
	All(path string) Endpoint
	Method(path string, methods ...string) Endpoint
}

type RouterGroup interface {
	PathRouter

	//With adds middleware that run on all routes/request in this group
	With(...Action) RouterGroup

	//Group is a convenience function to group paths that share the same middleware
	Group(path string) RouterGroup
}

type Router interface {
	Action
	PathRouter

	//With adds middleware that run on all routes/request in this router
	With(...Action) Router

	//Group is a convenience function to group paths that share the same middleware
	Group(path string) RouterGroup

	//OnNotFound defines what actions happens when a route is not found
	OnNotFound(...Action) Router
}

type Endpoint interface {
	Action

	//With adds middleware actions before the endpoint is executed
	With(...Action) Endpoint
	//When creates a decision what action to perform if all conditions match
	When(...Condition) Decision
	//Default action is executed if all the decision tree fails to find an action
	Default(...Action) Endpoint
	//WithTransform adds transformers that will modify the response for all routes
	WithTransform(...Transformer) Endpoint
	//DefaultTransform adds transformers that will modify the response for the default route
	//DefaultTransform(...Transformer) Endpoint
}

type endpoints map[string]map[string]Endpoint

type httpRouterGroup struct {
	middleware Actions
	endpoints  endpoints
	groups     map[string]*httpRouterGroup
}

func (r *httpRouterGroup) With(actions ...Action) RouterGroup {
	r.middleware = append(r.middleware, actions...)
	return r
}

func (r *httpRouterGroup) Get(path string) Endpoint {
	return r.Method(path, "GET")
}

func (r *httpRouterGroup) Post(path string) Endpoint {
	return r.Method(path, "POST")
}

func (r *httpRouterGroup) Put(path string) Endpoint {
	return r.Method(path, "PUT")
}

func (r *httpRouterGroup) Patch(path string) Endpoint {
	return r.Method(path, "PATCH")
}

func (r *httpRouterGroup) Delete(path string) Endpoint {
	return r.Method(path, "DELETE")
}

func (r *httpRouterGroup) Connect(path string) Endpoint {
	return r.Method(path, "CONNECT")
}

func (r *httpRouterGroup) Options(path string) Endpoint {
	return r.Method(path, "OPTIONS")
}

func (r *httpRouterGroup) All(path string) Endpoint {
	return r.Method(path, "GET", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS")
}

func (r *httpRouterGroup) Method(path string, methods ...string) Endpoint {
	ep := &endpoint{
		decisionTree: Decisions(),
		middleware:   nil,
	}

	if v, ok := r.endpoints[path]; !ok || v == nil {
		r.endpoints[path] = make(map[string]Endpoint)
	}
	for _, method := range methods {
		r.endpoints[path][method] = ep
	}
	return ep
}

func (r *httpRouterGroup) Group(path string) RouterGroup {
	group := &httpRouterGroup{
		middleware: nil,
		endpoints:  make(map[string]map[string]Endpoint),
		groups:     make(map[string]*httpRouterGroup),
	}
	r.groups[path] = group

	return group
}

type httpRouter struct {
	httpRouterGroup

	onNotFound Actions
}

func (r *httpRouter) With(actions ...Action) Router {
	r.middleware = append(r.middleware, actions...)
	return r
}

func (r *httpRouter) BuildHandler(ctx context.Context, next Handler) Handler {

	router := httprouter.New()

	buildEndpoints := func(basePath string, endpoints endpoints, middleware Actions) {
		for endpointPath, methods := range endpoints {
			for method, e := range methods {
				router.Handle(
					method,
					path.Join(basePath, endpointPath),
					middleware.BuildHandler(ctx, e.BuildHandler(ctx, next)),
				)
			}
		}
	}

	var buildGroups func(basePath string, groups map[string]*httpRouterGroup)
	buildGroups = func(basePath string, groups map[string]*httpRouterGroup) {
		for groupPath, group := range groups {
			groupBasePath := path.Join(basePath, groupPath)

			buildEndpoints(groupBasePath, group.endpoints, group.middleware)
			buildGroups(groupBasePath, group.groups)
		}
	}

	//build form the root
	buildEndpoints("", r.endpoints, nil)
	buildGroups("", r.groups)

	router.NotFound = r.onNotFound.BuildHandler(ctx, func(rw http.ResponseWriter, r *http.Request) {})

	next = router.BuildHandler(next)
	next = r.middleware.BuildHandler(ctx, next)

	return next
}

func (r *httpRouter) OnNotFound(actions ...Action) Router {
	r.onNotFound = append(r.onNotFound, actions...)
	return r
}

func Http() Router {
	return &httpRouter{
		httpRouterGroup: httpRouterGroup{
			middleware: nil,
			endpoints:  make(map[string]map[string]Endpoint),
			groups:     make(map[string]*httpRouterGroup),
		},
		onNotFound: nil,
	}
}

type endpoint struct {
	decisionTree           DecisionTree
	middleware             Actions
	middlewareTransformers Transformers
}

func (e *endpoint) BuildHandler(ctx context.Context, next Handler) Handler {
	//build decision tree
	next = e.decisionTree.BuildHandler(ctx, next)
	//wrap in middleware
	next = e.middleware.BuildHandler(ctx, next)

	//add middleware transformer that will be applied to all routes
	if e.middlewareTransformers != nil && len(e.middlewareTransformers) > 0 {
		next = TransformResponse(e.middlewareTransformers.Build()).BuildHandler(ctx, next)
	}

	return next
}

func (e *endpoint) With(actions ...Action) Endpoint {
	e.middleware = append(e.middleware, actions...)
	return e
}

func (e *endpoint) When(conditions ...Condition) Decision {
	return e.decisionTree.When(conditions...)
}

func (e *endpoint) Default(actions ...Action) Endpoint {
	e.decisionTree.Default(actions...)
	return e
}

func (e *endpoint) WithTransform(transformers ...Transformer) Endpoint {
	e.middlewareTransformers = append(e.middlewareTransformers, transformers...)
	return e
}

//-- Conditional DSL
func Query(name string, values ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
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

func SetQuery(name string, value interface{}) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
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

			next(rw, r)
		}
	})
}

func RemoveAuthorizationHeader() Action {
	return RemoveHeader("Authorization")
}

func RemoveHeader(name string) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			r.Header.Del(name)
			next(rw, r)
		}
	})
}

func PathParam(name string, values ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		params := httprouter.ParamsFromContext(r.Context())
		value := params.ByName(name)
		for _, matchValue := range values {
			if value == matchValue {
				return true
			}
		}
		return false
	})
}

func HasPathParam(name string, valuers ...Valuer) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		params := httprouter.ParamsFromContext(r.Context())
		value := params.ByName(name)
		for _, valuer := range valuers {

			v := valuer(r)
			switch matchValue := v.(type) {
			case string:
				if value == matchValue {
					return true
				}
			case []string:
				for _, val := range matchValue {
					if value == val {
						return true
					}
				}
			default:
				panic("unhandled type")
			}
		}
		return false
	})
}

func SetPathParam(name string, value interface{}) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			//todo implement
			//panic("not implemented yet: SetPathParam")
			next(rw, r)
		}
	})
}

func JsonParam(name string, values ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		params := httprouter.ParamsFromContext(r.Context())
		value := params.ByName(name)
		for _, matchValue := range values {
			if value == matchValue {
				return true
			}
		}
		return false
	})
}

func SetJsonParam(name string, value interface{}) Action {
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			//todo implement
			panic("not implemented yet: SetJsonParam")
			next(rw, r)
		}
	})
}

func AppendPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			appendPath := builder(r).(string)
			r.URL.Path = path.Join(r.URL.Path, appendPath)
			r.RequestURI = r.URL.String()
			next(rw, r)
		}
	})
}

func PrependPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			prependPath := builder(r).(string)
			r.URL.Path = path.Join(prependPath, r.URL.Path)
			r.RequestURI = r.URL.String()
			next(rw, r)
		}
	})
}

func SetPath(pathTemplate string, params ...Param) Action {
	builder := UrlBuilder(pathTemplate, params...)
	return HandleBuilder(func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			r.URL.Path = builder(r).(string)
			r.RequestURI = r.URL.String()
			next(rw, r)
		}
	})
}

func StripPrefix(prefix string) Action {
	return HandleBuilder(func(next Handler) Handler {
		return http.StripPrefix(prefix, next).ServeHTTP
	})
}

func IsMethod(method ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		for _, m := range method {
			if r.Method == m {
				return true
			}
		}
		return false
	})
}

func PathEquals(path ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		for _, p := range path {
			if r.URL.Path == p {
				return true
			}
		}
		return false
	})
}

func HasCookie(name string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		_, err := r.Cookie(name)
		return err == nil
	})
}
