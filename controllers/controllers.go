package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func oncourse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("winning"))
}
