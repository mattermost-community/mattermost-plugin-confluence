package main

import (
	"fmt"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/flow"
)

type Tracker interface {
	TrackEvent(event string, properties map[string]interface{})
	TrackUserEvent(event, userID string, properties map[string]interface{})
}

type FlowManager struct {
	client            *pluginapi.Client
	pluginID          string
	botUserID         string
	router            *mux.Router
	getConfiguration  func() *config.Configuration
	webhookURL        string
	confluenceBaseURL string
	tracker           Tracker
	setupFlow         *flow.Flow
	announcementFlow  *flow.Flow
}

func (p *Plugin) NewFlowManager() (*FlowManager, error) {
	webhookURL := util.GetPluginURL() + util.GetConfluenceServerWebhookURLPath()

	fm := &FlowManager{
		client:           p.client,
		pluginID:         manifest.Id,
		botUserID:        p.BotUserID,
		router:           p.Router,
		webhookURL:       webhookURL,
		getConfiguration: config.GetConfig,

		tracker: p,
	}

	setupFlow, err := fm.newFlow("setup")
	if err != nil {
		return nil, err
	}

	setupFlow.WithSteps(
		fm.stepWelcome(),
		fm.stepInstanceURL(),
		fm.stepServerVersionQuestion(),
		fm.stepCSversionGreaterthan9(),
		fm.stepCSversionLessthan9(),
		fm.stepOAuthInput(),
		fm.stepOAuthConnect(),
		fm.stepAnnouncementQuestion(),
		fm.stepAnnouncementConfirmation(),
		fm.stepDone(),
		fm.stepCancel("setup"),
	)
	fm.setupFlow = setupFlow

	announcementFlow, err := fm.newFlow("announcement")
	if err != nil {
		return nil, err
	}
	announcementFlow.WithSteps(
		fm.stepAnnouncementQuestion(),
		fm.stepAnnouncementConfirmation().Terminal(),

		fm.stepCancel("setup announcement"),
	)
	fm.announcementFlow = announcementFlow

	return fm, nil
}

func (fm *FlowManager) newFlow(name flow.Name) (*flow.Flow, error) {
	flow, err := flow.NewFlow(
		name,
		fm.client,
		fm.pluginID,
		fm.botUserID,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create flow %s", name)
	}

	flow.InitHTTP(fm.router)

	return flow, nil
}

const (
	stepWelcome                  flow.Name = "welcome"
	stepServerVersionQuestion    flow.Name = "server-verstion-question"
	stepConfluenceURL            flow.Name = "confluence-url"
	stepOAuthInput               flow.Name = "oauth-input"
	stepCSversionLessthan9       flow.Name = "server-version-less-than-9"
	stepCSversionGreaterthan9    flow.Name = "server-version-greater-than-9"
	stepAnnouncementQuestion     flow.Name = "announcement-question"
	stepAnnouncementConfirmation flow.Name = "announcement-confirmation"
	stepDone                     flow.Name = "done"
	stepCancel                   flow.Name = "cancel"
	stepOAuthConnect             flow.Name = "oauth-connect"

	keyConfluenceURL     = "ConfluenceURL"
	keyIsOAuthConfigured = "IsOAuthConfigured"
	redirectURL          = "dummyRedirectURL" // will be added in oauth PR
)

func cancelButton() flow.Button {
	return flow.Button{
		Name:    "Cancel setup",
		Color:   flow.ColorDanger,
		OnClick: flow.Goto(stepCancel),
	}
}

func (fm *FlowManager) stepCancel(command string) flow.Step {
	return flow.NewStep(stepCancel).
		Terminal().
		WithText(fmt.Sprintf("Confluence integration setup has stopped. Restart setup later by running `/confluence %s`. Learn more about the plugin [here](https://mattermost.gitbook.io/plugin-confluence/).", command)).
		WithColor(flow.ColorDanger)
}

func continueButtonF(f func(f *flow.Flow) (flow.Name, flow.State, error)) flow.Button {
	return flow.Button{
		Name:    "Continue",
		Color:   flow.ColorPrimary,
		OnClick: f,
	}
}

func continueButton(next flow.Name) flow.Button {
	return continueButtonF(flow.Goto(next))
}

func (fm *FlowManager) getBaseState() flow.State {
	config := fm.getConfiguration()
	isOAuthConfigured := config.ConfluenceOAuthClientID != "" || config.ConfluenceOAuthClientSecret != ""
	return flow.State{
		keyConfluenceURL:     config.ConfluenceURL,
		keyIsOAuthConfigured: isOAuthConfigured,
	}
}

