package service

import (
	"encoding/json"
	"fmt"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	generalDeleteError   = "Error occurred while deleting subscription with alias **%s**."
	subscriptionNotFound = "Subscription with alias **%s** not found."
)

func DeleteSubscription(channelID, alias string) error {
	subs, gErr := GetSubscriptions()
	if gErr != nil {
		return fmt.Errorf(generalDeleteError, alias)
	}
	if channelSubscriptions, valid := subs.ByChannelID[channelID]; valid {
		if subscription, ok := channelSubscriptions[alias]; ok {
			aErr := store.AtomicModify(store.GetSubscriptionKey(), func(initialBytes []byte) ([]byte, error) {
				subscriptions, err := serializer.SubscriptionsFromJSON(initialBytes)
				if err != nil {
					return nil, err
				}
				subscription.Remove(subscriptions)
				modifiedBytes, marshalErr := json.Marshal(subscriptions)
				if marshalErr != nil {
					return nil, marshalErr
				}
				return modifiedBytes, nil
			})
			return aErr
		}
	}
	return fmt.Errorf(subscriptionNotFound, alias)
}
