package http

import (
	"context"
	"github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/expr"
	"github.com/mbict/befe/utils/router"
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

type HttpRouter interface {
	Action
	PathRouter

	//With adds middleware that run on all routes/request in this router
	With(...Action) HttpRouter

	//Group is a convenience function to group paths that share the same middleware
	Group(path string) RouterGroup

	//OnNotFound defines what actions happens when a route is not found
	OnNotFound(...Action) HttpRouter
}

type Endpoint interface {
	Action

	//With adds middleware actions before the endpoint is executed
	With(...Action) Endpoint

	Then(...Action)
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
		actions:    Actions{},
		middleware: nil,
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

func (r *httpRouter) With(actions ...Action) HttpRouter {
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

	errHandler := r.onNotFound.BuildHandler(ctx, nil)
	router.NotFound = func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		_, err := errHandler(rw, r)
		return err == nil, err
	}

	next = router.BuildHandler(next)
	next = r.middleware.BuildHandler(ctx, next)

	return next
}

func (r *httpRouter) OnNotFound(actions ...Action) HttpRouter {
	r.onNotFound = append(r.onNotFound, actions...)
	return r
}

func Router() HttpRouter {
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
	actions    Actions
	middleware Actions
}

func (e *endpoint) BuildHandler(ctx context.Context, next Handler) Handler {
	//build decision tree
	next = e.actions.BuildHandler(ctx, next)
	//wrap in middleware
	next = e.middleware.BuildHandler(ctx, next)

	return next
}

func (e *endpoint) With(actions ...Action) Endpoint {
	e.middleware = append(e.middleware, actions...)
	return e
}

func (e *endpoint) Then(actions ...Action) {
	e.actions = append(e.actions, actions...)
}

// -- Conditional DSL

//

//	func PathParam(name string, values ...string) Condition {
//		return ConditionFunc(func(r *http.Request) bool {
//			params := httprouter.ParamsFromContext(r.Context())
//			value := params.ByName(name)
//			for _, matchValue := range values {
//				if value == matchValue {
//					return true
//				}
//			}
//			return false
//		})
//	}
//

//func HasPathParam(name string, valuers ...Valuer) Condition {
//	return ConditionFunc(func(r *http.Request) bool {
//		params := httprouter.ParamsFromContext(r.Context())
//		value := params.ByName(name)
//		for _, valuer := range valuers {
//
//			v := valuer(r)
//			switch matchValue := v.(type) {
//			case string:
//				if value == matchValue {
//					return true
//				}
//			case []string:
//				for _, val := range matchValue {
//					if value == val {
//						return true
//					}
//				}
//			default:
//				panic("unhandled type")
//			}
//		}
//		return false
//	})
//}
//
//	func SetPathParam(name string, value interface{}) Action {
//		return HandleBuilder(func(next Handler) Handler {
//			return func(rw http.ResponseWriter, r *http.Request) {
//				//todo implement
//				//panic("not implemented yet: SetPathParam")
//				next(rw, r)
//			}
//		})
//	}
//
//	func JsonParam(name string, values ...string) Condition {
//		return ConditionFunc(func(r *http.Request) bool {
//			params := httprouter.ParamsFromContext(r.Context())
//			value := params.ByName(name)
//			for _, matchValue := range values {
//				if value == matchValue {
//					return true
//				}
//			}
//			return false
//		})
//	}
//
//	func SetJsonParam(name string, value interface{}) Action {
//		return HandleBuilder(func(next Handler) Handler {
//			return func(rw http.ResponseWriter, r *http.Request) {
//				//todo implement
//				panic("not implemented yet: SetJsonParam")
//				next(rw, r)
//			}
//		})
//	}
//

func ParamFromPath(paramName string, pathParamName string) Param {
	return dsl.WithParam(paramName, ValueFromPath(pathParamName))
}

func ValueFromPath(name string) Valuer {
	return func(r *http.Request) interface{} {
		return httprouter.ParamsFromContext(r.Context()).ByName(name)
	}
}
