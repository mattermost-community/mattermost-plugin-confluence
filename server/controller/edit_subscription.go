package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

var editChannelSubscription = &Endpoint{
	RequiresAuth: true,
	Path:         "/subscription",
	Method:       http.MethodPut,
	Execute:      handleEditChannelSubscription,
}


func handleEditChannelSubscription(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	subscription := serializer.Subscription{}
	if err := body.Decode(&subscription); err != nil {
		config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}

	if err := subscription.IsValid(); err != nil {
		config.Mattermost.LogError(err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Header.Get(config.HeaderMattermostUserID)
	if errCode, err := service.EditSubscription(subscription, userID); err != nil {
		config.Mattermost.LogError(err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, err.Error(), errCode)
		return
	}
	ReturnStatusOK(w)
}
