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

const subscriptionEditSuccess = "Your subscription has been edited successfully."

func (p *Plugin) handleEditChannelSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params[ParamChannelID]
	subscriptionType := params[ParamSubscriptionType]
	oldSubscriptionType := params["oldSubscriptionType"]
	userID := r.Header.Get(config.HeaderMattermostUserID)

	instance, err := p.getInstanceFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

		newSubscription := subscription.(*serializer.SpaceSubscription).GetSubscription().UpdateSpaceIDAndUserID(strconv.FormatInt(resp.ID, 10), userID)
		subscription = newSubscription.GetSubscription()
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

		newSubscription := subscription.(*serializer.PageSubscription).GetSubscription().UpdateUserID(userID)
		subscription = newSubscription.GetSubscription()
	}

	var oldSubscription serializer.Subscription
	if oldSubscriptionType == serializer.SubscriptionTypeSpace {
		oldSubscription, err = serializer.OldSpaceSubscriptionFromJSON(body)
		if err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Error decoding request body for old space subscription.", err)
			return
		}
	} else if subscriptionType == serializer.SubscriptionTypePage {
		oldSubscription, err = serializer.OldPageSubscriptionFromJSON(body)
		if err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Error decoding request body for old page subscription.", err)
			return
		}
	}

	if err := p.DeleteSubscription(subscription.GetChannelID(), oldSubscription.GetAlias(), userID); err != nil {
		p.LogAndRespondError(w, http.StatusBadRequest, "not able to delete subscription.", err)
		return
	}

	if instance.Common().Type == ServerInstanceType {
		if err := p.CreateWebhook(instance, subscription, userID); err != nil {
			p.LogAndRespondError(w, http.StatusBadRequest, "Not able to create webhook.", err)
			return
		}
	}

	if err := service.EditSubscription(subscription); err != nil {
		p.LogAndRespondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   subscriptionEditSuccess,
	}
	_ = config.Mattermost.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}
