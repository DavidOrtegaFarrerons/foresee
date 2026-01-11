package main

import (
	"foresee/internal/web"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(
		web.NeuteredFileSystem(http.Dir("./ui/static")),
	)

	router.HandleFunc("/static", http.NotFound)
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("/", app.home)

	return router
}
