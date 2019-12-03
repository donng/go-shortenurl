package main

import "net/http"

type Middleware struct {}

func (m Middleware) templateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
