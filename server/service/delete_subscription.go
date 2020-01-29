package service

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-server/model"
)

func DeleteSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	alias := args[0]
	if err := store.Get(store.GetChannelSubscriptionKey(context.ChannelId), &channelSubscriptions); err != nil {
		util.PostCommandResponse(context, fmt.Sprintf("Error occured while deleting subscription with alias **%s**.", alias))
		return &model.CommandResponse{}
	}
	if subscription, ok := channelSubscriptions[alias]; ok {
		if err := deleteSubscriptionUtil(subscription, channelSubscriptions, alias); err != nil {
			util.PostCommandResponse(context, fmt.Sprintf("Error occured while deleting subscription with alias **%s**.", alias))
			return &model.CommandResponse{}
		}
		util.PostCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** deleted successfully.", alias))
		return &model.CommandResponse{}
	}
	util.PostCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** not found.", alias))
	return &model.CommandResponse{}
}

func deleteSubscriptionUtil(subscription serializer.Subscription, channelSubscriptions map[string]serializer.Subscription, alias string) error {
	key, kErr := store.GetURLSpaceKeyCombinationKey(subscription.BaseURL, subscription.SpaceKey)
	if kErr != nil {
		return kErr
	}
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return err
	}
	delete(keySubscriptions, subscription.ChannelID)
	delete(channelSubscriptions, alias)
	if err := store.Set(key, keySubscriptions); err != nil {
		return err
	}
	if err := store.Set(store.GetChannelSubscriptionKey(subscription.ChannelID), channelSubscriptions); err != nil {
		return err
	}
	return nil
}
