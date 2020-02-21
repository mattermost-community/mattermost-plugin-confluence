package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

type HandlerFunc func(context *model.CommandArgs, args ...string) *model.CommandResponse

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

const (
	specifyAlias              = "Please specify an alias."
	subscriptionDeleteSuccess = "**%s** has been deleted."
	noChannelSubscription     = "No subscriptions found for this channel."
	commonHelpText            = "###### Mattermost Confluence Plugin - Slash Command Help\n\n" +
		"* `/confluence subscribe` - Subscribe the current channel to notifications from Confluence.\n" +
		"* `/confluence unsubscribe \"<alias>\"` - Unsubscribe the current channel from notifications associated with the given alias.\n" +
		"* `/confluence list` - List all subscriptions for the current channel.\n" +
		"* `/confluence edit \"<alias>\"` - Edit the subscription settings associated with the given alias.\n"

	sysAdminHelpText = "\n###### For System Administrators:\n" +
		"Setup Instructions:\n" +
		"* `/confluence install cloud` - Connect Mattermost to a Confluence Cloud instance.\n" +
		"* `/confluence install server` - Connect Mattermost to a Confluence Server or Data Center instance.\n"

	invalidCommand         = "Invalid command parameters. Please use `/confluence help` for more information."
	installOnlySystemAdmin = "`/confluence install` can only be run by a system administrator."
)

const (
	installServerHelp = `
To configure the plugin, create a new app in your Confluence instance following these steps:
1. Navigate to **Settings > Apps > Manage Apps**.
  - For older versions of Confluence, navigate to **Administration > Applications > Add-ons > Manage add-ons**.
2. Click **Settings** at bottom of page, enable development mode, and apply this change.
  - Enabling development mode allows you to install apps that are not from the Atlassian Marketplace.
3. Click **Upload app**.
4. Chose 'From my computer' and upload the **Mattermost for Confluence OBR** file.
5. Wait for the app to install.
6. Use the 'configure' button to open the **Mattermost Configuration** page.
7. Enter the following URL as the **Webhook URL** and click on Save.
%s
`
	installCloudHelp = `
To finish the configuration, add a new app in your Confluence instance following these steps:
1. Navigate to **Settings > Apps > Manage Apps**.
2. Click **Settings** at bottom of page, enable development mode, and apply this change.
  - Enabling development mode allows you to install apps that are not from the Atlassian Marketplace.
3. Click **Upload app**.
4. In the **From this URL field**, enter: %s
5. Wait for the app to install. Once completed, you should see an "Installed and ready to go!" message.
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
		AutoCompleteDesc: "Available commands: subscribe, list, unsubscribe \"<alias>\", edit \"<alias>\", install cloud/server, help.",
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
	if err := service.DeleteSubscription(context.ChannelId, args[0]); err != nil {
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
