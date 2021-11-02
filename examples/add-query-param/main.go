package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	return With(
		SetQuery("foo", []string{"bar", "baz"}),
		ReverseProxy("http://localhost"),
	)
}
