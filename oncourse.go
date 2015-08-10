package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := routes()
	http.Handle("/", r)
}

func Data(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("winning"))
}

func routes() http.Handler {
	router := httprouter.New()

	router.GET("/data", Data)
	router.GET("/SetOAuthCredentials", SetOAuthCredentials)

	return router
}

// Only to be used manually to update OAuth Credentials.
// Never store any keys, tokens or secrets in the code.
func SetOAuthCredentials(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)

	t := &Credentials{
		ClientId:     "",
		ClientSecret: ""}

	k := datastore.NewKey(c, "Credentials", "OAuth", 0, nil)

	_, err := datastore.Put(c, k, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}
