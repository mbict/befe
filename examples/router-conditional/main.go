package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	Accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost"))

	jwk := JwkToken(FromEnvWithDefault("JWK_URI", "http://localhost/.well-known/jwks.json")).
		WithExpiredCheck()
	//all request will need a valid jwt token, checked against the jwk keyset provided by the issuer
	router := Http().With(jwk)
	{
		router.OnNotFound(NotFound())

		//if all of the conditional endpoints fail we deny the request
		router.Get("/foo").
			When(HasJwtClaim("role", "administrator")).
			Then(Accept)

		//when the query string contains test, we do something else
		//otherwise the request is accepted
		route := router.Get("/bar")
		{
			//when we have the param foo=test than all is good and we accept the connection
			route.When(Query("foo", "test")).Then(Accept)

			//when we have the param foo=notest than the user is so unauthorized
			route.When(Query("foo", "notest")).Then(Unauthorized())

			//by default deny all, no condition at all
			route.Default(Deny())
		}

		//the route baz is always accepted
		route = router.Get("/baz").Default(Accept)

		route = router.Get("/baz/:id")
		{
			//we deny all the request to id 123
			route.When(PathParam("id", "123")).
				Then(Print("deny for id 123"), Deny())

			//we alter the path param to id 100 when id = 456
			//this will be used in the proxy to call the upstream service
			route.When(PathParam("id", "456")).
				Then(Print("got into param 456 rewrite with accept"), SetPathParam("id", "100"), Accept)

			route.Default(Print("got into default"), Accept)
		}
	}

	//use the router as the entrypoint action handler
	return router
}
