package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	Cors := CORS().AllowAll()
	Accept := ReverseProxy(FromEnv("API_URI"))

	return With(Cors, Accept)
}
