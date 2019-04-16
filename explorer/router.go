package main

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {

	router := httprouter.New()

	for _, route := range routes {

		router.Handle(route.Method, route.Pattern, route.Handle)
	}

	return router
}
