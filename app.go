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
	Config     *Conf
	// storage engine interface, eg. redis,database
	Storage
}

type RequestParams struct {
	Url    string
	Expire int64
}

type ResponseParams struct {
	ShortLink string
}

func (a *App) Init() {
	a.Router = chi.NewRouter()
	a.Middleware = &Middleware{}
	a.Config = InitConfig()
	a.Storage = NewRedisClient(a.Config.Redis)
	a.InitRoutes()
}

func (a *App) InitRoutes() {
	a.Router.Use(a.Middleware.LoggingHandler)

	a.Router.Post("/api/shorten", a.createShortLink)
	a.Router.Get("/api/info", a.getShortLinkInfo)
	a.Router.Get("/{link:[a-zA-Z0-9]{1,11}}", a.redirect)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var params RequestParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("parse params error: %v", err),
		})
		return
	}
	defer r.Body.Close()

	sid, err := a.Storage.Shorten(params.Url, params.Expire)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJson(w, http.StatusOK, ResponseParams{ShortLink: sid})
	}
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	eid := r.Context().Value("link").(string)
	detail, err := a.Storage.ShortLinkInfo(eid)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJson(w, http.StatusOK, detail)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	eid := chi.URLParam(r, "link")
	url, err := a.Storage.UnShorten(eid)
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
