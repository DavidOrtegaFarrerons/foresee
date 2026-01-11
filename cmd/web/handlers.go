package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.home)

	return router
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html", &templateData{})
}
