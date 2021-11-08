package main

import (
	"fmt"
	. "github.com/mbict/befe/dsl"
	"net/http"
)

func Program() Action {
	return HandlerCallback(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("any of your custom http handler go code can live here")
	})
}
