package service

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const subscriptionEditSuccess = "Subscription updated successfully."

func EditSubscription(subscription serializer.Subscription, userID string) (int, error) {
	channelSubscriptions, cKey, gErr := GetChannelSubscriptions(subscription.ChannelID)
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	keySubscriptions, key, keyErr := GetURLSpaceKeyCombinationSubscriptions(subscription.BaseURL, subscription.SpaceKey)
	if keyErr != nil {
		return http.StatusInternalServerError, keyErr
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
