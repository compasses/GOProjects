package main

import (
	"github.com/julienschmidt/httprouter"
)

func ServerRouter() *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		httpHandle := Logger(route.HandleFunc, route.Name)

		router.Handle(
			route.Method,
			route.Pattern,
			httpHandle,
		)
	}

	router.NotFound = LoggerNotFound(NotFoundHandler)

	return router
}
