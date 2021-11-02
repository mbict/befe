package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	return With(
		SetPath("/foo/bar/baz"),
		ReverseProxy("http://localhost"),
	)
}
