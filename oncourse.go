package main

import (
	"net/http"

	"github.com/oncourse/oncourse/routers"
)

func init() {
	r := routers.Routes()
	http.Handle("/", r)
}
