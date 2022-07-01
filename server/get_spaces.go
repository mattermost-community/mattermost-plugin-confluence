package main

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
)

func (p *Plugin) handleGetSpacesForConfluenceURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(config.HeaderMattermostUserID)

	client, err := p.GetClientFromURL(r.URL.Path, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spaces, statusCode, err := client.GetSpaces()
	if err != nil {
		p.LogAndRespondError(w, statusCode, "Not able to get spaces for confluence url.", err)
		return
	}
	responseBody, err := json.Marshal(spaces)
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "Failed to marshal confluence spaces.", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(string(responseBody))); err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "failed to write response body.", err)
	}
}
