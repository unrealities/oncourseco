package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := routes()
	http.Handle("/", r)
}

func data(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("winning"))
}

func routes() http.Handler {
	router := httprouter.New()

	router.GET("/data", data)

	return router
}
