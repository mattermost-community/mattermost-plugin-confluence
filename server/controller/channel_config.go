package controller

import (
	"encoding/json"
	"net/http"
	url2 "net/url"

	"github.com/Brightscout/mattermost-plugin-confluence/server/util"

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
	channelSubscriptions := make(map[string]serializer.Subscription)
	if subscription.Alias == "" {
		http.Error(w, "Alias can not be empty", http.StatusBadRequest)
		return
	}
	if err := util.Get(subscription.ChannelID, &channelSubscriptions); err != nil {
		http.Error(w, "Failed to save subscription", http.StatusBadRequest)
		return
	}
	if _, ok := channelSubscriptions[subscription.Alias]; ok {
		http.Error(w, "Subscription already exist in the channel with this alias.", http.StatusBadRequest)
		return
	}
	if subscription.BaseURL == "" {
		http.Error(w, "Base url can not be empty.", http.StatusBadRequest)
		return
	}
	if _, err := url2.Parse(subscription.BaseURL); err != nil {
		http.Error(w, "Enter a valid url.", http.StatusBadRequest)
		return
	}
	if subscription.SpaceKey == "" {
		http.Error(w, "Space key can not be empty.", http.StatusBadRequest)
		return
	}
	key, _ := util.GetKey(subscription.BaseURL, subscription.SpaceKey)
	keySubscriptions := make(map[string][]string)
	if err := util.Get(key, &keySubscriptions); err != nil {
		http.Error(w, "Failed to save subscription", http.StatusBadRequest)
		return
	}
	if _, ok := keySubscriptions[subscription.ChannelID]; ok {
		http.Error(w, "Subscription already exist in the channel with url and space key combination.", http.StatusBadRequest)
		return
	}

	if err := platform.SaveSubscription(subscription, keySubscriptions, channelSubscriptions, key); err != nil {
		config.Mattermost.LogError("Failed to save subscription", "Error", err.Error(), "channelID", subscription.ChannelID)
		http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
		return
	}
}
