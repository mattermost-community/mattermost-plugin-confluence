package platform

import (
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

func SaveSubscription(subscription serializer.Subscription, keySubscriptions map[string][]string, channelSubscriptions map[string]serializer.Subscription, key string) error {
	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := util.Set(key, keySubscriptions); err != nil {
		return err
	}
	if err := util.Set(subscription.ChannelID, channelSubscriptions); err != nil {
		return err
	}

	return nil
}
