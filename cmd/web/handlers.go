package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.home)

	return router
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the homepage"))
}
