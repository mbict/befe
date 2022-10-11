package main

import (
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/http"
)

func Program() Expr {

	api := http.Client(FromEnvWithDefault("API_URI", "http://localhost"))

	return api.Get(String("/customers")).OnSuccess(
		Transform(
			//get the first active customer
			JsonPathFirst("$.data[?(@.active==true)]"),
			//we only allow to expose these fields
			IncludePath("id", "user_id", "customer_number", "firstname", "lastname"),
		),
		Ok(),
		WriteJson(GetResult()),
	).OnFailure(InternalServerError())
}
