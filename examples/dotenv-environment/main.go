package main

import (
	"fmt"
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	//run this in the start of your program, environment variables are directly available
	DotEnv("./examples/dotenv-environment/.env")

	fmt.Println("got this environment variable from the .env file API_URI = ", FromEnv("API_URI"))

	//this doesn't do anything but print the env loaded from the dotenv .env file to console
	return WriteResponseBody([]byte(FromEnv("API_URI")))
}
