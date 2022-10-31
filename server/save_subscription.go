package main

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils/types"
)

const (
	subscriptionSaveSuccess = "Your subscription has been saved."
	ParamChannelID          = "channelID"
	ParamSubscriptionType   = "type"
)

func (p *Plugin) handleSaveSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params[ParamChannelID]
	subscriptionType := params[ParamSubscriptionType]
	userID := r.Header.Get(config.HeaderMattermostUserID)

	instance, err := p.getInstanceFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.HasPermissionToManageSubscription(instance.GetURL(), userID, channelID)
	if err != nil {
		p.LogAndRespondError(w, http.StatusForbidden, "Don't have the permission to create subscription. Please contact your administrator", err)
		return
	}

	conn, err := p.userStore.LoadConnection(types.ID(instance.GetURL()), types.ID(userID))
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "Error in loading connection.", err)
		return
	}

	client, err := instance.GetClient(conn)
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "Not able to get Client.", err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var subscription serializer.Subscription
	if subscriptionType == serializer.SubscriptionTypeSpace {
		subscription, err = serializer.SpaceSubscriptionFromJSON(body)
		if err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Error decoding request body for space subscription.", err)
			return
		}

		spaceKey := subscription.(*serializer.SpaceSubscription).GetSubscription().SpaceKey
		resp, err := client.GetSpaceData(spaceKey)
		if err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Error getting space related data for space subscription.", err)
			return
		}

		updatedSubscrption := subscription.(*serializer.SpaceSubscription).GetSubscription().UpdateSpaceIDAndUserID(strconv.FormatInt(resp.ID, 10), userID)
		subscription = updatedSubscrption.GetSubscription()
	} else if subscriptionType == serializer.SubscriptionTypePage {
		subscription, err = serializer.PageSubscriptionFromJSON(body)
		if err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Error decoding request body for page subscription.", err)
			return
		}

		pageID, err := strconv.Atoi(subscription.(*serializer.PageSubscription).GetSubscription().PageID)
		if err != nil {
			p.LogAndRespondError(w, http.StatusInternalServerError, "Error converting pageID to integer.", err)
			return
		}

		_, err = client.GetPageData(pageID)
		if err != nil {
			p.LogAndRespondError(w, http.StatusInternalServerError, "Error getting page related data for page subscription.", err)
			return
		}

		updatedSubscrption := subscription.(*serializer.PageSubscription).UpdateUserID(userID)
		subscription = updatedSubscrption.GetSubscription()
	}

	if instance.Common().Type == ServerInstanceType {
		if err := p.CreateWebhook(instance, subscription, userID); err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Not able to create webhook.", err)
			return
		}
	}

	if statusCode, sErr := service.SaveSubscription(subscription); sErr != nil {
		p.LogAndRespondError(w, statusCode, sErr.Error(), sErr)
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   subscriptionSaveSuccess,
	}
	_ = p.API.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}

func (p *Plugin) LogAndRespondError(w http.ResponseWriter, statusCode int, errorLog string, err error) {
	p.API.LogError(errorLog, "Error", err.Error())
	http.Error(w, errorLog, statusCode)
}
