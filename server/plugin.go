package main

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/controller"
)

const (
	botUserName    = "confluence"
	botDisplayName = "Confluence"
	botDescription = "Bot for confluence plugin."
)

type Plugin struct {
	plugin.MattermostPlugin

	handler http.Handler
}

func (p *Plugin) OnActivate() error {
	config.Mattermost = p.API
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: botDescription,
	})
	if err != nil {
		config.Mattermost.LogError("Error in setting up bot user", "Error", err.Error())
		return err
	}
	config.BotUserID = botUserID

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.API.RegisterCommand(getCommand()); err != nil {
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

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	conf := config.GetConfig()
	if err := conf.IsValid(); err != nil {
		p.API.LogError("This plugin is not configured.", "Error", err.Error())
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	path := r.URL.Path
	endpoint := controller.Endpoints[path]

	if endpoint == nil {
		p.handler.ServeHTTP(w, r)
	} else if !endpoint.RequiresAuth || controller.Authenticated(w, r) {
		endpoint.Execute(w, r)
	}
}

func main() {
	plugin.ClientMain(&Plugin{})
}