func (fm *FlowManager) StartSetupWizard(userID string, delegatedFrom string) error {
	state := fm.getBaseState()

	err := fm.setupFlow.ForUser(userID).Start(state)
	if err != nil {
		return err
	}

	fm.client.Log.Debug("Started setup wizard", "userID", userID, "delegatedFrom", delegatedFrom)

	fm.trackStartSetupWizard(userID, delegatedFrom != "")

	return nil
}

func (fm *FlowManager) trackStartSetupWizard(userID string, fromInvite bool) {
	fm.tracker.TrackUserEvent("setup_wizard_start", userID, map[string]interface{}{
		"from_invite": fromInvite,
		"time":        model.GetMillis(),
	})
}

func (fm *FlowManager) trackCompleteSetupWizard(userID string) {
	fm.tracker.TrackUserEvent("setup_wizard_complete", userID, map[string]interface{}{
		"time": model.GetMillis(),
	})
}

func (fm *FlowManager) stepWelcome() flow.Step {
	welcomeText := ":wave: Welcome to your Confluence integration! [Learn more](https://mattermost.gitbook.io/plugin-confluence/)"
	welcomePretext := "Just a few configuration steps to go!"

	return flow.NewStep(stepWelcome).
		WithText(welcomeText).
		WithPretext(welcomePretext).
		WithButton(continueButton(stepConfluenceURL))
}

func (fm *FlowManager) stepServerVersionQuestion() flow.Step {
	delegateQuestionText := "Are you using confluence server version greater than or equal to 9?"
	return flow.NewStep(stepServerVersionQuestion).
		WithText(delegateQuestionText).
		WithButton(flow.Button{
			Name:    "Yes",
			Color:   flow.ColorPrimary,
			OnClick: flow.Goto(stepCSversionGreaterthan9),
		}).
		WithButton(flow.Button{
			Name:  "No",
			Color: flow.ColorDefault,
			OnClick: func(f *flow.Flow) (flow.Name, flow.State, error) {
				return stepCSversionLessthan9, nil, nil
			},
		})
}

func (fm *FlowManager) stepCSversionGreaterthan9() flow.Step {
	return flow.NewStep(stepCSversionGreaterthan9).
		WithText(
			fmt.Sprintf(
				"%s has been successfully added. To finish the configuration, add an Application Link in your Confluence instance following these steps:\n",
				fm.confluenceBaseURL,
			) +
				"1. Go to [**Settings > Applications > Application Links**]({{ .ConfluenceURL }}/plugins/servlet/applinks/listApplicationLinks)\n" +
				"   ![image](https://user-images.githubusercontent.com/90389917/202149868-a3044351-37bc-43c0-9671-aba169706917.png)\n" +
				"2. Select **Create link**.\n" +
				"3. On the **Create Link** screen, select **External Application** and **Incoming** as `Application type` and `Direction` respectively. Select **Continue**.\n" +
				"4. On the **Link Applications** screen, set the following values:\n" +
				"   - **Name**: `Mattermost`\n" +
				fmt.Sprintf("   - **Redirect URL**: `%s`\n", redirectURL) +
				"   - **Application Permissions**: `Admin`\n" +
				"   Select **Continue**.\n" +
				"5. Copy the `clientID` and `clientSecret` from **Settings**, and paste them into the modal in Mattermost which can be opened by using the `/confluence config add` slash command.\n" +
				"6. In Mattermost, use the `/confluence connect {{ .ConfluenceURL }} admin` slash command to connect your Mattermost account with your Confluence admin account and save the token of the admin to handle admin-restricted functions.\n" +
				"7. Use the `/confluence connect` slash command to connect your Mattermost account with your Confluence account for all other users.\n" +
				"If you see an option to create a Confluence issue, you're all set! If not, refer to our [documentation](https://mattermost.gitbook.io/plugin-confluence) for troubleshooting help.",
		).
		WithButton(continueButton(stepOAuthInput))
}

func (fm *FlowManager) stepCSversionLessthan9() flow.Step {
	return flow.NewStep(stepCSversionLessthan9).
		WithText(fmt.Sprintf(`
To configure the plugin, create a new app in your [Confluence Server](%s) following these steps:
1. Navigate to **Settings > Apps > Manage Apps**. For older versions of Confluence, navigate to **Administration > Applications > Add-ons > Manage add-ons**.
2. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.
3. Press **Upload app**.
4. Choose **From my computer** and upload the Mattermost for Confluence OBR file.
5. Once the app is installed, press **Configure** to open the configuration page.
6. In the **Webhook URL** field, enter: %s
7. Press **Save** to finish the setup.
`, fm.confluenceBaseURL, fm.webhookURL)).
		WithButton(continueButton(stepDone))
}

