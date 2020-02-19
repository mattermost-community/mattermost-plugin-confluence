package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

var getChannelSubscription = &Endpoint{
	RequiresAuth: true,
	Path:         "/subscription",
	Method:       http.MethodGet,
	Execute:      handleGetChannelSubscription,
}

func handleGetChannelSubscription(w http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channelID")
	alias := r.FormValue("alias")
	subscription, errCode, err := service.GetChannelSubscription(channelID, alias)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		return
	}
	b, _ := json.Marshal(subscription)
	fmt.Println("body=", string(b))
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(string(b)))
}
