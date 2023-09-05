// Code generated by 'yaegi extract github.com/mbict/befe/dsl'. DO NOT EDIT.

package buildin

import (
	"context"
	"github.com/mbict/befe/dsl"
	"github.com/mbict/befe/expr"
	"reflect"
)

func init() {
	Symbols["github.com/mbict/befe/dsl/dsl"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"All":                            reflect.ValueOf(dsl.All),
		"And":                            reflect.ValueOf(dsl.And),
		"Any":                            reflect.ValueOf(dsl.Any),
		"AppendPath":                     reflect.ValueOf(dsl.AppendPath),
		"CORS":                           reflect.ValueOf(dsl.CORS),
		"ConditionCallback":              reflect.ValueOf(dsl.ConditionCallback),
		"Created":                        reflect.ValueOf(dsl.Created),
		"Debug":                          reflect.ValueOf(dsl.Debug),
		"Decision":                       reflect.ValueOf(dsl.Decision),
		"DefaultParam":                   reflect.ValueOf(dsl.DefaultParam),
		"DefaultValue":                   reflect.ValueOf(dsl.DefaultValue),
		"Delay":                          reflect.ValueOf(dsl.Delay),
		"Deny":                           reflect.ValueOf(dsl.Deny),
		"DotEnv":                         reflect.ValueOf(dsl.DotEnv),
		"ExcludePath":                    reflect.ValueOf(dsl.ExcludePath),
		"FileServer":                     reflect.ValueOf(dsl.FileServer),
		"FromEnv":                        reflect.ValueOf(dsl.FromEnv),
		"FromEnvInt":                     reflect.ValueOf(dsl.FromEnvInt),
		"FromEnvWithDefault":             reflect.ValueOf(dsl.FromEnvWithDefault),
		"FromEnvWithDefaultInt":          reflect.ValueOf(dsl.FromEnvWithDefaultInt),
		"GetResult":                      reflect.ValueOf(dsl.GetResult),
		"HandlerCallback":                reflect.ValueOf(dsl.HandlerCallback),
		"HasCookie":                      reflect.ValueOf(dsl.HasCookie),
		"HasHeader":                      reflect.ValueOf(dsl.HasHeader),
		"HasQuery":                       reflect.ValueOf(dsl.HasQuery),
		"HeaderEquals":                   reflect.ValueOf(dsl.HeaderEquals),
		"HeaderHasValue":                 reflect.ValueOf(dsl.HeaderHasValue),
		"Host":                           reflect.ValueOf(dsl.Host),
		"HttpHandlerCallback":            reflect.ValueOf(dsl.HttpHandlerCallback),
		"IncludePath":                    reflect.ValueOf(dsl.IncludePath),
		"InjectOAuth2AccessToken":        reflect.ValueOf(dsl.InjectOAuth2AccessToken),
		"InternalServerError":            reflect.ValueOf(dsl.InternalServerError),
		"IsMethod":                       reflect.ValueOf(dsl.IsMethod),
		"JsonContentType":                reflect.ValueOf(dsl.JsonContentType),
		"JsonPath":                       reflect.ValueOf(dsl.JsonPath),
		"JsonPathFirst":                  reflect.ValueOf(dsl.JsonPathFirst),
		"JsonPathHas":                    reflect.ValueOf(dsl.JsonPathHas),
		"JsonPathHasValue":               reflect.ValueOf(dsl.JsonPathHasValue),
		"Lookup":                         reflect.ValueOf(dsl.Lookup),
		"MiddlewareHandlerCallback":      reflect.ValueOf(dsl.MiddlewareHandlerCallback),
		"Not":                            reflect.ValueOf(dsl.Not),
		"NotFound":                       reflect.ValueOf(dsl.NotFound),
		"OAuth2AccessToken":              reflect.ValueOf(dsl.OAuth2AccessToken),
		"OAuth2AuthorizationAccessToken": reflect.ValueOf(dsl.OAuth2AuthorizationAccessToken),
		"OAuthClientCredentials":         reflect.ValueOf(dsl.OAuthClientCredentials),
		"Ok":                             reflect.ValueOf(dsl.Ok),
		"Or":                             reflect.ValueOf(dsl.Or),
		"ParamFromJsonPath":              reflect.ValueOf(dsl.ParamFromJsonPath),
		"ParamFromQuery":                 reflect.ValueOf(dsl.ParamFromQuery),
		"ParamInt":                       reflect.ValueOf(dsl.ParamInt),
		"ParamString":                    reflect.ValueOf(dsl.ParamString),
		"ParamValue":                     reflect.ValueOf(dsl.ParamValue),
		"PathEndsWith":                   reflect.ValueOf(dsl.PathEndsWith),
		"PathEquals":                     reflect.ValueOf(dsl.PathEquals),
		"PathStartWith":                  reflect.ValueOf(dsl.PathStartWith),
		"PermanentRedirect":              reflect.ValueOf(dsl.PermanentRedirect),
		"PrependPath":                    reflect.ValueOf(dsl.PrependPath),
		"QueryEquals":                    reflect.ValueOf(dsl.QueryEquals),
		"Redirect":                       reflect.ValueOf(dsl.Redirect),
		"RemoveAuthorizationHeader":      reflect.ValueOf(dsl.RemoveAuthorizationHeader),
		"RemoveHeader":                   reflect.ValueOf(dsl.RemoveHeader),
		"RequestURI":                     reflect.ValueOf(dsl.RequestURI),
		"ResponseCode":                   reflect.ValueOf(dsl.ResponseCode),
		"ResultMerger":                   reflect.ValueOf(dsl.ResultMerger),
		"ResultsetMerger":                reflect.ValueOf(dsl.ResultsetMerger),
		"SetHeader":                      reflect.ValueOf(dsl.SetHeader),
		"SetPath":                        reflect.ValueOf(dsl.SetPath),
		"SetQuery":                       reflect.ValueOf(dsl.SetQuery),
		"StatusCodeIs":                   reflect.ValueOf(dsl.StatusCodeIs),
		"Stop":                           reflect.ValueOf(dsl.Stop),
		"String":                         reflect.ValueOf(dsl.String),
		"StripPrefix":                    reflect.ValueOf(dsl.StripPrefix),
		"TemporaryRedirect":              reflect.ValueOf(dsl.TemporaryRedirect),
		"Transform":                      reflect.ValueOf(dsl.Transform),
		"Unauthorized":                   reflect.ValueOf(dsl.Unauthorized),
		"UrlBuilder":                     reflect.ValueOf(dsl.UrlBuilder),
		"ValueFromEnv":                   reflect.ValueOf(dsl.ValueFromEnv),
		"ValueFromEnvWithDefault":        reflect.ValueOf(dsl.ValueFromEnvWithDefault),
		"ValueFromJsonPath":              reflect.ValueOf(dsl.ValueFromJsonPath),
		"ValueFromPattern":               reflect.ValueOf(dsl.ValueFromPattern),
		"ValueFromQuery":                 reflect.ValueOf(dsl.ValueFromQuery),
		"With":                           reflect.ValueOf(dsl.With),
		"WithParam":                      reflect.ValueOf(dsl.WithParam),
		"WriteJson":                      reflect.ValueOf(dsl.WriteJson),
		"WriteResponse":                  reflect.ValueOf(dsl.WriteResponse),

		// type definitions
		"CorsAction":       reflect.ValueOf((*dsl.CorsAction)(nil)),
		"Expr":             reflect.ValueOf((*dsl.Expr)(nil)),
		"JSON":             reflect.ValueOf((*dsl.JSON)(nil)),
		"Merger":           reflect.ValueOf((*dsl.Merger)(nil)),
		"OAuthClientToken": reflect.ValueOf((*dsl.OAuthClientToken)(nil)),

		// interface wrapper definitions
		"_CorsAction":       reflect.ValueOf((*_github_com_mbict_befe_dsl_CorsAction)(nil)),
		"_Expr":             reflect.ValueOf((*_github_com_mbict_befe_dsl_Expr)(nil)),
		"_Merger":           reflect.ValueOf((*_github_com_mbict_befe_dsl_Merger)(nil)),
		"_OAuthClientToken": reflect.ValueOf((*_github_com_mbict_befe_dsl_OAuthClientToken)(nil)),
	}
}

