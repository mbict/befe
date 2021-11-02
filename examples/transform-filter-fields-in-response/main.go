package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	Accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost"))

	//This example shows that we transform the response from the downstream
	router := Http()
	{
		//filtered with action
		router.Get("/filtered").
			Default(
				Print("filtered"),
				TransformResponse(
					Json(
						Filter("data.*.baz", "data.*.foo"),
					),
				),
				Accept)

		//filtered with transform on endpoint
		router.Get("/filtered-alternative").
			Default(Accept).
			WithTransform(
				Json(
					Filter("data.*.baz"),
				))

		//unfiltered
		router.Get("/").Default(Print("default"), Accept)

		//unkown path we reject with 404
		router.OnNotFound(NotFound())
	}
	return router
}
