package main

import (
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/http"
	"github.com/mbict/befe/dsl/jwt"
)

func Program() Expr {

	api := http.Client(FromEnvWithDefault("API_URI", "http://localhost/api/v1")) //api client to fetch request of external services

	//cors middleware
	cors := CORS().
		AllowedOrigins("http://localhost:8080").
		AllowAllMethods().
		AllowedHeaders("Authorization")

	//setting up endpoint security with JWK protection
	jwkUri := FromEnvWithDefault("JWK_URI", "http://localhost/.well-known/jwks.json")
	jwk := jwt.JwkToken(jwkUri).
		WithClaim("aud").   //the audience claim should always be available
		WithClaim("sub").   //the subject claim should always be available
		WithExpiredCheck(). //no expired jwt tokens are allowed
		WhenExpired(Deny(), WriteJson(JSON{"error": "expired_token"})).
		WhenDenied(Deny(), WriteJson(JSON{"error": "invalid_token"})).
		WhenNoToken(Unauthorized(), WriteJson(JSON{"error": "missing_token"}))

	router := http.Router().
		With(cors).
		OnNotFound(NotFound())

	routes := router.Group("").With(jwk) //to avoid chaining also on the error handlers we create a group with the jwt middleware
	{
		routes.Get("/profile").Then(
			api.Get(UrlBuilder("/accounts/{accountId}/customers?user_id={userId}",
				jwt.ParamFromClaim("accountId", "aud"), //fetch the accountId from the JWT clain aud (audience)
				jwt.ParamFromClaim("userId", "sub"),    //fetch the user id from the JWT claim sub (subject)
			)).OnSuccess( //do a rest api call to the internal customers api
				JsonPathHas("$.data[?(@.active==true)]").
					Then( //check if the result gives back atleast one response
						Transform( //transform is slightly an optimized way to pass transformers, and also functions as clear placeholder
							JsonPathFirst("$.data[?(@.active==true)]"),                               //get the first active customer
							IncludePath("id", "user_id", "customer_number", "firstname", "lastname"), //we only allow to expose these fields
						),
						Ok(),
						WriteJson(GetResult()),
					).Else(NotFound()), //if there is no data found in the condition we return a 404 response
			).OnFailure(
				Decision().
					When(StatusCodeIs(404)).
					Then(Deny(), WriteJson(JSON{"error": "denied"})).                               //When status code of api response is 404 - we say we denied the request
					Else(InternalServerError(), WriteJson(JSON{"error": "internal_server_error"})), //all other status will result in an internal server error to be logged
			),
		)
	}

	return router
}
