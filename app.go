package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

type App struct {
	Router *chi.Mux
	Middleware Middleware
}

func (a *App) Init() {
	a.Router = chi.NewRouter()
	a.Middleware = Middleware{}
	a.InitRoutes()
}

func (a *App) InitRoutes() {
	a.Router.Use(a.Middleware.templateMiddleware)

	a.Router.Post("/api/get_shortenurl", createShortUrl)
	a.Router.Get("/api/get_shortenurl_info", getShortUrlInfo)
	a.Router.Get("/{shorturl:[a-zA-Z0-9]+}", redirect)
}

func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}

func createShortUrl(w http.ResponseWriter, r *http.Request) {

}

func getShortUrlInfo(w http.ResponseWriter, r *http.Request) {

}

func redirect(w http.ResponseWriter, r *http.Request) {

}