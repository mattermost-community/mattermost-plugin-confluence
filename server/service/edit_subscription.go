package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/pkg/errors"
)

const openEditSubscriptionModalWebsocketEvent = "open_edit_subscription_modal"

func EditSubscription(subscription serializer.Subscription) (int, error) {
	channelSubscriptions := make(map[string]serializer.Subscription)

	if err := store.Get(store.GetChannelSubscriptionKey(subscription.ChannelID), &channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}
	// Error is ignored as the url is already parsed in isValid method
	key, _ := store.GetURLSpaceKeyCombinationKey(subscription.BaseURL, subscription.SpaceKey)
	keySubscriptions := make(map[string][]string)
	if err := store.Get(key, &keySubscriptions); err != nil {
		return http.StatusBadRequest, errors.New(generalError)
	}

	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}
	if err := store.Set(store.GetChannelSubscriptionKey(subscription.ChannelID), channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalError)
	}

	return http.StatusOK, nil
}

func OpenSubscriptionEditModal(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	alias := args[0]
	if err := store.Get(store.GetChannelSubscriptionKey(context.ChannelId), &channelSubscriptions); err != nil {
		util.PostCommandResponse(context, fmt.Sprintf("Error occured while editing subscription with alias **%s**.", alias))
		return &model.CommandResponse{}
	}
	if subscription, ok := channelSubscriptions[alias]; ok {
		bytes, err := json.Marshal(subscription)
		if err != nil {
			util.PostCommandResponse(context, fmt.Sprintf("Error occured while editing subscription with alias **%s**.", alias))
			return &model.CommandResponse{}
		}
		config.Mattermost.PublishWebSocketEvent(
			openEditSubscriptionModalWebsocketEvent,
			map[string]interface{}{
				"subscription": string(bytes),
			},
			&model.WebsocketBroadcast{
				UserId: context.UserId,
			},
		)
		return &model.CommandResponse{}
	}
	util.PostCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** not found.", alias))
	return &model.CommandResponse{}
}
