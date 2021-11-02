package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	cors := CORS().
		AllowedOrigins("http://localhost:8080").
		AllowAllMethods().
		AllowCredentials()
	jwk := JwkToken("http://localhost/.well-known/jwks.json").
		WithExpiredCheck().
		WhenExpired(Deny(), WriteResponseBody([]byte(`expired token`))).
		WhenDenied(Deny(), WriteResponseBody([]byte(`invalid token`))).
		WhenNoToken(Unauthorized(), WriteResponseBody([]byte(`no token`)))
	accept := ReverseProxy("http://localhost:8090")

	return With(cors, jwk, accept)
}
