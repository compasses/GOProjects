package main

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc httprouter.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"ATS Check",
		"POST",
		"/sbo/service/EShopService@getATS",
		ATS,
	},
}
