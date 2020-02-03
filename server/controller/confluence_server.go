package controller

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

var confluenceServerWebhook = &Endpoint{
	Path:         "/server/webhook",
	Method:       http.MethodPost,
	Execute:      handleConfluenceServerWebhook,
	RequiresAuth: false,
}

func handleConfluenceServerWebhook(w http.ResponseWriter, r *http.Request) {
	config.Mattermost.LogInfo("Received confluence server event.")

	event := serializer.ConfluenceServerEventFromJson(r.Body)
	service.SendConfluenceServerNotifications(event)

	w.Header().Set("Content-Type", "application/json")
	ReturnStatusOK(w)
}
