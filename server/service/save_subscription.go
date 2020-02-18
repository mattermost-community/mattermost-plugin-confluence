package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "An error occurred attempting to save a subscription."
	aliasAlreadyExist       = "A subscription with the same alias already exists."
	urlSpaceKeyAlreadyExist = "A subscription with the same url and space key already exists."
)

func SavePageSubscription(subscription *serializer.PageSubscription) (int, error) {
	key := store.GetSubscriptionKey()
	err := store.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		fmt.Println("bytes=", initialBytes)
		subscriptions, err := serializer.SubscriptionsFromJson(initialBytes)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", subscriptions, err)
		if err != nil {
			return nil, err
		}
		fmt.Println("list=", subscriptions)
		subscription.Add(subscriptions)
		fmt.Println("list1=", subscriptions)
		modifiedBytes, marshalErr := json.Marshal(subscriptions)
		if marshalErr != nil {
			return nil, marshalErr
		}
		return modifiedBytes, nil
	})

	return http.StatusOK,err
}
