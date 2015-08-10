package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

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

func GetOAuthCredentials(w http.ResponseWriter, r *http.Request) Credentials {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Credentials")
	tc := Credentials{}
	t := q.Run(c)
	for {
		_, err := t.Next(&tc)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	return tc
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}
