package platform

import (
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

func SaveSubscription(subscription serializer.Subscription) error {
	key, uErr := util.GetKey(subscription.BaseURL, subscription.SpaceKey)
	if uErr != nil {
		return uErr
	}
	keySubscriptions := make(map[string][]string)
	if err := util.Get(key, &keySubscriptions); err != nil {
		return err
	}
	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions := make(map[string]serializer.Subscription)
	if err := util.Get(subscription.ChannelID, &channelSubscriptions); err != nil {
		return err
	}
	channelSubscriptions[subscription.Alias] = subscription
	if err := util.Set(key, keySubscriptions); err != nil {
		return err
	}
	if err := util.Set(subscription.ChannelID, channelSubscriptions); err != nil {
		return err
	}

	return nil
}
