package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
)

const (
	ParamUserID = "userID"
)

func (p *Plugin) handleConfluenceConfig(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	channelID := pathParams[ParamChannelID]
	userID := pathParams[ParamUserID]
	if !model.IsValidId(channelID) || !model.IsValidId(userID) {
		http.Error(w, "Invalid channel or user id", http.StatusBadRequest)
		p.API.LogError("Invalid channel or user id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	submitRequest := &model.SubmitDialogRequest{}
	if err := decoder.Decode(&submitRequest); err != nil {
		p.API.LogError("Error decoding the submit dialog request.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := &serializer.ConfluenceConfig{
		ServerURL:    submitRequest.Submission[configServerURL].(string),
		ClientID:     submitRequest.Submission[configClientID].(string),
		ClientSecret: submitRequest.Submission[configClientSecret].(string),
	}

	if err := p.instanceStore.StoreInstanceConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.API.LogError("Error storing instance config.", "Error", err.Error())
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf("Your config is saved for confluence instance %s", config.ServerURL),
	}

	_ = p.API.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}
