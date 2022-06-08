package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "an error occurred attempting to save a subscription"
	aliasAlreadyExist       = "a subscription with the same name already exists in this channel"
	urlSpaceKeyAlreadyExist = "a subscription with the same url and space key already exists in this channel"
	urlPageIDAlreadyExist   = "a subscription with the same url and page id already exists in this channel"
)

func SaveSubscription(subscription serializer.Subscription) (int, error) {
	subs, err := GetSubscriptions()
	if err != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if err = subscription.ValidateSubscription(&subs); err != nil {
		return http.StatusBadRequest, err
	}
	key := store.GetSubscriptionKey()
	if err := store.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		subscriptions, err := serializer.SubscriptionsFromJSON(initialBytes)
		if err != nil {
			return nil, err
		}
		subscription.Add(subscriptions)
		modifiedBytes, marshalErr := json.Marshal(subscriptions)
		if marshalErr != nil {
			return nil, marshalErr
		}
		return modifiedBytes, nil
	}); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
