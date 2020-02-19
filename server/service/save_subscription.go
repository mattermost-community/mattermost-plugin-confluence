package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "An error occurred attempting to save a subscription."
	aliasAlreadyExist       = "A subscription with the same alias already exists."
	urlSpaceKeyAlreadyExist = "A subscription with the same url and space key already exists."
	urlPageIdAlreadyExist   = "A subscription with the same url and page id already exists."
)

func SaveSubscription(subscription serializer.Subscription) error {
	key := store.GetSubscriptionKey()
	err := store.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		subscriptions, err := serializer.SubscriptionsFromJson(initialBytes)
		if err != nil {
			return nil, err
		}
		subscription.Add(subscriptions)
		modifiedBytes, marshalErr := json.Marshal(subscriptions)
		if marshalErr != nil {
			return nil, marshalErr
		}
		return modifiedBytes, nil
	})
	return err
}

func ValidatePageSubscription(s serializer.PageSubscription) (int, error) {
	subs, gErr := GetSubscriptions()
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if channelSubscriptions, valid := subs.ByChannelID[s.ChannelID]; valid {
		if _, ok := channelSubscriptions[s.Alias]; ok {
			return http.StatusBadRequest, errors.New(aliasAlreadyExist)
		}
	}
	key := store.GetURLPageIDCombinationKey(s.BaseURL, s.PageID)
	if urlSpaceKeySubscriptions, valid := subs.ByURLPagID[key]; valid {
		if _, ok := urlSpaceKeySubscriptions[s.ChannelID]; ok {
			return http.StatusBadRequest, errors.New(urlPageIdAlreadyExist)
		}
	}
	return http.StatusOK, nil
}

func ValidateSpaceSubscription(s serializer.SpaceSubscription) (int, error) {
	subs, gErr := GetSubscriptions()
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if channelSubscriptions, valid := subs.ByChannelID[s.ChannelID]; valid {
		if _, ok := channelSubscriptions[s.Alias]; ok {
			return http.StatusBadRequest, errors.New(aliasAlreadyExist)
		}
	}
	key := store.GetURLPageIDCombinationKey(s.BaseURL, s.SpaceKey)
	if urlSpaceKeySubscriptions, valid := subs.ByURLSpaceKey[key]; valid {
		if _, ok := urlSpaceKeySubscriptions[s.ChannelID]; ok {
			return http.StatusBadRequest, errors.New(urlSpaceKeyAlreadyExist)
		}
	}
	return http.StatusOK, nil
}
