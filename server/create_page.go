package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
)

const (
	PageCreateSuccessMessage = "You created a page [%s](%s%s) in space [%s](%s%s)"
)

func (p *Plugin) handleCreatePage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	spaceKey := params["spaceKey"]
	channelID := params["channelID"]
	userID := r.Header.Get(config.HeaderMattermostUserID)

	client, err := p.GetClientFromURL(r.URL.Path, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	createPageResponse, statusCode, err := client.CreatePage(spaceKey, pageDetails)
	if err != nil {
		p.LogAndRespondError(w, statusCode, "not able to create page.", err)
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf(PageCreateSuccessMessage, pageDetails.Title, createPageResponse.Links.BaseURL, createPageResponse.Links.Self, spaceKey, createPageResponse.Links.BaseURL, createPageResponse.Space.Links.Self),
	}
	_ = p.API.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}
