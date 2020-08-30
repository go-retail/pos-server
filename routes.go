package main

import (
	"net/http"
)

//Route ..
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}



var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		home,
	},
	Route{
		"Create Transactions",
		"POST",
		"/txns/{storeID}/{posID}/create",
		createTxn,
	},
}
