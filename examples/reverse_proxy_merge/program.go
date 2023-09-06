package main

import (
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/http"
)

func Program() Expr {

	//rp := http.ReverseProxy("http://localhost/api/v1/")

	api := http.Client(FromEnvWithDefault("API_URI", "http://localhost/api/v1")) //reverse proxy

	return api.Get(UrlBuilder("/accounts")).
		OnSuccess(
			api.Get( //this api is paginated by offset, so we need to visit all pages to get a full result
				UrlBuilder("/bff/product-plans?account_id={account_id}&size={size}",
					ParamFromJsonPath("account_id", "$.data.*.id"), //gather all the account Ids and then fetch all plans for these account locations
					ParamValue("size", 1000),
				)).
				OnSuccess(
					ResultsetMerger(). //merge the current resultset of the lookup back to the first result set
								Target("data.*.product_plans").        //we merge the result into a new field called product_plans
								Source(ValueFromJsonPath("$.data.*")). //this is how we select the results from the source
								Matcher("$.account_id", "$.id"),       //as we fetch plans for multiple locations we need to check if this plan belongs to this location
				).
				OnFailure(
					InternalServerError(),
					Stop(), //we will not continue we will hard stop here
				),
			WriteJson(GetResult()), //write the results to the client
		).
		OnFailure(
			InternalServerError(),
			Stop(), //we will not continue we will hard stop here
		)

}
