package main

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Expr {
	//run this in the start of your program, environment variables are directly available
	DotEnv("./.env")

	//this doesn't do anything but print the env loaded from the dotenv .env file to console
	return With(
		Debug("got this environment variable from the .env file TEST_VAR = ", FromEnv("TEST_VAR")),
		WriteResponse([]byte("TEST_VAR = "+FromEnv("TEST_VAR"))),
	)
}
