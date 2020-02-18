package service

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const getChannelSubscriptionsError = " Error getting channel subscriptions."

func GetSubscriptions() (serializer.Subscriptions, error) {
	key := store.GetSubscriptionKey()
	var subscriptions serializer.Subscriptions
	if err := store.Get(key, &subscriptions); err != nil {
		return serializer.Subscriptions{}, errors.New(getChannelSubscriptionsError)
	}
	return subscriptions, nil
}

func GetSubscriptionsByChannelID(channelID string) (map[string]serializer.Subscription, error) {
	subscriptions, err := GetSubscriptions()
	fmt.Println("abc=", subscriptions)
	if err != nil {
		return nil, err
	}
	return subscriptions.ByChannelID[channelID], nil
}

func GetSubscriptionsByURLSpaceKey(url, spaceKey string) (map[string][]string, error) {
	key := store.GetURLSpaceKeyCombinationKey(url, spaceKey)
	subscriptions, err := GetSubscriptions()
	if err != nil {
		return nil, err
	}
	return subscriptions.ByURLSpaceKey[key], nil
}

func GetSubscriptionsByURLPageID(url, pageID string) (map[string][]string, error) {
	key := store.GetURLPageIDCombinationKey(url, pageID)
	subscriptions, err := GetSubscriptions()
	if err != nil {
		return nil, err
	}
	return subscriptions.ByURLPagID[key], nil
}
