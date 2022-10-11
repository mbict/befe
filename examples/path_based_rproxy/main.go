package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	"net/http"
)

func Program() Expr {
	backend := ReverseProxy(FromEnvWithDefault("BACKEND_URI", "http://localhost:1323"))
	frontend := ReverseProxy(FromEnvWithDefault("FRONTEND_URI", "http://localhost:3000"))

	//For simple condition checking the Decisions() method can be used
	//It mimics a switch / default statements
	d := Decision()

	d.When(
		IsMethod(http.MethodGet, http.MethodOptions),
		PathEquals("/login"),
		HasCookie("ssid"),
	).Then(SetPath("/check-session"), backend)

	d.When(
		IsMethod(http.MethodPost, http.MethodOptions),
		PathEquals("/login", "/register", "/recover", "/verify"),
	).Then(backend)

	d.When(
		IsMethod(http.MethodGet, http.MethodOptions),
		PathEquals("/consent", "/logout"),
	).Then(backend)

	d.When(
		PathEquals("/not-here"),
	).Then(TemporaryRedirect("/login"))

	d.Else(frontend)

	return d
}
