package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-plugin-confluence/server/util/types"
)

const (
	PathCurrentUser            = "/rest/api/user/current"
	PathCommentData            = "/rest/api/content/%s?expand=body.view,container,space,history"
	PathPageData               = "/rest/api/content/%s?status=any&expand=body.view,container,space,history"
	PathSpaceData              = "/rest/api/space/%s?status=any"
	PathAllSpaces              = "/rest/api/space"
	PathGetUserGroupsForServer = "/rest/api/user/memberof?username=%s"
	PathAdminData              = "/rest/api/audit"
)

const (
	Comment = "comment"
	Space   = "space"
	Page    = "page"
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

type SpaceResponse struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Links Links  `json:"_links"`
}

type CommentContainer struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Links Links  `json:"_links"`
}

type Links struct {
	Self string `json:"webui"`
}

type View struct {
	Value string `json:"value"`
}
type Body struct {
	View View `json:"view"`
}

type CreatedBy struct {
	Username string `json:"username"`
}

type History struct {
	CreatedBy CreatedBy `json:"createdBy"`
}

type CommentResponse struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Space     SpaceResponse    `json:"space"`
	Container CommentContainer `json:"container"`
	Body      Body             `json:"body"`
	Links     Links            `json:"_links"`
	History   History          `json:"history"`
}

type PageResponse struct {
	ID      string        `json:"id"`
	Title   string        `json:"title"`
	Space   SpaceResponse `json:"space"`
	Body    Body          `json:"body"`
	Links   Links         `json:"_links"`
	History History       `json:"history"`
}

type ConfluenceServerEvent struct {
	Comment *CommentResponse
	Page    *PageResponse
	Space   *SpaceResponse
	BaseURL string
}

func newServerClient(url string, httpClient *http.Client) Client {
	return &confluenceServerClient{
		URL:        url,
		HTTPClient: httpClient,
	}
}

func (csc *confluenceServerClient) GetSelf() (*types.ConfluenceUser, error) {
	confluenceServerUser := &ConfluenceServerUser{}
	if _, err := service.CallJSONWithURL(csc.URL, PathCurrentUser, http.MethodGet, nil, confluenceServerUser, csc.HTTPClient); err != nil {
		return nil, errors.Wrap(err, "confluence GetSelf. Error getting the current user")
	}

	confluenceUser := &types.ConfluenceUser{
		AccountID:   confluenceServerUser.UserKey,
		Name:        confluenceServerUser.Username,
		DisplayName: confluenceServerUser.DisplayName,
	}

	return confluenceUser, nil
}

func (csc *confluenceServerClient) GetEventData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*ConfluenceServerEvent, error) {
	var confluenceServerEvent ConfluenceServerEvent
	var err error
	if strings.Contains(webhookPayload.Event, Comment) {
		confluenceServerEvent.Comment, err = csc.GetCommentData(webhookPayload)
		if err != nil {
			return nil, errors.Wrap(err, "confluence GetEventData")
		}
	}
	if strings.Contains(webhookPayload.Event, Page) {
		confluenceServerEvent.Page, err = csc.GetPageData(int(webhookPayload.Page.ID))
		if err != nil {
			return nil, errors.Wrap(err, "confluence GetEventData")
		}
	}
	if strings.Contains(webhookPayload.Event, Space) {
		confluenceServerEvent.Space, err = csc.GetSpaceData(webhookPayload.Space.SpaceKey)
		if err != nil {
			return nil, errors.Wrap(err, "confluence GetEventData")
		}
	}

	return &confluenceServerEvent, nil
}

func (csc *confluenceServerClient) GetCommentData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*CommentResponse, error) {
	commentResponse := &CommentResponse{}
	_, err := service.CallJSONWithURL(csc.URL, fmt.Sprintf(PathCommentData, strconv.FormatInt(webhookPayload.Comment.ID, 10)), http.MethodGet, nil, commentResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetCommentData")
	}

	commentResponse.Body.View.Value = util.GetBodyForExcerpt(commentResponse.Body.View.Value)

	return commentResponse, nil
}

func (csc *confluenceServerClient) GetPageData(pageID int) (*PageResponse, error) {
	pageResponse := &PageResponse{}
	_, err := service.CallJSONWithURL(csc.URL, fmt.Sprintf(PathPageData, strconv.Itoa(pageID)), http.MethodGet, nil, pageResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetPageData")
	}

	pageResponse.Body.View.Value = util.GetBodyForExcerpt(pageResponse.Body.View.Value)

	return pageResponse, nil
}

func (csc *confluenceServerClient) GetSpaceData(spaceKey string) (*SpaceResponse, error) {
	spaceResponse := &SpaceResponse{}
	_, err := service.CallJSONWithURL(csc.URL, fmt.Sprintf(PathSpaceData, spaceKey), http.MethodGet, nil, spaceResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSpaceData")
	}

	return spaceResponse, nil
}

func (csc *confluenceServerClient) GetSpaceKeyFromSpaceID(spaceID int64) (string, error) {
	const pageSize = 10
	start := 0

	for {
		path := fmt.Sprintf("%s?start=%d&limit=%d", PathAllSpaces, start, pageSize)

		type apiResponse struct {
			Results []struct {
				ID   int64  `json:"id"`
				Key  string `json:"key"`
				Name string `json:"name"`
			} `json:"results"`
			Size int `json:"size"`
		}

		response := &apiResponse{}

		_, err := service.CallJSONWithURL(csc.URL, path, http.MethodGet, nil, response, csc.HTTPClient)
		if err != nil {
			return "", errors.Wrap(err, "confluence GetSpaceKeyFromSpaceID")
		}

		for _, space := range response.Results {
			if space.ID == spaceID {
				return space.Key, nil
			}
		}

		if len(response.Results) < pageSize {
			break
		}

		start += pageSize
	}

	return "", fmt.Errorf("confluence GetSpaceKeyFromSpaceID: no space found for the space id %v", spaceID)
}

func (csc *confluenceServerClient) GetUserGroups(connection *types.Connection) ([]*types.UserGroup, error) {
	userGroups := types.UserGroups{}

	_, err := service.CallJSONWithURL(csc.URL, PathGetUserGroupsForServer, http.MethodGet, nil, &userGroups, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetUserGroups")
	}

	return userGroups.Groups, nil
}
