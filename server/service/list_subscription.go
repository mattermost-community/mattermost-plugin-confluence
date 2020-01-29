package service

import (
	"fmt"
	"strings"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-server/model"
)

var eventTypes = map[string]string{
	"comment_create": "Comment Create",
	"comment_update": "Comment Update",
	"comment_delete": "Comment Delete",
	"page_create":    "Page Create",
	"page_update":    "Page Update",
	"page_delete":    "Page Delete",
}

func ListChannelSubscriptions(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	if err := store.Get(store.GetChannelSubscriptionKey(context.ChannelId), &channelSubscriptions); err != nil {
		util.PostCommandResponse(context, "Encountered an error getting channel subscriptions.")
		return &model.CommandResponse{}
	}
	if len(channelSubscriptions) == 0 {
		util.PostCommandResponse(context, "No subscription found for this channel.")
		return &model.CommandResponse{}
	}
	text := fmt.Sprintf("| Alias | Base Url | Space Key | Events|\n| :----: |:--------:| :--------:| :-----:|")
	for _, subscription := range channelSubscriptions {
		var events []string
		for _, event := range subscription.Events {
			events = append(events, eventTypes[event])
		}
		text += fmt.Sprintf("\n|%s|%s|%s|%s|", subscription.Alias, subscription.BaseURL, subscription.SpaceKey, strings.Join(events, ", "))
	}
	util.PostCommandResponse(context, text)
	return &model.CommandResponse{}
}
