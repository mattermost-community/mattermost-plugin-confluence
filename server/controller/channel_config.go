package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/platform"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

var SaveChannelSubscription = &Endpoint{
	RequiresAuth: true,
	Path:         "/save-channel-subscription",
	Execute:      saveChannelSubscription,
}

func saveChannelSubscription(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	subscription := serializer.Subscription{}
	if err := body.Decode(&subscription); err != nil {
		config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}
	if err := platform.SaveSubscription(subscription); err != nil {
		config.Mattermost.LogError("Unable to save subscription", "Error", err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
