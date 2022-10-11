package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

func Program() Expr {
	apiClient := Client(FromEnv("API_URI"))

	//oauth token creation, for injecting it into the client call
	oauthClient := OAuthClientCredentials(
		FromEnv("OAUTH2_CLIENT_ID"),
		FromEnv("OAUTH2_CLIENT_SECRET"),
		FromEnv("OAUTH2_TOKEN_URL"),
		[]string{},
	).WhenDenied(Deny(), WriteJson(JSON{"error": "access denied, obtaining access token", "code": "error_obtaining_access_token"})).
		WhenError(InternalServerError(), WriteJson(JSON{"error": "internal error, obtaining access token", "code": "internal_error"}))

	//Inject the oauth token into the client header
	apiCall := apiClient.WithHeader("Authorization", OAuth2AuthorizationAccessToken()).
		Get(String("/test")).
		OnSuccess(WriteJson(GetResult())).
		OnFailure(InternalServerError(), WriteJson(JSON{"error": "fetch"}))

	//if the JWT is valid, we create an oauth access token with the client credentials flow and
	//pass that to the target service
	return With(oauthClient, apiCall)
}
