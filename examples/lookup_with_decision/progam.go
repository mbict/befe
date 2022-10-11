package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

// example lookup with decision
func Program() Expr {

	apiClient := Client(FromEnvWithDefault("API_URI", "http://localhost:8081"))

	return apiClient.Get(UrlBuilder("/test/{id}", ParamFromQuery("id", "object_id"))).
		OnSuccess(
			Decision().
				When(
					StatusCodeIs(200),
					JsonPathHasValue("$.customer_id", String("1234")), //fake customer id, you can get it for example from the JWT token
				).Then(
				WriteResponse([]byte(`success`)),
			).Else(
				NotFound(),
				WriteJson(JSON{"error": "object_not_found"}),
			),
		).OnFailure(
		NotFound(),
		WriteJson(JSON{"error": "object_not_found"}),
	)

}
