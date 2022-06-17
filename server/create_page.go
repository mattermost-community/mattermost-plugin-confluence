package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils/types"
)

const (
	pageCreateSuccess = "You created a page [%s](%s) in space [%s](%s)"
)

func (p *Plugin) handleCreatePage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	spaceKey := params["spaceKey"]
	channelID := params["channelID"]
	userID := r.Header.Get(config.HeaderMattermostUserID)

	instance, err := p.getInstanceFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := p.userStore.LoadConnection(types.ID(instance.GetURL()), types.ID(userID))
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "error in loading connection.", err)
		return
	}

	client, err := instance.GetClient(conn)
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "not able to get Client.", err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageDetails, err := serializer.PageDetailsFromJSON(body)
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "error decoding request body for page details.", err)
		return
	}

	createPageResponse, err := client.CreatePage(spaceKey, pageDetails)
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "not able to create page.", err)
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf(pageCreateSuccess, pageDetails.Title, fmt.Sprintf("%s%s", createPageResponse.Links.BaseURL, createPageResponse.Links.Self), spaceKey, fmt.Sprintf("%s%s", createPageResponse.Links.BaseURL, createPageResponse.Space.Links.Self)),
	}
	_ = p.API.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}
