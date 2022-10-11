package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	. "github.com/mbict/befe/dsl/jwt"
)

func Program() Expr {
	cors := CORS().
		AllowedOrigins(FromEnvWithDefault("API_URI", "http://localhost/")).
		AllowAllMethods().
		AllowedHeaders("Authorization")

	jwkUri := FromEnvWithDefault("JWK_URI", "http://localhost/.well-known/jwks.json")
	jwk := JwkToken(jwkUri).
		WithExpiredCheck().
		WhenExpired(Deny(), WriteJson(JSON{"error": "expired_token"})).
		WhenDenied(Deny(), WriteJson(JSON{"error": "invalid_token"})).
		WhenNoToken(Unauthorized(), WriteJson(JSON{"error": "missing_token"}))

	accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost/"))

	return With(cors, jwk, accept)
}
