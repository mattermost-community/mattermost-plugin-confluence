package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/telemetry"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
)

const (
	botUserName    = "confluence"
	botDisplayName = "Confluence"
	botDescription = "Bot for confluence plugin."
)

type Plugin struct {
	plugin.MattermostPlugin
	client *pluginapi.Client

	BotUserID string

	Router *mux.Router

	flowManager *FlowManager

	// templates are loaded on startup
	templates map[string]*template.Template

	pluginVersion string

	ConfluenceClient *Client

	telemetryClient telemetry.Client
	tracker         telemetry.Tracker
}

func (p *Plugin) OnActivate() error {
	config.Mattermost = p.API
	p.client = pluginapi.NewClient(p.API, p.Driver)

	if err := p.setUpBotUser(); err != nil {
		config.Mattermost.LogError("Failed to create a bot user", "Error", err.Error())
		return err
	}

	p.Router = p.InitAPI()
	p.initializeTelemetry()

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "couldn't get bundle path")
	}

	templates, err := p.loadTemplates(filepath.Join(bundlePath, "assets", "templates"))
	if err != nil {
		return err
	}
	p.templates = templates

	p.pluginVersion = manifest.Version

	flowManager, err := p.NewFlowManager()
	if err != nil {
		return errors.Wrap(err, "failed to create flow manager")
	}
	p.flowManager = flowManager

	cmd, err := GetCommand(p.API)
	if err != nil {
		return errors.Wrap(err, "failed to get command")
	}

	if err := p.API.RegisterCommand(cmd); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	// If OnActivate has not been run yet.
	if config.Mattermost == nil {
		return nil
	}
	var configuration config.Configuration

	if err := config.Mattermost.LoadPluginConfiguration(&configuration); err != nil {
		config.Mattermost.LogError("Error in LoadPluginConfiguration.", "Error", err.Error())
		return err
	}

	if err := configuration.ProcessConfiguration(); err != nil {
		config.Mattermost.LogError("Error in ProcessConfiguration.", "Error", err.Error())
		return err
	}

	if err := configuration.IsValid(); err != nil {
		config.Mattermost.LogError("Error in Validating Configuration.", "Error", err.Error())
		return err
	}

	config.SetConfig(&configuration)
	return nil
}

func (p *Plugin) setUpBotUser() error {
	botUserID, err := p.client.Bot.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: botDescription,
	})
	if err != nil {
		config.Mattermost.LogError("Error in setting up bot user", "Error", err.Error())
		return err
	}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return err
	}

	profileImage, err := os.ReadFile(filepath.Join(bundlePath, "assets", "icon.png"))
	if err != nil {
		return err
	}

	if appErr := p.API.SetProfileImage(botUserID, profileImage); appErr != nil {
		return errors.Wrap(appErr, "couldn't set profile image")
	}

	config.BotUserID = botUserID
	p.BotUserID = botUserID
	return nil
}

func (p *Plugin) ExecuteCommand(context *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args, argErr := util.SplitArgs(commandArgs.Command)
	if argErr != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         argErr.Error(),
		}, nil
	}
	return ConfluenceCommandHandler.Handle(p, commandArgs, args[1:]...), nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)

	conf := config.GetConfig()
	if err := conf.IsValid(); err != nil {
		p.API.LogError("This plugin is not configured.", "Error", err.Error())
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	p.Router.ServeHTTP(w, r)
}

func (p *Plugin) debugf(f string, args ...interface{}) {
	p.API.LogDebug(fmt.Sprintf(f, args...))
}

func (p *Plugin) errorf(f string, args ...interface{}) {
	p.API.LogError(fmt.Sprintf(f, args...))
}

func (p *Plugin) track(name, userID string) {
	p.trackWithArgs(name, userID, nil)
}

func (p *Plugin) trackWithArgs(name, userID string, args map[string]interface{}) {
	if args == nil {
		args = map[string]interface{}{}
	}
	args["time"] = model.GetMillis()
	_ = p.tracker.TrackUserEvent(name, userID, args)
}

func main() {
	plugin.ClientMain(&Plugin{})
}
