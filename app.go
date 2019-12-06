package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
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
	Url             string `json:"url" validate:"required"`
	ExpireInMinutes int64  `json:"expire_in_minutes" validate:"required"`
}

type ResponseParams struct {
	ShortLink string `json:"short_link"`
}

func (a *App) InitApp() {
	a.Router = chi.NewRouter()
	a.Middleware = &Middleware{}
	a.Config = InitConfig()
	a.Storage = NewRedisClient(a.Config.Redis)
	a.InitRoutes()
}

func (a *App) InitRoutes() {
	a.Router.Use(a.Middleware.LoggingHandler, a.Middleware.RecoverHandler)

	a.Router.Post("/api/shorten", a.createShortLink)
	a.Router.Get("/api/info/{link:[a-zA-Z0-9]{1,11}}", a.getShortLinkInfo)
	a.Router.Get("/{link:[a-zA-Z0-9]{1,11}}", a.redirect)
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", a.Config.Server.HttpPort), a.Router))
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var params RequestParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("json parse error: %v", params),
		})
		return
	}
	// validate request params
	validate := validator.New()
	if err := validate.Struct(&params); err != nil {
		respondWithError(w, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("json validate error: %v", params),
		})
		return
	}

	defer r.Body.Close()

	encodeId, err := a.Storage.Shorten(params.Url, params.ExpireInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJson(w, http.StatusOK, ResponseParams{ShortLink: encodeId})
	}
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	encodeId := chi.URLParam(r, "link")
	detail, err := a.Storage.ShortLinkInfo(encodeId)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJson(w, http.StatusOK, detail)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	encodeId := chi.URLParam(r, "link")
	url, err := a.Storage.UnShorten(encodeId)
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
