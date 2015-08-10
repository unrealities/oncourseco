package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/oncourse/oncourse/models"
	"google.golang.org/appengine"
)

func oncourse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := []models.Info{}

	c := appengine.NewContext(r)

	w.Header().Set("Content-Type", "application/json")

	js, err := JSONMarshal(info, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
