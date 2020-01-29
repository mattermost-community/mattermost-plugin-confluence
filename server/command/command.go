package command

import (
	"strings"

	"github.com/Brightscout/mattermost-plugin-confluence/server/service"

	"github.com/mattermost/mattermost-server/model"
)

type HandlerFunc func(context *model.CommandArgs, args ...string) *model.CommandResponse

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

var ConfluenceCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"list":        service.ListChannelSubscriptions,
		"unsubscribe": service.DeleteSubscription,
		"edit":        service.OpenSubscriptionEditModal,
	},
	defaultHandler: executeConfluenceDefault,
}

func GetCommand() *model.Command {
	return &model.Command{
		Trigger:          "confluence",
		DisplayName:      "Confluence",
		Description:      "Integration with Confluence.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: subscribe, list, unsubscribe \"<alias>\"",
		AutoCompleteHint: "[command]",
	}
}

// TODO : Show help text instead of invalid command.
func executeConfluenceDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Invalid command",
	}
}

func (ch Handler) Handle(context *model.CommandArgs, args ...string) *model.CommandResponse {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}
