package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/command"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/store"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-plugin-confluence/server/util/types"
)

type PluginAPI interface {
	GetBundlePath() (string, error)
}

type HandlerFunc func(p *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

const (
	specifyAlias              = "Please specify a subscription name."
	subscriptionDeleteSuccess = "Subscription **%s** has been deleted."
	noChannelSubscription     = "No subscriptions found for this channel."
	commonHelpText            = "###### Mattermost Confluence Plugin - Slash Command Help\n\n" +
		"* `/confluence subscribe` - Subscribe the current channel to notifications from Confluence.\n" +
		"* `/confluence unsubscribe \"<name>\"` - Unsubscribe the current channel from notifications associated with the given subscription name.\n" +
		"* `/confluence list` - List all subscriptions for the current channel.\n" +
		"* `/confluence edit \"<name>\"` - Edit the subscription settings associated with the given subscription name.\n"

	sysAdminHelpText = "\n###### For System Administrators:\n" +
		"Setup Instructions:\n" +
		"* `/confluence install cloud` - Connect Mattermost to a Confluence Cloud instance.\n" +
		"* `/confluence install server` - Connect Mattermost to a Confluence Server or Data Center instance.\n"

	invalidCommand          = "Invalid command."
	installOnlySystemAdmin  = "`/confluence install` can only be run by a system administrator."
	commandsOnlySystemAdmin = "`/confluence` commands can only be run by a system administrator."
)

const (
	installCloudHelp = `
To finish the configuration, add a new app in your Confluence Cloud instance following these steps:
1. Navigate to **Settings > Apps > Manage Apps**.
2. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.
3. Press **Upload App**.
4. In **From this URL**, enter: %s
5. Once installed, you will see the "Installed and ready to go!" message.
`
)

var ConfluenceCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"list":           listChannelSubscription,
		"unsubscribe":    deleteSubscription,
		"install/cloud":  showInstallCloudHelp,
		"install/server": showInstallServerHelp,
		"connect":        executeConnect,
		"disconnect":     executeDisconnect,
		"help":           confluenceHelpCommand,
	},
	defaultHandler: executeConfluenceDefault,
}

func GetCommand(pAPI PluginAPI) (*model.Command, error) {
	iconData, err := command.GetIconData(pAPI, "assets/icon.svg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get icon data")
	}

	return &model.Command{
		Trigger:              "confluence",
		DisplayName:          "Confluence",
		Description:          "Integration with Confluence.",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: subscribe, list, unsubscribe, edit, install, help.",
		AutoCompleteHint:     "[command]",
		AutocompleteData:     getAutoCompleteData(),
		AutocompleteIconData: iconData,
	}, nil
}

func getAutoCompleteData() *model.AutocompleteData {
	confluence := model.NewAutocompleteData("confluence", "[command]", "Available commands: subscribe, list, unsubscribe, edit, install, help")

	install := model.NewAutocompleteData("install", "", "Connect Mattermost to a Confluence instance")
	installItems := []model.AutocompleteListItem{{
		HelpText: "Connect Mattermost to a Confluence Cloud instance",
		Item:     "cloud",
	}, {
		HelpText: "Connect Mattermost to a Confluence Server or Data Center instance",
		Item:     "server",
	}}
	install.AddStaticListArgument("", false, installItems)
	confluence.AddCommand(install)

	list := model.NewAutocompleteData("list", "", "List all subscriptions for the current channel")
	confluence.AddCommand(list)

	edit := model.NewAutocompleteData("edit", "[name]", "Edit the subscription settings associated with the given subscription name")
	edit.AddDynamicListArgument("name", "api/v1/autocomplete/GetChannelSubscriptions", false)
	confluence.AddCommand(edit)

	subscribe := model.NewAutocompleteData("subscribe", "", "Subscribe the current channel to notifications from Confluence")
	confluence.AddCommand(subscribe)

	unsubscribe := model.NewAutocompleteData("unsubscribe", "[name]", "Unsubscribe the current channel from notifications associated with the given subscription name")
	unsubscribe.AddDynamicListArgument("name", "api/v1/autocomplete/GetChannelSubscriptions", false)
	confluence.AddCommand(unsubscribe)

	help := model.NewAutocompleteData("help", "", "Show confluence slash command help")
	confluence.AddCommand(help)

	connect := model.NewAutocompleteData("connect", "", "Connect your Mattermost account to your Confluence account")
	confluence.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your Mattermost account from your Confluence account")
	confluence.AddCommand(disconnect)

	return confluence
}

func executeConfluenceDefault(p *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	out := invalidCommand + "\n\n"
	out += getFullHelpText(context, args...)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         out,
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

func (ch Handler) Handle(p *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, commandsOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	if len(args) == 0 {
		return ch.handlers["help"](p, context, "")
	}
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(p, context, args[n:]...)
		}
	}
	return ch.defaultHandler(p, context, args...)
}

