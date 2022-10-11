package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

// this example of the reverse proxy will directly passthrough the response without any modification
func Program() Expr {
	Cors := CORS().AllowAll()
	Accept := ReverseProxy(FromEnvWithDefault("API_URI", "http://localhost:8081"))

	return With(Cors, Accept) //<- be aware that the accept call is the last in the stack and no next actions will be added, it is internally optimized to not read the body and parse and pass it, but directly write it to the output
}
