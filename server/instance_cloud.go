package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-confluence/server/utils"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils/types"
)

type cloudInstance struct {
	CloudID string

	*InstanceCommon
}

func (p *Plugin) installCloudInstance(rawURL string) (string, *cloudInstance, error) {
	confluenceURL, err := utils.NormalizeConfluenceURL(rawURL)
	if err != nil {
		return "", nil, err
	}
	if !strings.HasPrefix(confluenceURL, "https://") {
		return "", nil, errors.New("a secure HTTPS URL is required")
	}

	instance := &cloudInstance{
		InstanceCommon: newInstanceCommon(p, CloudInstanceType, types.ID(confluenceURL)),
	}

	err = p.InstallInstance(instance)
	if err != nil {
		return "", nil, err
	}

	return confluenceURL, instance, err
}

func (ci *cloudInstance) GetOAuth2Config(isAdmin bool) (*oauth2.Config, error) {
	config, ok := ci.Plugin.getConfig().ParsedConfluenceConfig[ci.GetURL()]
	if !ok {
		return nil, fmt.Errorf(configNotFoundError, ci.InstanceID, ci.InstanceID)
	}

	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s%s", ci.GetPluginURL(), instancePath(routeUserComplete, ci.InstanceID)),
		Scopes: []string{
			"read:confluence-user",
			"read:confluence-content.all",
			"read:content-details:confluence",
			"write:confluence-content",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.atlassian.com/authorize?audience=api.atlassian.com",
			TokenURL: "https://auth.atlassian.com/oauth/token",
		},
	}, nil
}

func (ci *cloudInstance) GetURL() string {
	return ci.InstanceID.String()
}

func (ci *cloudInstance) GetManageAppsURL() string {
	return fmt.Sprintf("%s/plugins/servlet/upm", ci.GetURL())
}

func (ci *cloudInstance) GetManageWebhooksURL() string {
	return cloudManageWebhooksURL(ci.GetURL())
}

func cloudManageWebhooksURL(confluenceURL string) string {
	return fmt.Sprintf("%s/plugins/servlet/webhooks", confluenceURL)
}

func (ci *cloudInstance) GetClient(connection *Connection, mattermostUserID types.ID) (Client, error) {
	token, err := ci.Plugin.ParseAuthToken(connection.OAuth2Token)
	if err != nil {
		return nil, err
	}

	token, err = ci.checkAndRefreshToken(token, connection, ci.InstanceID, mattermostUserID)
	if err != nil {
		return nil, err
	}

	oconf, err := ci.GetOAuth2Config(connection.IsAdmin)
	if err != nil {
		return nil, err
	}
	httpClient := oconf.Client(context.Background(), token)

	baseURL := "https://api.atlassian.com"
	if ci.CloudID != "" {
		baseURL = fmt.Sprintf("%s/ex/confluence/%s/wiki", baseURL, ci.CloudID)
	}

	return newCloudClient(baseURL, ci.InstanceID.String(), httpClient), nil
}

func (ci *cloudInstance) checkAndRefreshToken(token *oauth2.Token, connection *Connection, instanceID types.ID, mattermostUserID types.ID) (*oauth2.Token, error) {
	// If there is only one minute left for the token to expire, we are refreshing the token.
	// We don't want the token to expire between the time when we decide that the old token is valid
	// and the time at which we create the request. We are handling that by not letting the token expire.
	if time.Until(token.Expiry) <= 1*time.Minute {
		oconf, err := ci.GetOAuth2Config(connection.IsAdmin)
		if err != nil {
			return nil, err
		}
		src := oconf.TokenSource(context.Background(), token)
		newToken, err := src.Token() // this actually goes and renews the tokens
		if err != nil {
			return nil, errors.Wrap(err, "unable to get the new refreshed token")
		}
		if newToken.AccessToken != token.AccessToken {
			encryptedToken, err := ci.Plugin.NewEncodedAuthToken(newToken)
			if err != nil {
				return nil, err
			}
			connection.OAuth2Token = encryptedToken

			err = ci.Plugin.userStore.StoreConnection(instanceID, mattermostUserID, connection)
			if err != nil {
				return nil, err
			}
			return newToken, nil
		}
	}

	return token, nil
}