func (p *Plugin) postCommandResponse(context *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: context.ChannelId,
		Message:   text,
	}
	_ = config.Mattermost.SendEphemeralPost(context.UserId, post)
}

func (p *Plugin) responsef(commandArgs *model.CommandArgs, format string, args ...interface{}) *model.CommandResponse {
	p.postCommandResponse(commandArgs, fmt.Sprintf(format, args...))
	return &model.CommandResponse{}
}

func executeConnect(p *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	isAdmin := util.IsSystemAdmin(context.UserId)
	if !isAdmin {
		return p.responsef(context, "Command is required for admins only")
	}

	pluginConfig := config.GetConfig()
	if pluginConfig.ConfluenceURL == "" || pluginConfig.ConfluenceOAuthClientID == "" || pluginConfig.ConfluenceOAuthClientSecret == "" {
		return p.responsef(context, "Oauth config not set for confluence plugin. Please run `/confluence install server`")
	}
	confluenceURL := pluginConfig.ConfluenceURL
	confluenceURL = strings.TrimSuffix(confluenceURL, "/")

	conn, err := store.LoadConnection(types.ID(confluenceURL), types.ID(context.UserId), p.pluginVersion)
	if err == nil && len(conn.ConfluenceAccountID()) != 0 {
		return p.responsef(context,
			"You already have a Confluence account linked to your Mattermost account. Please use `/confluence disconnect` to disconnect.")
	}

	link := fmt.Sprintf("%s/oauth2/connect", util.GetPluginURL())
	return p.responsef(context, "[Click here to link your Confluence account](%s)", link)
}

func executeDisconnect(p *Plugin, commArgs *model.CommandArgs, args ...string) *model.CommandResponse {
	user, err := store.LoadUser(types.ID(commArgs.UserId))
	if err != nil {
		return p.responsef(commArgs, "Could not complete the **disconnection** request. Error: %v", err)
	}
	confluenceURL := user.InstanceURL.String()

	disconnected, err := p.DisconnectUser(confluenceURL, types.ID(commArgs.UserId))
	if errors.Cause(err) == store.ErrNotFound {
		errorStr := "Your account is not connected to Confluence. Please use `/confluence connect <instance url>` to connect your account."
		if confluenceURL != "" {
			errorStr = fmt.Sprintf("You don't have a Confluence account at %s linked to your Mattermost account currently. Please use `/confluence connect` to connect your account.", confluenceURL)
		}
		return p.responsef(commArgs, errorStr)
	}
	if err != nil {
		return p.responsef(commArgs, "Could not complete the **disconnection** request. Error: %v", err)
	}
	return p.responsef(commArgs, "You have successfully disconnected your Confluence account (**%s**).", disconnected.DisplayName)
}

func showInstallCloudHelp(_ *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, installOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	cloudURL := util.GetPluginURL() + util.GetAtlassianConnectURLPath()
	postCommandResponse(context, fmt.Sprintf(installCloudHelp, cloudURL))
	return &model.CommandResponse{}
}

func showInstallServerHelp(p *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, installOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	err := p.flowManager.StartSetupWizard(context.UserId, "")
	if err != nil {
		return &model.CommandResponse{}
	}

	postCommandResponse(context, "Please continue with confluence bot DM for the setup.")

	return &model.CommandResponse{}
}

func deleteSubscription(_ *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	if len(args) == 0 {
		postCommandResponse(context, specifyAlias)
		return &model.CommandResponse{}
	}
	alias := strings.Join(args, " ")
	if err := service.DeleteSubscription(context.ChannelId, alias); err != nil {
		postCommandResponse(context, err.Error())
		return &model.CommandResponse{}
	}
	postCommandResponse(context, fmt.Sprintf(subscriptionDeleteSuccess, alias))
	return &model.CommandResponse{}
}

func listChannelSubscription(_ *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions, gErr := service.GetSubscriptionsByChannelID(context.ChannelId)
	if gErr != nil {
		postCommandResponse(context, gErr.Error())
		return &model.CommandResponse{}
	}

	if len(channelSubscriptions) == 0 {
		postCommandResponse(context, noChannelSubscription)
		return &model.CommandResponse{}
	}
	list := serializer.FormattedSubscriptionList(channelSubscriptions)
	postCommandResponse(context, list)
	return &model.CommandResponse{}
}

func confluenceHelpCommand(_ *Plugin, context *model.CommandArgs, args ...string) *model.CommandResponse {
	helpText := getFullHelpText(context, args...)

	postCommandResponse(context, helpText)
	return &model.CommandResponse{}
}

func getFullHelpText(context *model.CommandArgs, _ ...string) string {
	helpText := commonHelpText
	if util.IsSystemAdmin(context.UserId) {
		helpText += sysAdminHelpText
	}
	return helpText
}
