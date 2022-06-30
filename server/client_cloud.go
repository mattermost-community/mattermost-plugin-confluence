package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/utils"
)

const (
	PathAccessibleResources   = "/oauth/token/accessible-resources"
	PathGetUserGroupsForCloud = "/rest/api/user/memberof?accountId=%s"
	PathGetSpacesForCloud     = "/rest/api/space"
	PathCreatePageForCloud    = "/rest/api/content"
	AccountID                 = "accountId"
)

type confluenceCloudClient struct {
	URL        string
	InstanceID string
	HTTPClient *http.Client
}

type ConfluenceCloudEvent struct {
	Comment *CommentResponse
	Page    *PageResponse
	Space   *SpaceResponse
}

type AccessibleResources struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	URL       string   `json:"url"`
	Scopes    []string `json:"scopes"`
	AvatarURL string   `json:"avatarUrl"`
}

type ConfluenceCloudUser struct {
	AccountID   string `json:"accountId"`
	PublicName  string `json:"publicName"`
	DisplayName string `json:"displayName"`
}

func newCloudClient(url, instanceID string, httpClient *http.Client) Client {
	return &confluenceCloudClient{
		URL:        url,
		InstanceID: instanceID,
		HTTPClient: httpClient,
	}
}

func (ccc *confluenceCloudClient) GetSelf() (*ConfluenceUser, int, error) {
	confluenceCloudUser := &ConfluenceCloudUser{}
	url, err := utils.GetEndpointURL(ccc.URL, PathCurrentUser)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSelf")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, confluenceCloudUser, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSelf")
	}

	confluenceUser := &ConfluenceUser{
		AccountID:   confluenceCloudUser.AccountID,
		Name:        confluenceCloudUser.PublicName,
		DisplayName: confluenceCloudUser.DisplayName,
	}
	return confluenceUser, statusCode, nil
}

func (ccc *confluenceCloudClient) GetCloudID() (string, int, error) {
	accessibleResources := []*AccessibleResources{}
	url, err := utils.GetEndpointURL(ccc.URL, PathAccessibleResources)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(err, "confluence GetCloudID")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, &accessibleResources, ccc.HTTPClient)
	if err != nil {
		return "", statusCode, errors.Wrap(err, "confluence GetCloudID")
	}

	for _, accessibleResource := range accessibleResources {
		if accessibleResource.URL == ccc.InstanceID {
			return accessibleResource.ID, statusCode, nil
		}
	}

	return "", statusCode, fmt.Errorf("cloudID not found for cloud instance: %s", ccc.InstanceID)
}

func (ccc *confluenceCloudClient) GetEventData(confluenceCloudEvent *serializer.ConfluenceCloudEvent, eventType string) (*ConfluenceCloudEvent, int, error) {
	var err error
	var statusCode int
	var confluenceCloudEventResponse ConfluenceCloudEvent
	if strings.Contains(eventType, Comment) {
		confluenceCloudEventResponse.Comment, statusCode, err = ccc.GetCommentData(confluenceCloudEvent)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}

		confluenceCloudEvent.Comment.Body.View.Value = confluenceCloudEventResponse.Comment.Body.View.Value
		confluenceCloudEvent.Comment.UserName = confluenceCloudEventResponse.Comment.History.CreatedBy.Username
	}
	if strings.Contains(eventType, Page) {
		confluenceCloudEventResponse.Page, statusCode, err = ccc.GetPageData(confluenceCloudEvent.Page.ID)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}

		confluenceCloudEvent.Page.Body.View.Value = confluenceCloudEventResponse.Page.Body.View.Value
		confluenceCloudEvent.Page.UserName = confluenceCloudEventResponse.Page.History.CreatedBy.Username
	}
	if strings.Contains(eventType, Space) {
		confluenceCloudEventResponse.Space, statusCode, err = ccc.GetSpaceData(confluenceCloudEvent.Space.SpaceKey)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}
	}

	return &confluenceCloudEventResponse, statusCode, nil
}

func (ccc *confluenceCloudClient) GetCommentData(confluenceCloudEvent *serializer.ConfluenceCloudEvent) (*CommentResponse, int, error) {
	commentResponse := &CommentResponse{}
	url, err := utils.GetEndpointURL(ccc.URL, fmt.Sprintf(PathCommentData, strconv.Itoa(confluenceCloudEvent.Comment.ID)))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetCommentData")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, commentResponse, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetCommentData")
	}
	commentResponse.Body.View.Value = utils.GetBodyForExcerpt(commentResponse.Body.View.Value)

	return commentResponse, statusCode, nil
}

func (ccc *confluenceCloudClient) GetPageData(pageID int) (*PageResponse, int, error) {
	pageResponse := &PageResponse{}
	url, err := utils.GetEndpointURL(ccc.URL, fmt.Sprintf(PathPageData, strconv.Itoa(pageID)))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetPageData")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, pageResponse, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetPageData")
	}
	pageResponse.Body.View.Value = utils.GetBodyForExcerpt(pageResponse.Body.View.Value)

	return pageResponse, statusCode, nil
}

func (ccc *confluenceCloudClient) GetSpaceData(spaceKey string) (*SpaceResponse, int, error) {
	spaceResponse := &SpaceResponse{}
	url, err := utils.GetEndpointURL(ccc.URL, fmt.Sprintf(PathSpaceData, spaceKey))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaceData")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, spaceResponse, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSpaceData")
	}

	return spaceResponse, statusCode, nil
}

func (ccc *confluenceCloudClient) GetUserGroups(connection *Connection) ([]*UserGroup, int, error) {
	userGroups := UserGroups{}
	url, err := utils.GetEndpointURL(ccc.URL, PathGetUserGroupsForCloud)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetUserGroups")
	}

	url, err = utils.AddQueryParams(url, map[string]interface{}{
		AccountID: connection.AccountID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, &userGroups, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetUserGroups")
	}

	return userGroups.Groups, statusCode, nil
}

func (ccc *confluenceCloudClient) GetSpaces() ([]*Spaces, int, error) {
	spacesForConfluenceURL := SpacesForConfluenceURL{}
	url, err := utils.GetEndpointURL(ccc.URL, PathGetSpacesForCloud)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	url, err = utils.AddQueryParams(url, map[string]interface{}{
		Limit: SpaceLimit,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodGet, url, nil, &spacesForConfluenceURL, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSpaces")
	}
	return spacesForConfluenceURL.Spaces, statusCode, nil
}

func (ccc *confluenceCloudClient) CreatePage(spaceKey string, pageDetails *serializer.PageDetails) (*CreatePageResponse, int, error) {
	requestBody := GetRequestBodyForCreatePage(spaceKey, pageDetails)
	createPageResponse := &CreatePageResponse{}
	url, err := utils.GetEndpointURL(ccc.URL, PathCreatePageForCloud)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence CreatePage")
	}
	_, statusCode, err := utils.CallJSON(ccc.URL, http.MethodPost, url, requestBody, createPageResponse, ccc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence CreatePage")
	}
	return createPageResponse, http.StatusOK, nil
}
