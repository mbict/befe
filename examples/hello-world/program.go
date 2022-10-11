package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Expr {
	return WriteResponse([]byte(`hello world`))
}
