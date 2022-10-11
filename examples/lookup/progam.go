package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

// example lookup with lookup middleware
func Program() Expr {

	apiClient := Client(FromEnvWithDefault("API_URI", "http://localhost:8081"))

	return With(
		Lookup(apiClient.Get(UrlBuilder("/test/{id}", ParamFromQuery("id", "object_id")))).
			Must(
				StatusCodeIs(200),
				JsonPathHasValue("$.customer_id", String("1234")), //fake customer id, you can get it for example from the JWT token
			).
			OnFailure(
				NotFound(),
				WriteJson(JSON{"error": "object_not_found"}),
			),
		WriteResponse([]byte(`success`)),
	)

}
