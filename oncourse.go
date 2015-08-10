package main

import (
	"net/http"

	"github.com/unrealities/oncourseco/routers"
)

func init() {
	r := routers.Routes()
	http.Handle("/", r)
}
