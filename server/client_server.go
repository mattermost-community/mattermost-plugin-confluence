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
	PathCurrentUser            = "/rest/api/user/current"
	PathCreateWebhook          = "/rest/api/webhooks"
	PathCommentData            = "/rest/api/content/%s?expand=body.view,container,space,history"
	PathPageData               = "/rest/api/content/%s?status=any&expand=body.view,container,space,history"
	PathSpaceData              = "/rest/api/space/%s?status=any"
	PathDeleteWebhook          = "/rest/api/webhooks/%s"
	PathGetUserGroupsForServer = "/rest/api/user/memberof"
	PathAdminData              = "/rest/api/audit/retention"
	PathGetSpacesForServer     = "/rest/api/space"
	PathCreatePageForServer    = "/rest/api/content"
)

const (
	Comment    = "comment"
	Space      = "space"
	Page       = "page"
	Limit      = "limit"
	UserName   = "username"
	SpaceLimit = 100
)

var webhookEvents = []string{"space_created", "space_removed", "space_updated",
	"page_created", "page_removed", "page_restored", "page_trashed",
	"page_updated", "comment_created", "comment_removed", "comment_updated"}

type WebhookConfiguration struct {
	Secret string `json:"secret"`
}

type CreateWebhookRequestBody struct {
	Name          string                `json:"name"`
	Events        []string              `json:"events"`
	URL           string                `json:"url"`
	Active        string                `json:"active"`
	Configuration *WebhookConfiguration `json:"configuration"`
}

type PageCreateSpace struct {
	Key string `json:"key"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type PageBody struct {
	Storage Storage `json:"storage"`
}

type PageRequestBody struct {
	Title string          `json:"title"`
	Type  string          `json:"type"`
	Space PageCreateSpace `json:"space"`
	Body  PageBody        `json:"body"`
}

type CreatePageResponse struct {
	Space SpaceResponse `json:"space"`
	Links PageLinks     `json:"_links"`
}

type PageLinks struct {
	Self    string `json:"webui"`
	BaseURL string `json:"base"`
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

type confluenceServerClient struct {
	URL        string
	HTTPClient *http.Client
}

type ConfluenceServerUser struct {
	UserKey     string `json:"userKey"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type AdminData struct {
	Number int    `json:"number"`
	Units  string `json:"units"`
}

type WebhookResponse struct {
	ID int `json:"id"`
}

func newServerClient(url string, httpClient *http.Client) Client {
	return &confluenceServerClient{
		URL:        url,
		HTTPClient: httpClient,
	}
}

func (csc *confluenceServerClient) CreateWebhook(subscription serializer.Subscription, redirectURL, secret string) (*WebhookResponse, int, error) {
	requestBody := &CreateWebhookRequestBody{
		Name:   subscription.GetAlias(),
		Events: webhookEvents,
		URL:    redirectURL,
		Active: "true",
		Configuration: &WebhookConfiguration{
			Secret: secret,
		},
	}
	webhookResponse := &WebhookResponse{}
	url, err := utils.GetEndpointURL(csc.URL, PathCreateWebhook)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence CreateWebhook")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodPost, url, requestBody, webhookResponse, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence CreateWebhook")
	}
	return webhookResponse, statusCode, nil
}

func (csc *confluenceServerClient) DeleteWebhook(webhookID string) (int, error) {
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathDeleteWebhook, webhookID))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "confluence DeleteWebhook")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodDelete, url, nil, nil, csc.HTTPClient)
	if err != nil && err.Error() != utils.ErrorStatusNotFound {
		return statusCode, errors.Wrap(err, "confluence DeleteWebhook")
	}
	return statusCode, nil
}

func (csc *confluenceServerClient) GetSelf() (*ConfluenceUser, int, error) {
	confluenceServerUser := &ConfluenceServerUser{}
	url, err := utils.GetEndpointURL(csc.URL, PathCurrentUser)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSelf")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, confluenceServerUser, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSelf")
	}

	confluenceUser := &ConfluenceUser{
		AccountID:   confluenceServerUser.UserKey,
		Name:        confluenceServerUser.Username,
		DisplayName: confluenceServerUser.DisplayName,
	}

	return confluenceUser, statusCode, nil
}

