package dsl

import (
	. "github.com/mbict/befe/expr"
)

func FileServer(path string) HTTPFileServer {
	return NewFileServer(path)
}
