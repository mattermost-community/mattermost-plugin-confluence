package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	openEditSubscriptionModalWebsocketEvent = "open_edit_subscription_modal"
	generalError                            = "Some error occurred. Please try again after sometime."
	subscriptionEditSuccess                 = "Subscription updated successfully."
)

func EditSubscription(subscription serializer.Subscription, userID string) (int, error) {
	channelSubscriptions, cKey, gErr := GetChannelSubscriptions(subscription.ChannelID)
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	keySubscriptions, key, kErr := GetURLSpaceKeyCombinationSubscriptions(subscription.BaseURL, subscription.SpaceKey)
	if kErr != nil {
		return http.StatusInternalServerError, kErr
	}

	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if err := store.Set(cKey, channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: subscription.ChannelID,
		Message:   subscriptionEditSuccess,
	}
	_ = config.Mattermost.SendEphemeralPost(userID, post)

	return http.StatusOK, nil
}

func OpenSubscriptionEditModal(channelID, userID, alias string) error {
	channelSubscriptions, _, gErr := GetChannelSubscriptions(channelID)
	if gErr != nil {
		return errors.New(generalError)
	}
	if subscription, ok := channelSubscriptions[alias]; ok {
		bytes, err := json.Marshal(subscription)
		if err != nil {
			return errors.New(generalError)
		}
		config.Mattermost.PublishWebSocketEvent(
			openEditSubscriptionModalWebsocketEvent,
			map[string]interface{}{
				"subscription": string(bytes),
			},
			&model.WebsocketBroadcast{
				UserId: userID,
			},
		)
		return nil
	}

	return errors.New(fmt.Sprintf(subscriptionNotFound, alias))
}
