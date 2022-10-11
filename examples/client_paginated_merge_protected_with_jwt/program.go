package main

import (
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/http"
)

func Program() Expr {

	var maxResultsPerRequest = FromEnvWithDefaultInt("MAX_RESULTS_PER_REQUEST", 10)

	cors := CORS().
		AllowedOrigins("http://localhost:8080").
		AllowAllMethods().
		AllowedHeaders("Authorization")

	api := http.Client(FromEnvWithDefault("API_URI", "http://localhost/api/v1")) //reverse proxy

	router := http.Router().With(cors)
	{
		router.OnNotFound(NotFound())

		router.Get("/locations").Then( //get the list with all account locations
			api.Get(UrlBuilder("/accounts")).
				OnSuccess(
					http.OffsetPaginatedResults(maxResultsPerRequest, ValueFromJsonPath("$.metadata.total"),
						api.Get( //this api is paginated by offset, so we need to visit all pages to get a full result
							UrlBuilder("/bff/product-plans?account_id={account_id}&offset={offset}&size={size}",
								ParamFromJsonPath("account_id", "$.data.*.id"), //gather all the account Ids and then fetch all plans for these account locations
								http.ParamFromPaginator("offset", "offset"),    //we fetch the current offset based on the pagination helper
								ParamValue("size", maxResultsPerRequest),
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
					),
					WriteJson(GetResult()), //write the results to the client
				).
				OnFailure(
					InternalServerError(),
					Stop(), //we will not continue we will hard stop here
				),
		)

		router.Get("/locations/{accountId}").Then( //get the list with all account locations
			api.Get(UrlBuilder("/accounts/{accountId}")).
				OnSuccess(
					http.OffsetPaginatedResults(1000, ValueFromJsonPath("$.data.metadata.total"), api.Get( //this api is paginated by offset, so we need to visit all pages to get a full result
						UrlBuilder("/bff/product-plans?account_id={account_id}&offset={offset}&size=1000",
							ParamFromJsonPath("accountIds", "$.data.*.id"), //gather all the account Ids and then fetch all plans for these account locations
							http.ParamFromPaginator("offset", "offset"),    //we fetch the current offset based on the pagination helper
						)).
						OnSuccess(
							ResultsetMerger(). //merge the current resultset of the lookup back to the first result set
										Target("product_plans").             //this is the target field/path to put our results into
										Source(ValueFromJsonPath("$.data")), //this is how we select the results from the source
						).
						OnFailure(
							InternalServerError(),
						),
					),
				).
				OnFailure(
					InternalServerError(),
				),
		)
	}

	return router
}