// _github_com_mbict_befe_dsl_CorsAction is an interface wrapper for CorsAction type
type _github_com_mbict_befe_dsl_CorsAction struct {
	IValue               interface{}
	WAllowAll            func() dsl.CorsAction
	WAllowAllMethods     func() dsl.CorsAction
	WAllowCredentials    func() dsl.CorsAction
	WAllowedHeaders      func(a0 ...string) dsl.CorsAction
	WAllowedMethods      func(a0 ...string) dsl.CorsAction
	WAllowedOrigins      func(a0 ...string) dsl.CorsAction
	WBuildHandler        func(ctx context.Context, next expr.Handler) expr.Handler
	WDisallowCredentials func() dsl.CorsAction
	WExposedHeaders      func(a0 ...string) dsl.CorsAction
	WMaxAge              func(a0 int) dsl.CorsAction
	WOptionsPassthrough  func() dsl.CorsAction
}

func (W _github_com_mbict_befe_dsl_CorsAction) AllowAll() dsl.CorsAction {
	return W.WAllowAll()
}
func (W _github_com_mbict_befe_dsl_CorsAction) AllowAllMethods() dsl.CorsAction {
	return W.WAllowAllMethods()
}
func (W _github_com_mbict_befe_dsl_CorsAction) AllowCredentials() dsl.CorsAction {
	return W.WAllowCredentials()
}
func (W _github_com_mbict_befe_dsl_CorsAction) AllowedHeaders(a0 ...string) dsl.CorsAction {
	return W.WAllowedHeaders(a0...)
}
func (W _github_com_mbict_befe_dsl_CorsAction) AllowedMethods(a0 ...string) dsl.CorsAction {
	return W.WAllowedMethods(a0...)
}
func (W _github_com_mbict_befe_dsl_CorsAction) AllowedOrigins(a0 ...string) dsl.CorsAction {
	return W.WAllowedOrigins(a0...)
}
func (W _github_com_mbict_befe_dsl_CorsAction) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {
	return W.WBuildHandler(ctx, next)
}
func (W _github_com_mbict_befe_dsl_CorsAction) DisallowCredentials() dsl.CorsAction {
	return W.WDisallowCredentials()
}
func (W _github_com_mbict_befe_dsl_CorsAction) ExposedHeaders(a0 ...string) dsl.CorsAction {
	return W.WExposedHeaders(a0...)
}
func (W _github_com_mbict_befe_dsl_CorsAction) MaxAge(a0 int) dsl.CorsAction {
	return W.WMaxAge(a0)
}
func (W _github_com_mbict_befe_dsl_CorsAction) OptionsPassthrough() dsl.CorsAction {
	return W.WOptionsPassthrough()
}

