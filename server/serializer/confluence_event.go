package serializer

import "github.com/mattermost/mattermost-server/v6/model"

type ConfluenceEvent interface {
	GetNotificationPost(string, string, string) *model.Post
	GetPageID() string
	GetSpaceKey() string
	GetURL() string
}
