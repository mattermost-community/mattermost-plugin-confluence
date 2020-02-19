package service

import (
	"encoding/json"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "An error occurred attempting to save a subscription."
	aliasAlreadyExist       = "A subscription with the same alias already exists."
	urlSpaceKeyAlreadyExist = "A subscription with the same url and space key already exists."
)

func SaveSubscription(subscription serializer.Subscription) (int, error) {
	key := store.GetSubscriptionKey()
	err := store.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		subscriptions, err := serializer.SubscriptionsFromJson(initialBytes)
		if err != nil {
			return nil, err
		}
		subscription.Add(subscriptions)
		modifiedBytes, marshalErr := json.Marshal(*subscriptions)
		if marshalErr != nil {
			return nil, marshalErr
		}
		return modifiedBytes, nil
	})

	return http.StatusOK,err
}