// _github_com_mbict_befe_dsl_Expr is an interface wrapper for Expr type
type _github_com_mbict_befe_dsl_Expr struct {
	IValue        interface{}
	WBuildHandler func(ctx context.Context, next expr.Handler) expr.Handler
}

func (W _github_com_mbict_befe_dsl_Expr) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {
	return W.WBuildHandler(ctx, next)
}

// _github_com_mbict_befe_dsl_Merger is an interface wrapper for Merger type
type _github_com_mbict_befe_dsl_Merger struct {
	IValue        interface{}
	WBuildHandler func(ctx context.Context, next expr.Handler) expr.Handler
	WMatcher      func(sourcePattern string, targetPattern string) dsl.Merger
	WSource       func(sourcer expr.Valuer) dsl.Merger
	WTarget       func(path string) dsl.Merger
}

func (W _github_com_mbict_befe_dsl_Merger) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {
	return W.WBuildHandler(ctx, next)
}
func (W _github_com_mbict_befe_dsl_Merger) Matcher(sourcePattern string, targetPattern string) dsl.Merger {
	return W.WMatcher(sourcePattern, targetPattern)
}
func (W _github_com_mbict_befe_dsl_Merger) Source(sourcer expr.Valuer) dsl.Merger {
	return W.WSource(sourcer)
}
func (W _github_com_mbict_befe_dsl_Merger) Target(path string) dsl.Merger {
	return W.WTarget(path)
}

// _github_com_mbict_befe_dsl_OAuthClientToken is an interface wrapper for OAuthClientToken type
type _github_com_mbict_befe_dsl_OAuthClientToken struct {
	IValue        interface{}
	WBuildHandler func(ctx context.Context, next expr.Handler) expr.Handler
	WInjectToken  func() dsl.OAuthClientToken
	WWhenDenied   func(actions ...expr.Action) dsl.OAuthClientToken
	WWhenError    func(actions ...expr.Action) dsl.OAuthClientToken
}

func (W _github_com_mbict_befe_dsl_OAuthClientToken) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {
	return W.WBuildHandler(ctx, next)
}
func (W _github_com_mbict_befe_dsl_OAuthClientToken) InjectToken() dsl.OAuthClientToken {
	return W.WInjectToken()
}
func (W _github_com_mbict_befe_dsl_OAuthClientToken) WhenDenied(actions ...expr.Action) dsl.OAuthClientToken {
	return W.WWhenDenied(actions...)
}
func (W _github_com_mbict_befe_dsl_OAuthClientToken) WhenError(actions ...expr.Action) dsl.OAuthClientToken {
	return W.WWhenError(actions...)
}