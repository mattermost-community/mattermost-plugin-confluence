package platform

import (
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/pkg/errors"
	"net/http"
)

func SaveSubscription(subscription serializer.Subscription) (error, int) {
	channelSubscriptions := make(map[string]serializer.Subscription)

	if err := store.Get(util.GetChannelSubscriptionKey(subscription.ChannelID), &channelSubscriptions); err != nil {
		return errors.New("Failed to save subscription"), http.StatusInternalServerError
	}
	if _, ok := channelSubscriptions[subscription.Alias]; ok {
		return errors.New("Subscription already exist in the channel with this alias."), http.StatusBadRequest
	}
	// Error is ignored as the url is already parsed in isValid method
	key, _ := util.GetURLSpaceKeyCombinationKey(subscription.BaseURL, subscription.SpaceKey)
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return errors.New("Failed to save subscription"), http.StatusBadRequest
	}
	if _, ok := keySubscriptions[subscription.ChannelID]; ok {
		return errors.New("Subscription already exist in the channel with url and space key combination."), http.StatusBadRequest
	}
	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return errors.New("Failed to save subscription"), http.StatusInternalServerError
	}
	if err := store.Set(util.GetChannelSubscriptionKey(subscription.ChannelID), channelSubscriptions); err != nil {
		return errors.New("Failed to save subscription"), http.StatusInternalServerError
	}

	return errors.New("Failed to save subscription"), http.StatusOK
}
