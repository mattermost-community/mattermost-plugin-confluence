package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

type Endpoint struct {
	Path         string
	Method       string
	Execute      func(w http.ResponseWriter, r *http.Request)
	RequiresAuth bool
}

// Endpoints is a map of endpoint key to endpoint object
// Usage: getEndpointKey(GetMetadata): GetMetadata
var Endpoints = map[string]*Endpoint{}

// Uniquely identifies an endpoint using path and method
func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

// InitAPI initializes the REST API
func InitAPI() *mux.Router {
	r := mux.NewRouter()
	handleStaticFiles(r)

	s := r.PathPrefix("/api/v1").Subrouter()
	for _, endpoint := range Endpoints {
		handler := endpoint.Execute
		if endpoint.RequiresAuth {
			handler = handleAuth(endpoint)
		}
		s.HandleFunc(endpoint.Path, handler).Methods(endpoint.Method)
	}

	return r
}

// handleStaticFiles handles the static files under the assets directory.
func handleStaticFiles(r *mux.Router) {
	bundlePath, err := config.Mattermost.GetBundlePath()
	if err != nil {
		config.Mattermost.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve files under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

func handleAuth(endpoint *Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Authenticated(w, r) {
			endpoint.Execute(w, r)
		}
	}
}

// Authenticated verifies if provided request is performed by a logged-in Mattermost user.
func Authenticated(w http.ResponseWriter, r *http.Request) bool {
	userID := r.Header.Get(config.HeaderMattermostUserID)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
