package platform

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/pkg/errors"
)

func SaveSubscription(subscription serializer.Subscription) (int, error) {
	channelSubscriptions := make(map[string]serializer.Subscription)

	if err := store.Get(util.GetChannelSubscriptionKey(subscription.ChannelID), &channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New("failed to save subscription")
	}
	if _, ok := channelSubscriptions[subscription.Alias]; ok {
		return http.StatusBadRequest, errors.New("subscription already exist in the channel with this alias")
	}
	// Error is ignored as the url is already parsed in isValid method
	key, _ := util.GetURLSpaceKeyCombinationKey(subscription.BaseURL, subscription.SpaceKey)
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return http.StatusBadRequest, errors.New("failed to save subscription")
	}
	if _, ok := keySubscriptions[subscription.ChannelID]; ok {
		return http.StatusBadRequest, errors.New("subscription already exist in the channel with url and space key combination")
	}
	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New("failed to save subscription")
	}
	if err := store.Set(util.GetChannelSubscriptionKey(subscription.ChannelID), channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New("failed to save subscription")
	}

	return http.StatusOK, errors.New("failed to save subscription")
}
