package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

// this example of the reverse proxy can modify and transform the response body
func Program() Expr {
	Cors := CORS().AllowAll()
	Accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost:8081"))

	return With(Cors, Accept, Transform(IncludePath("abc")))
}
