package routers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes() http.Handler {
	router := httprouter.New()

	return router
}
