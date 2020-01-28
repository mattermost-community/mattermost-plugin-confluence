package config

import (
	"github.com/mattermost/mattermost-server/plugin"
	"go.uber.org/atomic"
)

const (
	CommandPrefix             = PluginName
	URLMappingKeyPrefix       = "url_"
	ConfluenceSubscriptionKeyPrefix = "confluence_subscription"
	ServerExeToWebappRootPath = "/../webapp"

	URLPluginBase = "/plugins/" + PluginName
	URLStaticBase = URLPluginBase + "/static"

	HeaderMattermostUserID = "Mattermost-User-Id"
)

var (
	config     atomic.Value
	Mattermost plugin.API
	BotUserID  string
)

type Configuration struct {
}

func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

func SetConfig(c *Configuration) {
	config.Store(c)
}

func (c *Configuration) ProcessConfiguration() error {
	// any post-processing on configurations goes here

	return nil
}

func (c *Configuration) IsValid() error {
	// Add config validations here.
	// Check for required fields, formats, etc.

	return nil
}