func (csc *confluenceServerClient) GetEventData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*ConfluenceServerEvent, int, error) {
	var confluenceServerEvent ConfluenceServerEvent
	var err error
	var statusCode int
	if strings.Contains(webhookPayload.Event, Comment) {
		confluenceServerEvent.Comment, statusCode, err = csc.GetCommentData(webhookPayload)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}
	}
	if strings.Contains(webhookPayload.Event, Page) {
		confluenceServerEvent.Page, statusCode, err = csc.GetPageData(webhookPayload)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}
	}
	if strings.Contains(webhookPayload.Event, Space) {
		confluenceServerEvent.Space, statusCode, err = csc.GetSpaceData(webhookPayload.Space.SpaceKey)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "confluence GetEventData")
		}
	}

	return &confluenceServerEvent, statusCode, nil
}

func (csc *confluenceServerClient) GetCommentData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*CommentResponse, int, error) {
	commentResponse := &CommentResponse{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathCommentData, strconv.FormatInt(webhookPayload.Comment.ID, 10)))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetCommentData")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, commentResponse, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetCommentData")
	}
	commentResponse.Body.View.Value = utils.GetBodyForExcerpt(commentResponse.Body.View.Value)

	return commentResponse, statusCode, nil
}

func (csc *confluenceServerClient) GetPageData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*PageResponse, int, error) {
	pageResponse := &PageResponse{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathPageData, strconv.FormatInt(webhookPayload.Page.ID, 10)))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetPageData")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, pageResponse, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetPageData")
	}
	pageResponse.Body.View.Value = utils.GetBodyForExcerpt(pageResponse.Body.View.Value)

	return pageResponse, statusCode, nil
}

func (csc *confluenceServerClient) GetSpaceData(spaceKey string) (*SpaceResponse, int, error) {
	spaceResponse := &SpaceResponse{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathSpaceData, spaceKey))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaceData")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, spaceResponse, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSpaceData")
	}

	return spaceResponse, statusCode, nil
}

func (csc *confluenceServerClient) CheckConfluenceAdmin() (*AdminData, int, error) {
	adminData := &AdminData{}
	url, err := utils.GetEndpointURL(csc.URL, PathAdminData)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence CheckConfluenceAdmin")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, adminData, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence CheckConfluenceAdmin")
	}

	return adminData, statusCode, nil
}

func (csc *confluenceServerClient) GetUserGroups(connection *Connection) ([]*UserGroup, int, error) {
	userGroups := UserGroups{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathGetUserGroupsForServer, connection.Name))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetUserGroups")
	}
	url, err = utils.AddQueryParams(url, map[string]interface{}{
		UserName: connection.Name,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, &userGroups, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetUserGroups")
	}
	return userGroups.Groups, statusCode, nil
}

func (csc *confluenceServerClient) GetSpaces() ([]*Spaces, int, error) {
	spacesForConfluenceURL := SpacesForConfluenceURL{}
	url, err := utils.GetEndpointURL(csc.URL, PathGetSpacesForServer)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	url, err = utils.AddQueryParams(url, map[string]interface{}{
		Limit: SpaceLimit,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence GetSpaces")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodGet, url, nil, &spacesForConfluenceURL, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence GetSpaces")
	}

	return spacesForConfluenceURL.Spaces, statusCode, nil
}

func (csc *confluenceServerClient) CreatePage(spaceKey string, pageDetails *serializer.PageDetails) (*CreatePageResponse, int, error) {
	requestBody := GetRequestBodyForCreatePage(spaceKey, pageDetails)
	createPageResponse := &CreatePageResponse{}
	url, err := utils.GetEndpointURL(csc.URL, PathCreatePageForServer)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "confluence CreatePage")
	}
	_, statusCode, err := utils.CallJSON(csc.URL, http.MethodPost, url, requestBody, createPageResponse, csc.HTTPClient)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "confluence CreatePage")
	}
	return createPageResponse, http.StatusOK, nil
}
