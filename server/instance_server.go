package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-confluence/server/utils"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils/types"
)

type serverInstance struct {
	*InstanceCommon
}

func (p *Plugin) installServerInstance(rawURL string) (string, *serverInstance, error) {
	confluenceURL, err := utils.CheckConfluenceURL(p.GetSiteURL(), rawURL, false)
	if err != nil {
		return "", nil, err
	}
	if utils.IsConfluenceCloudURL(confluenceURL) {
		return "", nil, errors.Errorf("`%s` is not a Confluence server URL, it refers to Confluence Cloud", confluenceURL)
	}

	instance := &serverInstance{
		InstanceCommon: newInstanceCommon(p, ServerInstanceType, types.ID(confluenceURL)),
	}

	err = p.InstallInstance(instance)
	if err != nil {
		return "", nil, err
	}

	return confluenceURL, instance, err
}

func (si *serverInstance) GetURL() string {
	return si.InstanceID.String()
}

func (si *serverInstance) GetManageAppsURL() string {
	return fmt.Sprintf("%s/plugins/servlet/applinks/listApplicationLinks", si.GetURL())
}

func (si *serverInstance) GetManageWebhooksURL() string {
	return fmt.Sprintf("%s/plugins/servlet/webhooks", si.GetURL())
}

func (si *serverInstance) GetOAuth2Config(isAdmin bool) (*oauth2.Config, error) {
	config, ok := si.Plugin.getConfig().ParsedConfluenceConfig[si.GetURL()]
	if !ok {
		return nil, fmt.Errorf(configNotFoundError, si.InstanceID, si.InstanceID)
	}

	var scopes []string
	if isAdmin {
		scopes = []string{
			"ADMIN",
		}
	} else {
		scopes = []string{
			"READ",
			"WRITE",
		}
	}
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s%s", si.GetPluginURL(), instancePath(routeUserComplete, si.InstanceID)),
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/rest/oauth2/latest/authorize", si.GetURL()),
			TokenURL: fmt.Sprintf("%s/rest/oauth2/latest/token", si.GetURL()),
		},
	}, nil
}

func (si *serverInstance) GetClient(connection *Connection, mattermostUserID types.ID) (Client, error) {
	token, err := si.Plugin.ParseAuthToken(connection.OAuth2Token)
	if err != nil {
		return nil, err
	}
	token, err = si.checkAndRefreshToken(token, connection, si.InstanceID, mattermostUserID)
	if err != nil {
		return nil, err
	}

	oconf, err := si.GetOAuth2Config(connection.IsAdmin)
	if err != nil {
		return nil, err
	}
	httpClient := oconf.Client(context.Background(), token)

	return newServerClient(si.GetURL(), httpClient), nil
}

func (si *serverInstance) checkAndRefreshToken(token *oauth2.Token, connection *Connection, instanceID types.ID, mattermostUserID types.ID) (*oauth2.Token, error) {
	// If there is only one minute left for the token to expire, we are refreshing the token.
	// We don't want the token to expire between the time when we decide that the old token is valid
	// and the time at which we create the request. We are handling that by not letting the token expire.
	if time.Until(token.Expiry) <= 1*time.Minute {
		oconf, err := si.GetOAuth2Config(connection.IsAdmin)
		if err != nil {
			return nil, err
		}
		src := oconf.TokenSource(context.Background(), token)
		newToken, err := src.Token() // this actually goes and renews the tokens
		if err != nil {
			return nil, errors.Wrap(err, "unable to get the new refreshed token")
		}
		if newToken.AccessToken != token.AccessToken {
			encryptedToken, err := si.Plugin.NewEncodedAuthToken(newToken)
			if err != nil {
				return nil, err
			}
			connection.OAuth2Token = encryptedToken

			err = si.Plugin.userStore.StoreConnection(instanceID, mattermostUserID, connection)
			if err != nil {
				return nil, err
			}
			return newToken, nil
		}
	}

	return token, nil
}
