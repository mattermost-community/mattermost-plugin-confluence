package service

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const generalError = "Some error occurred. Please try again after sometime."

func GetChannelSubscription(channelID, alias string) (serializer.Subscription, error, int){
	channelSubscriptions, _, gErr := GetChannelSubscriptions(channelID)
	if gErr != nil {
		return serializer.Subscription{}, errors.New(generalError), http.StatusInternalServerError
	}
	subscription, found := channelSubscriptions[alias]
	if !found {
		return serializer.Subscription{}, errors.New(fmt.Sprintf(subscriptionNotFound, alias)), http.StatusBadRequest
	}
	return subscription, nil, http.StatusOK
}
