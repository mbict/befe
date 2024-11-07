package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	"github.com/mbict/befe/dsl/oidc"
)

func Program() Expr {
	frontendEndpoint := ReverseProxy("http://localhost:3000")

	SSO := oidc.SingleSignOn(
		"https://demo.dev.stalling.app",
		"fd27c635-9ed0-4fcb-baf6-524e1b1065bb",
		"secret",
		"http://localhost:8080",
		oidc.WithCookieHttpOnly(false),
	).WithSameIssuer().
		//WhenExpired(tokenExpired()). //let it automatically refresh
		WhenDenied(invalidToken()).
		WhenInvalidToken(invalidToken()).
		WhenNoToken(oidc.AuthTokenRedirect())

	return With(SSO, frontendEndpoint)
}

func invalidToken() Expr {
	return With(
		Deny(),
		WriteResponse([]byte(`your session/token is invalid`)),
	)
}
