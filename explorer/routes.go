package main

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  httprouter.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"GetLinkchainOverview",
		"GET",
		"/api/v1/explorer/linkchian/overview",
		GetLinkchainOverview,
	},
	Route{
		"GetNodes",
		"GET",
		"/api/v1/explorer/linkchian/nodes",
		GetNodes,
	},
	Route{
		"GetAuthNodeDetail",
		"GET",
		"/api/v1/explorer/linkchian/authnodes/:node",
		GetAuthNodeDetail,
	},
	Route{
		"GetScanNodes",
		"GET",
		"/api/v1/explorer/linkchian/scannodes",
		GetScanNodes,
	},
}
