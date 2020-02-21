package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

var confluenceCloudWebhook = &Endpoint{
	RequiresAuth: false,
	Path:         "/cloud/{event:[A-Za-z0-9_]+}",
	Method:       http.MethodPost,
	Execute:      handleConfluenceCloudWebhook,
}

func handleConfluenceCloudWebhook(w http.ResponseWriter, r *http.Request) {
	config.Mattermost.LogInfo("Received confluence cloud event.")

	if status, err := verifyHTTPSecret(config.GetConfig().Secret, r.FormValue("secret")); err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	params := mux.Vars(r)
	event := serializer.ConfluenceCloudEventFromJSON(r.Body)
	go service.SendConfluenceNotifications(event, params["event"])

	w.Header().Set("Content-Type", "application/json")
	ReturnStatusOK(w)
}
