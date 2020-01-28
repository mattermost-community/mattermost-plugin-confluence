package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type CommandHandlerFunc func(context *model.CommandArgs, args ...string) *model.CommandResponse

type CommandHandler struct {
	handlers       map[string]CommandHandlerFunc
	defaultHandler CommandHandlerFunc
}

const (
	OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT = "open_edit_subscription_modal"
)

var (
	confluenceCommandHandler = CommandHandler{
		handlers: map[string]CommandHandlerFunc{
			"list":        listChannelSubscriptions,
			"unsubscribe": deleteSubscription,
			"edit":        editSubscription,
		},
		defaultHandler: executeConflunceDefault,
	}

	eventTypes = map[string]string{
		"comment_create": "Comment Create",
		"comment_update": "Comment Update",
		"comment_delete": "Comment Delete",
		"page_create":    "Page Create",
		"page_update":    "Page Update",
		"page_delete":    "Page Delete",
	}
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "confluence",
		DisplayName:      "Confluence",
		Description:      "Integration with Confluence.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: subscribe, list, unsubscribe \"<alias>\"",
		AutoCompleteHint: "[command]",
	}
}

func executeConflunceDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Invalid command",
	}
}

func postCommandResponse(context *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: context.ChannelId,
		Message:   text,
	}
	_ = config.Mattermost.SendEphemeralPost(context.UserId, post)
}

func (ch CommandHandler) Handle(context *model.CommandArgs, args ...string) *model.CommandResponse {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}

func (p *Plugin) ExecuteCommand(context *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args, argErr := util.SplitArgs(commandArgs.Command)
	if argErr != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         argErr.Error(),
		}, nil
	}
	return confluenceCommandHandler.Handle(commandArgs, args[1:]...), nil
}

func listChannelSubscriptions(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	if err := util.Get(context.ChannelId, &channelSubscriptions); err != nil {
		postCommandResponse(context, "Encountered an error getting channel subscriptions.")
		return &model.CommandResponse{}
	}
	if len(channelSubscriptions) == 0 {
		postCommandResponse(context, "No subscription found for this channel.")
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
	postCommandResponse(context, text)
	return &model.CommandResponse{}
}

func deleteSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	alias := args[0]
	if err := util.Get(context.ChannelId, &channelSubscriptions); err != nil {
		postCommandResponse(context, fmt.Sprintf("Error occured while deleting subscription with alias **%s**.", alias))
		return &model.CommandResponse{}
	}
	if subscription, ok := channelSubscriptions[alias]; ok {
		if err := deleteSubscriptionUtil(subscription, channelSubscriptions, alias); err != nil {
			postCommandResponse(context, fmt.Sprintf("Error occured while deleting subscription with alias **%s**.", alias))
			return &model.CommandResponse{}
		}
		postCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** deleted successfully.", alias))
		return &model.CommandResponse{}
	} else {
		postCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** not found.", alias))
		return &model.CommandResponse{}
	}
}

func deleteSubscriptionUtil(subscription serializer.Subscription, channelSubscriptions map[string]serializer.Subscription, alias string) error {
	key, kErr := util.GetKey(subscription.BaseURL, subscription.SpaceKey)
	if kErr != nil {
		return kErr
	}
	keySubscriptions := make(map[string][]string)
	if err := util.Get(key, &keySubscriptions); err != nil {
		return err
	}
	delete(keySubscriptions, subscription.ChannelID)
	delete(channelSubscriptions, alias)
	if err := util.Set(key, keySubscriptions); err != nil {
		return err
	}
	if err := util.Set(subscription.ChannelID, channelSubscriptions); err != nil {
		return err
	}
	return nil
}

func editSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions := make(map[string]serializer.Subscription)
	alias := args[0]
	if err := util.Get(context.ChannelId, &channelSubscriptions); err != nil {
		postCommandResponse(context, fmt.Sprintf("Error occured while editing subscription with alias **%s**.", alias))
		return &model.CommandResponse{}
	}
	if subscription, ok := channelSubscriptions[alias]; ok {
		bytes, err := json.Marshal(subscription)
		if err != nil {
			postCommandResponse(context, fmt.Sprintf("Error occured while editing subscription with alias **%s**.", alias))
			return &model.CommandResponse{}
		}
		config.Mattermost.PublishWebSocketEvent(
			OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT,
			map[string]interface{}{
				"subscription": string(bytes),
			},
			&model.WebsocketBroadcast{
				UserId: context.UserId,
			},
		)
		return &model.CommandResponse{}
	} else {
		postCommandResponse(context, fmt.Sprintf("Subscription with alias **%s** not found.", alias))
		return &model.CommandResponse{}
	}
}