func (fm *FlowManager) stepInstanceURL() flow.Step {
	enterpriseText := "Click the Continue button below to open a dialog to enter the **Confluence URL**"
	return flow.NewStep(stepConfluenceURL).
		WithText(enterpriseText).
		WithButton(flow.Button{
			Name:  "Continue",
			Color: flow.ColorPrimary,
			Dialog: &model.Dialog{
				Title:            "Confluence URL",
				IntroductionText: "Enter the **Confluence URL** of your Confluence instance (Example: https://confluence.example.com).",
				SubmitLabel:      "Save & continue",
				Elements: []model.DialogElement{
					{

						DisplayName: "Confluence URL",
						Name:        "confluence_url",
						Type:        "text",
						SubType:     "url",
						Placeholder: "Enter Confluence URL",
					},
				},
			},
			OnDialogSubmit: fm.submitConfluenceURL,
		}).
		WithButton(cancelButton())
}

func (fm *FlowManager) submitConfluenceURL(f *flow.Flow, submitted map[string]interface{}) (flow.Name, flow.State, map[string]string, error) {
	errorList := map[string]string{}

	confluenceURLRaw, ok := submitted["confluence_url"]
	if !ok {
		return "", nil, nil, errors.New("confluence_url missing")
	}
	confluenceURL, ok := confluenceURLRaw.(string)
	if !ok {
		return "", nil, nil, errors.New("confluence_url is not a string")
	}

	// _, err := service.CheckConfluenceURL(fm.MMSiteURL, confluenceURL, false)
	// if err != nil {
	// 	errorList["confluence_url"] = err.Error()
	// }

	if len(errorList) != 0 {
		return "", nil, errorList, nil
	}

	config := fm.getConfiguration()
	config.ConfluenceURL = confluenceURL
	config.Sanitize()

	configMap, err := config.ToMap()
	if err != nil {
		return "", nil, nil, err
	}

	err = fm.client.Configuration.SavePluginConfig(configMap)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "failed to save plugin config")
	}

	fm.confluenceBaseURL = confluenceURL

	return stepServerVersionQuestion, flow.State{
		keyConfluenceURL: config.ConfluenceURL,
	}, nil, nil
}

func (fm *FlowManager) stepOAuthInput() flow.Step {
	return flow.NewStep(stepOAuthInput).
		WithText("Click the Continue button below to open a dialog to enter the **Application ID** and **Secret**.").
		WithButton(flow.Button{
			Name:  "Continue",
			Color: flow.ColorPrimary,
			Dialog: &model.Dialog{
				Title:            "Confluence OAuth Credentials",
				IntroductionText: "Please enter the **Application ID** and **Secret** you copied in a previous step.{{ if .IsOAuthConfigured }}\n\n**Any existing OAuth configuration will be overwritten.**{{end}}",
				SubmitLabel:      "Save & continue",
				Elements: []model.DialogElement{
					{
						DisplayName: "Confluence OAuth Application ID",
						Name:        "client_id",
						Type:        "text",
						SubType:     "text",
						Placeholder: "Enter Confluence OAuth Application ID",
					},
					{
						DisplayName: "Confluence OAuth Secret",
						Name:        "client_secret",
						Type:        "text",
						SubType:     "text",
						Placeholder: "Enter Confluence OAuth Secret",
					},
				},
			},
			OnDialogSubmit: fm.submitOAuthConfig,
		}).
		WithButton(cancelButton())
}

func (fm *FlowManager) submitOAuthConfig(f *flow.Flow, submitted map[string]interface{}) (flow.Name, flow.State, map[string]string, error) {
	errorList := map[string]string{}

	clientIDRaw, ok := submitted["client_id"]
	if !ok {
		return "", nil, nil, errors.New("client_id missing")
	}
	clientID, ok := clientIDRaw.(string)
	if !ok {
		return "", nil, nil, errors.New("client_id is not a string")
	}

	clientID = strings.TrimSpace(clientID)

	if len(clientID) < 64 {
		errorList["client_id"] = "Client ID should be at least 64 characters long"
	}

	clientSecretRaw, ok := submitted["client_secret"]
	if !ok {
		return "", nil, nil, errors.New("client_secret missing")
	}
	clientSecret, ok := clientSecretRaw.(string)
	if !ok {
		return "", nil, nil, errors.New("client_secret is not a string")
	}

	clientSecret = strings.TrimSpace(clientSecret)

	if len(clientSecret) < 64 {
		errorList["client_secret"] = "Client Secret should be at least 64 characters long"
	}

	if len(errorList) != 0 {
		return "", nil, errorList, nil
	}

	config := fm.getConfiguration()
	config.ConfluenceOAuthClientID = clientID
	config.ConfluenceOAuthClientSecret = clientSecret

	configMap, err := config.ToMap()
	if err != nil {
		return "", nil, nil, err
	}

	err = fm.client.Configuration.SavePluginConfig(configMap)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "failed to save plugin config")
	}

	return stepOAuthConnect, nil, nil, nil
}

