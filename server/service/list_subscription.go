package service

import (
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const getChannelSubscriptionsError = " Error getting channel subscriptions."

func GetChannelSubscriptions(channelID string) (map[string]serializer.Subscription, string, error) {
	key := store.GetChannelSubscriptionKey(channelID)
	channelSubscriptions := make(map[string]serializer.Subscription)
	if err := store.Get(key, &channelSubscriptions); err != nil {
		return nil, "", errors.New(getChannelSubscriptionsError)
	}
	return channelSubscriptions, key, nil
}

func GetURLSpaceKeyCombinationSubscriptions(baseURL, spaceKey string) (map[string][]string, string, error) {
	// Error is ignored as the url is already parsed in isValid method
	key, _ := store.GetURLSpaceKeyCombinationKey(baseURL, spaceKey)
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return nil, "", err
	}
	return keySubscriptions, key, nil
}
