package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
)

var confluenceServerWebhook = &Endpoint{
	Path:          "/server/webhook",
	Method:        http.MethodPost,
	Execute:       handleConfluenceServerWebhook,
	RequiresAdmin: false,
}

func handleConfluenceServerWebhook(w http.ResponseWriter, r *http.Request) {
	config.Mattermost.LogInfo("Received confluence server event.")

	if status, err := verifyHTTPSecret(config.GetConfig().Secret, r.FormValue("secret")); err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	event := serializer.ConfluenceServerEventFromJSON(r.Body)
	go service.SendConfluenceNotifications(event, event.Event)

	w.Header().Set("Content-Type", "application/json")
	ReturnStatusOK(w)
}
