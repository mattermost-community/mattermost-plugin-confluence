package main

import "net/http"

const (
	routeUserConnect        = "/oauth2/connect"
	routeUserComplete       = "/oauth2/complete.html"
	routeUserConnectionInfo = "/user-connection-info"
)

var userConnect = &Endpoint{
	Path:    routeUserConnect,
	Method:  http.MethodGet,
	Execute: httpOAuth2Connect,
}

var userConnectComplete = &Endpoint{
	Path:    routeUserComplete,
	Method:  http.MethodGet,
	Execute: httpOAuth2Complete,
}

var userConnectionInfo = &Endpoint{
	Path:    routeUserConnectionInfo,
	Method:  http.MethodGet,
	Execute: httpGetUserInfo,
}
