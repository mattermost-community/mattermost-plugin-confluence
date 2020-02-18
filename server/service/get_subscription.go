package service

import (
	"net/http"
)

const generalError = "Some error occurred. Please try again after sometime."

func GetChannelSubscription(channelID, alias string) ( int, error) {
	// channelSubscriptions, _, gErr := GetChannelSubscriptions(channelID)
	// if gErr != nil {
	// 	return serializer.Subscription{}, http.StatusInternalServerError, errors.New(generalError)
	// }
	// subscription, found := channelSubscriptions[alias]
	// if !found {
	// 	return serializer.Subscription{}, http.StatusBadRequest, errors.New(fmt.Sprintf(subscriptionNotFound, alias))
	// }
	return  http.StatusOK, nil
}
