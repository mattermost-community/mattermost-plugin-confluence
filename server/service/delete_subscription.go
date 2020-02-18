package service

import (
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const (
	generalDeleteError   = "Error occurred while deleting subscription with alias **%s**."
	subscriptionNotFound = "Subscription with alias **%s** not found."
)

func DeleteSubscription(channelID, alias string) error {
	// channelSubscriptions, cKey, gErr := GetChannelSubscriptions(channelID)
	// if gErr != nil {
	// 	return errors.New(fmt.Sprintf(generalDeleteError, alias))
	// }
	// if subscription, ok := channelSubscriptions[alias]; ok {
	// 	if err := deleteSubscriptionUtil(subscription, channelSubscriptions, cKey, alias); err != nil {
	// 		return errors.New(fmt.Sprintf(generalDeleteError, alias))
	// 	}
	// 	return nil
	// }
	// return errors.New(fmt.Sprintf(subscriptionNotFound, alias))
	return nil
}

func deleteSubscriptionUtil(subscription serializer.Subscription, channelSubscriptions map[string]serializer.Subscription, cKey, alias string) error {
	// keySubscriptions, key, gErr := GetURLSpaceKeyCombinationSubscriptions(subscription.BaseURL, subscription.SpaceKey)
	// if gErr != nil {
	// 	return gErr
	// }
	// delete(keySubscriptions, subscription.ChannelID)
	// delete(channelSubscriptions, alias)
	// if err := store.Set(key, keySubscriptions); err != nil {
	// 	return err
	// }
	// if err := store.Set(cKey, channelSubscriptions); err != nil {
	// 	return err
	// }
	return nil
}
