package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "An error occurred attempting to save a subscription."
	aliasAlreadyExist       = "A subscription with the same alias already exists."
	urlSpaceKeyAlreadyExist = "A subscription with the same url and space key already exists."
	urlPageIDAlreadyExist   = "A subscription with the same url and page id already exists."
)

func SaveSubscription(subscription serializer.Subscription) (int, error) {
	subs, gErr := GetSubscriptions()
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if vErr := subscription.ValidateSubscription(&subs); vErr != nil {
		return http.StatusBadRequest, vErr
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
