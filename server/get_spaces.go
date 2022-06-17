package main

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils/types"
)

func (p *Plugin) handleGetSpacesForConfluenceURL(w http.ResponseWriter, r *http.Request) {
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

	spaces, err := client.GetSpacesForConfluenceURL()
	if err != nil {
		p.LogAndRespondError(w, http.StatusInternalServerError, "not able to get Spaces for confluence url.", err)
		return
	}
	b, _ := json.Marshal(spaces)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(string(b)))
}
