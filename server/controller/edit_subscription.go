package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

const (
	specifyAlias         = "Please specify alias."
	subscriptionNotFound = "Subscription with alias **%s** not found."
	generalError         = "Some error occurred. Please try again after sometime."
)

var (
	EditChannelSubscription = &Endpoint{
		RequiresAuth: true,
		Path:         "/subscription",
		Method:       http.MethodPut,
		Execute:      editChannelSubscription,
	}
	OpenEditSubscriptionModal = &Endpoint{
		RequiresAuth: true,
		Path:         "/open-edit-subscription",
		Method:       http.MethodPost,
		Execute:      openEditSubscriptionModal,
	}
)

func editChannelSubscription(w http.ResponseWriter, r *http.Request) {
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

func openEditSubscriptionModal(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	data := serializer.EditSubscription{}
	if err := body.Decode(&data); err != nil {
		config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}
	args, aErr := util.SplitArgs(data.Message)
	if aErr != nil {
		http.Error(w, aErr.Error(), http.StatusBadRequest)
		return
	}
	if len(args) < 3 {
		http.Error(w, specifyAlias, http.StatusBadRequest)
		return
	}
	channelSubscriptions, _, gErr := service.GetChannelSubscriptions(data.ChannelID)
	if gErr != nil {
		http.Error(w, generalError, http.StatusInternalServerError)
		return
	}
	subscription, found := channelSubscriptions[args[2]]
	if found {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(subscription.ToJSON()))
		return
	}
	http.Error(w, fmt.Sprintf(subscriptionNotFound, args[2]), http.StatusInternalServerError)
}
