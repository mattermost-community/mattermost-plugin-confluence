package service

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/pkg/errors"
)

const (
	generalError            = "Failed to save subscription"
	aliasAlreadyExist       = "Subscription already exist in the channel with this alias"
	urlSpaceKeyAlreadyExist = "Subscription already exist in the channel with url and space key combination"
)

func SaveSubscription(subscription serializer.Subscription) (int, error) {
	channelSubscriptions := make(map[string]serializer.Subscription)

	if err := store.Get(store.GetChannelSubscriptionKey(subscription.ChannelID), &channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}
	if _, ok := channelSubscriptions[subscription.Alias]; ok {
		return http.StatusBadRequest, errors.New(aliasAlreadyExist)
	}
	// Error is ignored as the url is already parsed in isValid method
	key, _ := store.GetURLSpaceKeyCombinationKey(subscription.BaseURL, subscription.SpaceKey)
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return http.StatusBadRequest, errors.New(generalError)
	}
	if _, ok := keySubscriptions[subscription.ChannelID]; ok {
		return http.StatusBadRequest, errors.New(urlSpaceKeyAlreadyExist)
	}

	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}
	if err := store.Set(store.GetChannelSubscriptionKey(subscription.ChannelID), channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}

	return http.StatusOK, nil
}
