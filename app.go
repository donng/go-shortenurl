package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type App struct {
	Router     *chi.Mux
	Middleware *Middleware
}

type RequestParams struct {
	Url    string
	Expire int
}

func (a *App) Init() {
	a.Router = chi.NewRouter()
	a.Middleware = &Middleware{}
	a.InitRoutes()
}

func (a *App) InitRoutes() {
	a.Router.Use(a.Middleware.LoggingHandler, a.Middleware.RecoverHandler)

	a.Router.Post("/api/shorten_url", a.createShortUrl)
	a.Router.Get("/api/shorten_url_info", getShortUrlInfo)
	a.Router.Get("/{shorten_url:[a-zA-Z0-9]{1,11}}", redirect)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) createShortUrl(w http.ResponseWriter, r *http.Request) {
	var params RequestParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("parse params error: %v", err),
		})
		return
	}
	defer r.Body.Close()

	w.Write([]byte(fmt.Sprintf("hi %v, %v", params.Url, params.Expire)))
}

func getShortUrlInfo(w http.ResponseWriter, r *http.Request) {

}

func redirect(w http.ResponseWriter, r *http.Request) {

}
