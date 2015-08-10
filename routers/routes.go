package routers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.FileServer(http.Dir("www/")).ServeHTTP

	return router
}
