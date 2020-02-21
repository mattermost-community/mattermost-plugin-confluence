package serializer

import (
	"github.com/mattermost/mattermost-server/model"
)

type ConfluenceEvent interface {
	GetNotificationPost(string) *model.Post
	GetURL() string
	GetSpaceKey() string
	GetPageID() string
}
