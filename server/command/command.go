package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/command"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
)

type PluginAPI interface {
	GetBundlePath() (string, error)
}

type HandlerFunc func(context *model.CommandArgs, args ...string) *model.CommandResponse

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
	installServerHelp = `
To configure the plugin, create a new app in your Confluence Server following these steps:
1. Navigate to **Settings > Apps > Manage Apps**. For older versions of Confluence, navigate to **Administration > Applications > Add-ons > Manage add-ons**.
2. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.
3. Press **Upload app**.
4. Choose **From my computer** and upload the Mattermost for Confluence OBR file.
5. Once the app is installed, press **Configure** to open the configuration page.
6. In the **Webhook URL** field, enter: %s
7. Press **Save** to finish the setup.
`
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
	return confluence
}

func executeConfluenceDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
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

func (ch Handler) Handle(context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, commandsOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	if len(args) == 0 {
		return ch.handlers["help"](context, "")
	}
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}

func showInstallCloudHelp(context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, installOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	cloudURL := util.GetPluginURL() + util.GetAtlassianConnectURLPath()
	postCommandResponse(context, fmt.Sprintf(installCloudHelp, cloudURL))
	return &model.CommandResponse{}
}

func showInstallServerHelp(context *model.CommandArgs, args ...string) *model.CommandResponse {
	if !util.IsSystemAdmin(context.UserId) {
		postCommandResponse(context, installOnlySystemAdmin)
		return &model.CommandResponse{}
	}

	serverURL := util.GetPluginURL() + util.GetConfluenceServerWebhookURLPath()
	postCommandResponse(context, fmt.Sprintf(installServerHelp, serverURL))
	return &model.CommandResponse{}
}

func deleteSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
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

func listChannelSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
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

func confluenceHelpCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	helpText := getFullHelpText(context, args...)

	postCommandResponse(context, helpText)
	return &model.CommandResponse{}
}

func getFullHelpText(context *model.CommandArgs, args ...string) string {
	helpText := commonHelpText
	if util.IsSystemAdmin(context.UserId) {
		helpText += sysAdminHelpText
	}
	return helpText
}
