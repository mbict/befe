package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

func Program() Expr {
	return With(
		SetPath("/foo/bar/baz"),
		ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost:8081")),
	)
}
