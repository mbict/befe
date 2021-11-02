package vue_file_server

import (
	. "github.com/mbict/befe/dsl"
)

func Program() Action {
	fs := FileServer("/home/michael/Projects/befe/examples/vue-file-server/dist/").
		AllowDirectoryListing(false).
		WhenNotFound(NotFound(), WriteResponseBody([]byte("the page does not exist 404 html here")))

	defaultIndex := fs.ServeFile("index.html")

	//we have history mode for vue application and we allow a few specific paths directly
	//you could also directly use the fileServer and always server the index.html on the notFound
	//but i wanted that non-existing paths should return a real 404 from the server.
	router := Http()
	{
		router.Get("/accounts/*filepath").Default(defaultIndex)
		router.Get("/accounts").Default(defaultIndex)
		router.Get("/users/*filepath").Default(defaultIndex)
		router.Get("/users").Default(defaultIndex)
		router.Get("/financial/*filepath").Default(defaultIndex)
		router.Get("/financial").Default(defaultIndex)
		router.Get("/settings/*filepath").Default(defaultIndex)
		router.Get("/settings").Default(defaultIndex)

		router.OnNotFound(fs)
	}
	return router
}
