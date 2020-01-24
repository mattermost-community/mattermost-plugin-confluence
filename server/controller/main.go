package controller

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
)

type Endpoint struct {
	Path         string
	Execute      func(w http.ResponseWriter, r *http.Request)
	RequiresAuth bool
}

var Endpoints = map[string]*Endpoint{
	SaveChannelSubscription.Path: SaveChannelSubscription,
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
