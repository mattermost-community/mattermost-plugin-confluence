package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
)

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

	invalidCommand         = "Invalid command parameters. Please use `/confluence help` for more information."
	installOnlySystemAdmin = "`/confluence install` can only be run by a system administrator."
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
		"help":           confluenceHelp,
	},
	defaultHandler: executeConfluenceDefault,
}

func GetCommand() *model.Command {
	return &model.Command{
		Trigger:          "confluence",
		DisplayName:      "Confluence",
		Description:      "Integration with Confluence.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: subscribe, list, unsubscribe \"<name>\", edit \"<name>\", install cloud/server, help.",
		AutoCompleteHint: "[command]",
	}
}

func executeConfluenceDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         invalidCommand,
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
	alias := args[0]
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

func confluenceHelp(context *model.CommandArgs, args ...string) *model.CommandResponse {
	helpText := commonHelpText
	if util.IsSystemAdmin(context.UserId) {
		helpText += sysAdminHelpText
	}

	postCommandResponse(context, helpText)
	return &model.CommandResponse{}
}
