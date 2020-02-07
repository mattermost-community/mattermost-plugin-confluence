package controller

import (
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
	channelID :=  r.FormValue("channelID")
	alias := r.FormValue("alias")
	subscription, err, errCode := service.GetChannelSubscription(channelID, alias)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(subscription.ToJSON()))
}
