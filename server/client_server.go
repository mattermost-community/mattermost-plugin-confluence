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
	PathGetUserGroupsForServer = "/rest/api/user/memberof?username=%s"
	PathAdminData              = "/rest/api/audit/retention"
	PathGetSpacesForServer     = "/rest/api/space?limit=100"
	PathCreatePage             = "/rest/api/content"
)

const (
	Comment = "comment"
	Space   = "space"
	Page    = "page"
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
type SpaceForPageCreate struct {
	Key string `json:"key"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type BodyForPageCreate struct {
	Storage Storage `json:"storage"`
}

type CreatePageRequestBody struct {
	Title string             `json:"title"`
	Type  string             `json:"type"`
	Space SpaceForPageCreate `json:"space"`
	Body  BodyForPageCreate  `json:"body"`
}

type CreatePageResponse struct {
	Space SpaceResponse      `json:"space"`
	Links LinksForPageCreate `json:"_links"`
}

type LinksForPageCreate struct {
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

func (csc *confluenceServerClient) CreateWebhook(subscription serializer.Subscription, redirectURL, secret string) (*WebhookResponse, error) {
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
		return nil, errors.Wrap(err, "confluence CreateWebhook")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodPost, url, requestBody, webhookResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CreateWebhook")
	}
	return webhookResponse, nil
}

func (csc *confluenceServerClient) DeleteWebhook(webhookID string) error {
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathDeleteWebhook, webhookID))
	if err != nil {
		return errors.Wrap(err, "confluence DeleteWebhook")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodDelete, url, nil, nil, csc.HTTPClient)
	if err != nil && err.Error() != utils.ErrorStatusNotFound {
		return errors.Wrap(err, "confluence DeleteWebhook")
	}
	return nil
}

func (csc *confluenceServerClient) GetSelf() (*ConfluenceUser, error) {
	confluenceServerUser := &ConfluenceServerUser{}
	url, err := utils.GetEndpointURL(csc.URL, PathCurrentUser)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSelf")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, confluenceServerUser, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSelf")
	}

	confluenceUser := &ConfluenceUser{
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
		confluenceServerEvent.Page, err = csc.GetPageData(webhookPayload)
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
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathCommentData, strconv.FormatInt(webhookPayload.Comment.ID, 10)))
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetCommentData")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, commentResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetCommentData")
	}
	commentResponse.Body.View.Value = utils.GetBodyForExcerpt(commentResponse.Body.View.Value)

	return commentResponse, nil
}

func (csc *confluenceServerClient) GetPageData(webhookPayload *serializer.ConfluenceServerWebhookPayload) (*PageResponse, error) {
	pageResponse := &PageResponse{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathPageData, strconv.FormatInt(webhookPayload.Page.ID, 10)))
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetPageData")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, pageResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetPageData")
	}
	pageResponse.Body.View.Value = utils.GetBodyForExcerpt(pageResponse.Body.View.Value)

	return pageResponse, nil
}

func (csc *confluenceServerClient) GetSpaceData(spaceKey string) (*SpaceResponse, error) {
	spaceResponse := &SpaceResponse{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathSpaceData, spaceKey))
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSpaceData")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, spaceResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSpaceData")
	}

	return spaceResponse, nil
}

func (csc *confluenceServerClient) CheckConfluenceAdmin() (*AdminData, error) {
	adminData := &AdminData{}
	url, err := utils.GetEndpointURL(csc.URL, PathAdminData)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CheckConfluenceAdmin")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, adminData, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CheckConfluenceAdmin")
	}

	return adminData, nil
}

func (csc *confluenceServerClient) GetUserGroups(connection *Connection) ([]*UserGroup, error) {
	userGroups := UserGroups{}
	url, err := utils.GetEndpointURL(csc.URL, fmt.Sprintf(PathGetUserGroupsForServer, connection.Name))
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetUserGroups")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, &userGroups, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetUserGroups")
	}
	return userGroups.Groups, nil
}

func (csc *confluenceServerClient) GetSpacesForConfluenceURL() ([]*SpaceForConfluenceURL, error) {
	spacesForConfluenceURL := SpacesForConfluenceURL{}
	url, err := utils.GetEndpointURL(csc.URL, PathGetSpacesForServer)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSpacesForConfluenceURL")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodGet, url, nil, &spacesForConfluenceURL, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence GetSpacesForConfluenceURL")
	}
	return spacesForConfluenceURL.Spaces, nil
}

func (csc *confluenceServerClient) CreatePage(spaceKey string, pageDetails *serializer.PageDetails) (*CreatePageResponse, error) {
	requestBody := &CreatePageRequestBody{
		Title: pageDetails.Title,
		Type:  "page",
		Space: SpaceForPageCreate{
			Key: spaceKey,
		},
		Body: BodyForPageCreate{
			Storage: Storage{
				Value:          pageDetails.Description,
				Representation: "storage",
			},
		},
	}
	createPageResponse := &CreatePageResponse{}
	url, err := utils.GetEndpointURL(csc.URL, PathCreatePage)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CreatePage")
	}
	_, err = utils.CallJSON(csc.URL, http.MethodPost, url, requestBody, createPageResponse, csc.HTTPClient)
	if err != nil {
		return nil, errors.Wrap(err, "confluence CreatePage")
	}
	return createPageResponse, nil
}
