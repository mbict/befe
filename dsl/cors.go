package dsl

import (
	"context"
	"github.com/rs/cors"
	"net/http"
)

type CorsAction interface {
	//Condition
	Action

	// AllowAll configures the Cors handler with permissive configuration allowing all
	// origins with all standard methods with any headers.
	AllowAll() CorsAction

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins(...string) CorsAction

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods(...string) CorsAction

	// AllowAllMethods allowa all methods the client is to use with cross-domain requests.
	AllowAllMethods() CorsAction

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders(...string) CorsAction

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders(...string) CorsAction

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge(int) CorsAction

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials() CorsAction

	// DisallowCredentials indicates whether the request should not include user credentials.
	DisallowCredentials() CorsAction

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough() CorsAction
}

func CORS() CorsAction {
	return &corsAction{}
}

type corsAction struct {
	corsOptions cors.Options
}

func (c *corsAction) BuildHandler(ctx context.Context, next Handler) Handler {
	corsHandler := cors.New(c.corsOptions)
	return func(rw http.ResponseWriter, r *http.Request) {
		corsHandler.ServeHTTP(rw, r, next)
	}
}

func (c *corsAction) AllowAll() CorsAction {
	c.corsOptions.AllowedOrigins = []string{"*"}
	c.AllowAllMethods()
	c.corsOptions.AllowedHeaders = []string{"*"}
	c.corsOptions.AllowCredentials = false
	return c
}

func (c *corsAction) AllowAllMethods() CorsAction {
	c.corsOptions.AllowedMethods = []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
	return c
}

func (c *corsAction) AllowedOrigins(origins ...string) CorsAction {
	c.corsOptions.AllowedOrigins = origins
	return c
}

func (c *corsAction) AllowedMethods(methods ...string) CorsAction {
	c.corsOptions.AllowedMethods = methods
	return c
}

func (c *corsAction) AllowedHeaders(headers ...string) CorsAction {
	c.corsOptions.AllowedHeaders = headers
	return c
}

func (c *corsAction) ExposedHeaders(headers ...string) CorsAction {
	c.corsOptions.ExposedHeaders = headers
	return c
}

func (c *corsAction) MaxAge(seconds int) CorsAction {
	c.corsOptions.MaxAge = seconds
	return c
}

func (c *corsAction) AllowCredentials() CorsAction {
	c.corsOptions.AllowCredentials = true
	return c
}

func (c *corsAction) DisallowCredentials() CorsAction {
	c.corsOptions.AllowCredentials = false
	return c
}

func (c *corsAction) OptionsPassthrough() CorsAction {
	c.corsOptions.OptionsPassthrough = true
	return c
}
