package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	"github.com/mbict/befe/dsl/oidc"
)

func Program() Expr {
	frontend := ReverseProxy(FromEnvWithDefault("FRONTEND_URI", "http://localhost:3000"))
	cors := CORS().AllowAll()

	sso := oidc.SingleSignon(
		FromEnvWithDefault("BACKEND_URI", "http://localhost:8080"),
		"webapp",
		"secret")

	return With(cors, sso, frontend)
}
