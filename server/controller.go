package main

import (
	"crypto/subtle"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/mux"
	model "github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
)

type Endpoint struct {
	Path          string
	Method        string
	Execute       func(w http.ResponseWriter, r *http.Request)
	RequiresAdmin bool
}

// Endpoints is a map of endpoint key to endpoint object
// Usage: getEndpointKey(GetMetadata): GetMetadata
var Endpoints = map[string]*Endpoint{
	getEndpointKey(atlassianConnectJSON):    atlassianConnectJSON,
	getEndpointKey(confluenceCloudWebhook):  confluenceCloudWebhook,
	getEndpointKey(saveChannelSubscription): saveChannelSubscription,
	getEndpointKey(editChannelSubscription): editChannelSubscription,
	getEndpointKey(confluenceServerWebhook): confluenceServerWebhook,
	getEndpointKey(getChannelSubscription):  getChannelSubscription,

	getEndpointKey(autocompleteGetChannelSubscriptions): autocompleteGetChannelSubscriptions,
}

// Uniquely identifies an endpoint using path and method
func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	handleStaticFiles(r)

	s := r.PathPrefix("/api/v1").Subrouter()
	for _, endpoint := range Endpoints {
		handler := endpoint.Execute
		if endpoint.RequiresAdmin {
			handler = handleAdminRequired(endpoint)
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

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

func handleAdminRequired(endpoint *Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if IsAdmin(w, r) {
			endpoint.Execute(w, r)
		}
	}
}

// IsAdmin verifies if provided request is performed by a logged-in Mattermost user.
func IsAdmin(w http.ResponseWriter, r *http.Request) bool {
	userID := r.Header.Get(config.HeaderMattermostUserID)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return util.IsSystemAdmin(userID)
}

func ReturnStatusOK(w io.Writer) {
	m := make(map[string]string)
	m[model.STATUS] = model.StatusOk
	_, _ = w.Write([]byte(model.MapToJSON(m)))
}

func verifyHTTPSecret(expected, got string) (status int, err error) {
	for {
		if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1 {
			break
		}

		unescaped, _ := url.QueryUnescape(got)
		if unescaped == got {
			return http.StatusForbidden, errors.New("request URL: secret did not match")
		}
		got = unescaped
	}

	return 0, nil
}
