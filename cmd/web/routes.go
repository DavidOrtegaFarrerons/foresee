package main

import (
	"foresee/internal/web"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	baseChain := web.Chain{app.sessionManager.LoadAndSave, app.logRequest, app.authenticate, app.noSurf}
	authChain := append(baseChain, app.requiresAuthentication)

	fileServer := http.FileServer(
		web.NeuteredFileSystem(http.Dir("./ui/static")),
	)

	router.Handle("GET /favicon.ico", http.NotFoundHandler())

	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.Handle("GET /", http.HandlerFunc(app.home))

	router.Handle("GET /signup", http.HandlerFunc(app.signup))
	router.Handle("POST /signup", http.HandlerFunc(app.signupPost))

	router.Handle("GET /login", http.HandlerFunc(app.login))
	router.Handle("POST /login", http.HandlerFunc(app.loginPost))

	router.Handle("GET /account", http.HandlerFunc(app.account))

	router.Handle("GET /markets/create", authChain.ThenFunc(app.createMarket))
	router.Handle("POST /markets", authChain.ThenFunc(app.createMarketPost))
	router.Handle("GET /markets/{id}", http.HandlerFunc(app.viewMarket))
	router.Handle("POST /markets/{id}/bets", authChain.ThenFunc(app.createBetPost))
	router.Handle("GET /markets/{id}/resolve", authChain.ThenFunc(app.resolveMarket))
	router.Handle("POST /markets/{id}/resolve", authChain.ThenFunc(app.resolveMarketPost))

	router.Handle("POST /users/me/daily-claim", authChain.ThenFunc(app.dailyClaimPost))

	return baseChain.Then(router)
}
