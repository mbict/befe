package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	return WriteResponseBody([]byte(`hello world !`))
}
