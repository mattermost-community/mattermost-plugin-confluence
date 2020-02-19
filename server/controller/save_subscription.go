package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

const subscriptionSaveSuccess = "Your subscription has been saved."

var saveChannelSubscription = &Endpoint{
	RequiresAuth: true,
	Path:         "/subscription",
	Method:       http.MethodPost,
	Execute:      handleSavePageSubscription,
}

func handleSavePageSubscription(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	subscription := serializer.PageSubscription{}
	if err := body.Decode(&subscription); err != nil {
		config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}
	fmt.Println("sub=", subscription)
	if err := subscription.IsValid(); err != nil {
		config.Mattermost.LogError(err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errCode, err := service.SaveSubscription(&subscription); err != nil {
		config.Mattermost.LogError(err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, err.Error(), errCode)
		return
	}

	userID := r.Header.Get(config.HeaderMattermostUserID)
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: subscription.ChannelID,
		Message:   subscriptionSaveSuccess,
	}
	_ = config.Mattermost.SendEphemeralPost(userID, post)

	ReturnStatusOK(w)
}
