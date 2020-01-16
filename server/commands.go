package main

import (
	"fmt"
	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"strings"
)

type CommandHandlerFunc func(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) *model.CommandResponse

type CommandHandler struct {
	handlers       map[string]CommandHandlerFunc
	defaultHandler CommandHandlerFunc
}


var confluenceCommandHandler = CommandHandler{
	handlers: map[string]CommandHandlerFunc{
		"config": executeConfig,
	},
	defaultHandler: executeConflunceDefault,
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "confluence",
		DisplayName:      "Confluence",
		Description:      "Integration with Confluence.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: config, help",
		AutoCompleteHint: "[command]",
	}
}

func executeConfig(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) *model.CommandResponse {
	dialogRequest := model.OpenDialogRequest{
		TriggerId: header.TriggerId,
		URL:       fmt.Sprintf("%s/plugins/%s/abc", *config.Mattermost.GetConfig().ServiceSettings.SiteURL, config.PluginName),
		Dialog: model.Dialog{
			Title: "Configure confluence",
			Elements: []model.DialogElement{
				{
					DisplayName: "URL",
					Name:        "url",
					Type:        "text",
					Placeholder: "Please enter the url.",
				},
				{
					DisplayName: "Space Key",
					Name:        "spaceKey",
					Type:        "text",
					Placeholder: "Please enter the space key.",

				},
				{
					DisplayName: "Notification",
					Name:        "abc",
					Type:        "bool",
					Placeholder: "comments",
				},
			},
		},
	}
	if err := config.Mattermost.OpenInteractiveDialog(dialogRequest); err != nil {
		config.Mattermost.LogError("Error opening config modal: ", err.Error())
	}
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "",
	}
}

func executeConflunceDefault(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "",
	}
}

func (ch CommandHandler) Handle(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) *model.CommandResponse {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(p, c, header, args[n:]...)
		}
	}
	return ch.defaultHandler(p, c, header, args...)
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args := strings.Fields(commandArgs.Command)
	return confluenceCommandHandler.Handle(p, c, commandArgs, args[1:]...), nil
}
