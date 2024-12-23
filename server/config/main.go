package config

import (
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

var (
	config     atomic.Value
	Mattermost plugin.API
	BotUserID  string
)

type Configuration struct {
	Secret                      string `json:"Secret"`
	ConfluenceOAuthClientID     string
	ConfluenceOAuthClientSecret string
	ConfluenceURL               string
}

func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

func SetConfig(c *Configuration) {
	config.Store(c)
}

func (c *Configuration) ProcessConfiguration() error {
	c.Secret = strings.TrimSpace(c.Secret)

	return nil
}

func (c *Configuration) IsValid() error {
	if c.Secret == "" {
		return errors.New("please provide the Webhook Secret")
	}

	return nil
}

func (c *Configuration) Sanitize() {
	// Ensure Confluence ends with a slash
	c.ConfluenceURL = strings.TrimRight(c.ConfluenceURL, "/")

	c.ConfluenceOAuthClientID = strings.TrimSpace(c.ConfluenceOAuthClientID)
	c.ConfluenceOAuthClientSecret = strings.TrimSpace(c.ConfluenceOAuthClientSecret)
}

func (c *Configuration) IsOAuthConfigured() bool {
	return (c.ConfluenceOAuthClientID != "" && c.ConfluenceOAuthClientSecret != "")
}

func (c *Configuration) ToMap() (map[string]interface{}, error) {
	var out map[string]interface{}
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
