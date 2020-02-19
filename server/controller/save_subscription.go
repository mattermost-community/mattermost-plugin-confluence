package controller

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

const subscriptionSaveSuccess = "Your subscription has been saved."

var saveChannelSubscription = &Endpoint{
	RequiresAuth: true,
	Path:         "/subscription",
	Method:       http.MethodPost,
	Execute:      handleSaveSubscription,
}

func handleSaveSubscription(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	subscriptionType := r.FormValue("type")
	userID := r.Header.Get(config.HeaderMattermostUserID)
	if subscriptionType == serializer.SubscriptionTypeSpace {
		subscription := serializer.SpaceSubscription{}
		if err := body.Decode(&subscription); err != nil {
			config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
			http.Error(w, "Could not decode request body", http.StatusBadRequest)
			return
		}
		if errCode, err := saveSubscription(subscription, subscription.ChannelID, userID); err != nil {
			config.Mattermost.LogError(err.Error())
			http.Error(w, err.Error(), errCode)
			return
		}
	} else if subscriptionType == serializer.SubscriptionTypePage {
		subscription := serializer.PageSubscription{}
		if err := body.Decode(&subscription); err != nil {
			config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
			http.Error(w, "Could not decode request body", http.StatusBadRequest)
			return
		}
		if errCode, err := saveSubscription(subscription, subscription.ChannelID, userID); err != nil {
			config.Mattermost.LogError(err.Error())
			http.Error(w, err.Error(), errCode)
			return
		}
	}
	ReturnStatusOK(w)
}

func saveSubscription(subscription serializer.Subscription, channelID, userID string) (int, error) {
	if err := subscription.IsValid(); err != nil {
		return http.StatusBadRequest, err
	}

	if subscription.Name() == serializer.SubscriptionTypeSpace {
		if errCode, err := service.ValidateSpaceSubscription(subscription.(serializer.SpaceSubscription)); err != nil {
			return errCode, err
		}
	} else if subscription.Name() == serializer.SubscriptionTypePage {
		if errCode, err := service.ValidatePageSubscription(subscription.(serializer.PageSubscription)); err != nil {
			return errCode, err
		}
	}

	if err := service.SaveSubscription(subscription); err != nil {
		return http.StatusInternalServerError, err
	}
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: channelID,
		Message:   subscriptionSaveSuccess,
	}
	_ = config.Mattermost.SendEphemeralPost(userID, post)

	return http.StatusOK, nil
}