func (fm *FlowManager) stepOAuthConnect() flow.Step {
	connectPretext := "##### :white_check_mark: Step 1: Connect your Confluence account"
	connectURL := fmt.Sprintf("%s/oauth/connect", util.GetPluginURL())
	connectText := fmt.Sprintf("Go [here](%s) to connect your account.", connectURL)
	return flow.NewStep(stepOAuthConnect).
		WithText(connectText).
		WithPretext(connectPretext)
}

func (fm *FlowManager) StartAnnouncementWizard(userID string) error {
	state := fm.getBaseState()

	err := fm.announcementFlow.ForUser(userID).Start(state)
	if err != nil {
		return err
	}

	fm.trackStartAnnouncementWizard(userID)

	return nil
}

func (fm *FlowManager) trackStartAnnouncementWizard(userID string) {
	fm.tracker.TrackUserEvent("announcement_wizard_start", userID, map[string]interface{}{
		"time": model.GetMillis(),
	})
}

func (fm *FlowManager) trackCompletAnnouncementWizard(userID string) {
	fm.tracker.TrackUserEvent("announcement_wizard_complete", userID, map[string]interface{}{
		"time": model.GetMillis(),
	})
}

func (fm *FlowManager) stepAnnouncementQuestion() flow.Step {
	defaultMessage := "Hi team,\n" +
		"\n" +
		"We've set up the Mattermost Confluence plugin to enable notifications from Confluence in Mattermost. To get started, run the `/confluence connect` slash command from any channel within Mattermost to connect that channel with Confluence. See the [documentation](https://mattermost.gitbook.io/plugin-confluence/) for details on using the Confluence plugin."

	return flow.NewStep(stepAnnouncementQuestion).
		WithText("Want to let your team know?").
		WithButton(flow.Button{
			Name:  "Send Message",
			Color: flow.ColorPrimary,
			Dialog: &model.Dialog{
				Title:       "Notify your team",
				SubmitLabel: "Send message",
				Elements: []model.DialogElement{
					{
						DisplayName: "To",
						Name:        "channel_id",
						Type:        "select",
						Placeholder: "Select channel",
						DataSource:  "channels",
					},
					{
						DisplayName: "Message",
						Name:        "message",
						Type:        "textarea",
						Default:     defaultMessage,
						HelpText:    "You can edit this message before sending it.",
					},
				},
			},
			OnDialogSubmit: fm.submitChannelAnnouncement,
		}).
		WithButton(flow.Button{
			Name:    "Not now",
			Color:   flow.ColorDefault,
			OnClick: flow.Goto(stepDone),
		})
}

func (fm *FlowManager) stepAnnouncementConfirmation() flow.Step {
	return flow.NewStep(stepAnnouncementConfirmation).
		WithText("Message to ~{{ .ChannelName }} was sent.").
		Next("").
		OnRender(func(f *flow.Flow) { fm.trackCompletAnnouncementWizard(f.UserID) })
}

func (fm *FlowManager) submitChannelAnnouncement(f *flow.Flow, submitted map[string]interface{}) (flow.Name, flow.State, map[string]string, error) {
	channelIDRaw, ok := submitted["channel_id"]
	if !ok {
		return "", nil, nil, errors.New("channel_id missing")
	}
	channelID, ok := channelIDRaw.(string)
	if !ok {
		return "", nil, nil, errors.New("channel_id is not a string")
	}

	channel, err := fm.client.Channel.Get(channelID)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "failed to get channel")
	}

	messageRaw, ok := submitted["message"]
	if !ok {
		return "", nil, nil, errors.New("message is not a string")
	}
	message, ok := messageRaw.(string)
	if !ok {
		return "", nil, nil, errors.New("message is not a string")
	}

	post := &model.Post{
		UserId:    f.UserID,
		ChannelId: channel.Id,
		Message:   message,
	}
	err = fm.client.Post.CreatePost(post)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "failed to create announcement post")
	}

	return stepAnnouncementConfirmation, flow.State{
		"ChannelName": channel.Name,
	}, nil, nil
}

func (fm *FlowManager) stepDone() flow.Step {
	return flow.NewStep(stepDone).
		Terminal().
		WithText(":tada: You successfully installed Confluence.").
		OnRender(fm.onDone)
}

func (fm *FlowManager) onDone(f *flow.Flow) {
	fm.trackCompleteSetupWizard(f.UserID)
}
