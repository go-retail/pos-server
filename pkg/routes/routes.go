package routes

import (
	"net/http"

	"github.com/go-retail/pos-server/pkg/handlers"
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
		handlers.Home,
	},
	Route{
		"Create Transactions",
		"POST",
		"/txns/{storeID}/{posID}/create",
		handlers.CreateTxn,
	},
}
