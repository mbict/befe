package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Expr {
	return WriteJson(JSON{"hello": "world", "abc": "foo"})
}
