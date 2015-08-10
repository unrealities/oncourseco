package routers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrealities/oncourseco/controllers"
)

func Routes() http.Handler {
	router := httprouter.New()

	router.GET("/data", controllers.oncourse)

	return router
}
