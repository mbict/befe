package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	Accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost"))

	jwk := JwkToken("http://localhost/.well-known/jwks.json").
		WithExpiredCheck().
		WhenExpired(Deny(), WriteResponseBody([]byte(`expired token`))).
		WhenDenied(Deny(), WriteResponseBody([]byte(`invalid token`))).
		WhenNoToken(Unauthorized(), WriteResponseBody([]byte(`no token`)))

	//all request needs a valid jwt and should not be expired
	return With(jwk, Accept)
}
