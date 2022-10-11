package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
)

func Program() Expr {

	fs := FileServer("./dist/").
		AllowDirectoryListing(false).
		WhenNotFound(NotFound(), WriteResponse([]byte("the page does not exist 404 html here")))

	defaultIndex := fs.ServeFile("index.html")

	//We have history mode on for the vue route application, and we allow to navigate directly to a specific page.
	//You could also directly use the fileServer and always serve the index.html on the notFound,
	//but I wanted that non-existing paths should return a real 404 from the server.

	router := Router()
	{
		router.Get("/about").Then(defaultIndex)
		router.Get("/test/*filepath").Then(defaultIndex)
		router.Get("/test").Then(defaultIndex)

		router.OnNotFound(fs)
	}
	return router
}
