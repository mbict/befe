package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Expr {
	//For simple condition checking the Decisions() method can be used
	//It mimics a switch / default statements
	d := Decision()

	//We accept the request on path `/foo`
	d.When(PathEquals("/foo")).Then(WriteResponse([]byte(`hey you!`)))

	//When no match was made we deny the request!
	d.Else(Deny(), WriteResponse([]byte(`nope, denied!`)))

	return d
}
