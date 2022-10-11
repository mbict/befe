package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	. "github.com/mbict/befe/dsl/jwt"
)

// Program is an example on how to project an server with JWT while the upstream server requires a fixed set of oauth2 client credentials
func Program() Expr {

	//endpoint jwk protection
	jwk := JwkToken(FromEnv("JWK_URI")).
		WithExpiredCheck().
		WhenExpired(Deny(), WriteJson(JSON{"error": "expired access token", "code": "expired_token"})).
		WhenDenied(Deny(), WriteJson(JSON{"error": "invalid or malformed access token", "code": "malformed_token"})).
		WhenNoToken(Unauthorized(), WriteJson(JSON{"error": "access token required", "code": "access_token_required"}))

	//oauth token creation, for injecting it into the request to the reverse proxy
	oauthClient := OAuthClientCredentials(
		FromEnv("OAUTH2_CLIENT_ID"),
		FromEnv("OAUTH2_CLIENT_SECRET"),
		FromEnv("OAUTH2_TOKEN_URL"),
		[]string{},
	).
		InjectToken().
		WhenDenied(Deny(), WriteJson(JSON{"error": "access denied, obtaining access token", "code": "error_obtaining_access_token"})).
		WhenError(InternalServerError(), WriteJson(JSON{"error": "internal error, obtaining access token", "code": "internal_error"}))

	//the target server where the token should be injected
	reverseProxy := ReverseProxy(FromEnv("API_URI"))

	//if the JWT is valid, we create an oauth access token with the client credentials flow and
	//pass that to the target service
	return With(jwk, oauthClient, reverseProxy)
}
