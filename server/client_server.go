package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/util/types"
)

const (
	PathCurrentUser = "/rest/api/user/current"
	PathAdminData   = "/rest/api/audit"
)

type confluenceServerClient struct {
	URL        string
	HTTPClient *http.Client
}

type ConfluenceServerUser struct {
	UserKey     string `json:"userKey"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

type AdminData struct {
	Number int    `json:"number"`
	Units  string `json:"units"`
}

func newServerClient(url string, httpClient *http.Client) Client {
	return &confluenceServerClient{
		URL:        url,
		HTTPClient: httpClient,
	}
}

func (csc *confluenceServerClient) GetSelf() (*types.ConfluenceUser, error) {
	confluenceServerUser := &ConfluenceServerUser{}
	_, err := service.CallJSONWithURL(csc.URL, PathCurrentUser, http.MethodGet, nil, confluenceServerUser, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSelf")
	}

	confluenceUser := &types.ConfluenceUser{
		AccountID:   confluenceServerUser.UserKey,
		Name:        confluenceServerUser.Username,
		DisplayName: confluenceServerUser.DisplayName,
	}

	return confluenceUser, nil
}

func (csc *confluenceServerClient) CheckConfluenceAdmin() (*AdminData, error) {
	adminData := &AdminData{}
	_, err := service.CallJSONWithURL(csc.URL, PathAdminData, http.MethodGet, nil, adminData, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CheckConfluenceAdmin")
	}

	return adminData, nil
}
