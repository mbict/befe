package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	//For simple condition checking the Decisions() method can be used
	//It mimics a switch / default statements
	d := Decisions()

	//We accept the request on path `/foo`
	d.When(PathEquals("/foo")).Then(WriteResponseBody([]byte(`hey you!`)))

	//When no match was made we deny the request!
	d.Default(Deny(), WriteResponseBody([]byte(`nope, denied!`)))

	return d
}
