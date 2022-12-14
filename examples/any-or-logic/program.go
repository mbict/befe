package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Expr {

	return Decision().

		//All condition
		When(All(
			PathStartWith("/test"), // (takes precedence over next any condition)
			QueryEquals("all", "true"),
		)).
		Then(WriteResponse([]byte(`All matched`))).

		//any condition
		When(Any(
			PathStartWith("/test"),
			PathEquals("/blup"),
		)).
		Then(WriteResponse([]byte(`Any matched`))).

		//if all fail
		Else(WriteResponse([]byte(`None matched`)))
}
