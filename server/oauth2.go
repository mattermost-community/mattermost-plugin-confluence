package main

import "net/http"

const (
	routeUserConnect  = "/oauth2/connect"
	routeUserComplete = "/oauth2/complete.html"
)

var userConnect = &Endpoint{
	Path:          routeUserConnect,
	Method:        http.MethodGet,
	Execute:       httpOAuth2Connect,
	RequiresAdmin: false,
}

var userConnectComplete = &Endpoint{
	Path:          routeUserComplete,
	Method:        http.MethodGet,
	Execute:       httpOAuth2Complete,
	RequiresAdmin: false,
}
