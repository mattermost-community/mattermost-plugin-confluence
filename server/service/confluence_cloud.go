package service

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-server/model"
)

const (
	pageCreateMessage = "Page [%s](%s) created in space **%s**"
	commentCreateMessage= "[Comment](%s) created in page [%s](%s)"
	pageUpdateMessage = "Page [%s](%s) updated in space **%s**"
	commentUpdateMessage= "[Comment](%s) updated in page [%s](%s)"
	pageDeleteMessage = "Page [%s](%s) removed in space **%s**"
	commentDeleteMessage= "Comment deleted in page [%s](%s)"
)

func SendConfluenceCloudNotification(event *serializer.ConfluenceCloudEvent, eventType string) {
	post := generateConfluenceCloudNotificationPost(event, eventType)
	if post == nil {
		return
	}

	if event.Comment != nil {
		SendConfluenceNotifications(post, event.Comment.Self, event.Comment.SpaceKey, eventType)
	} else if event.Page != nil {
		SendConfluenceNotifications(post, event.Page.Self, event.Page.SpaceKey, eventType)
	}
}

func generateConfluenceCloudNotificationPost(event *serializer.ConfluenceCloudEvent, eventType string) *model.Post {
	message := ""
	page := event.Page
	comment := event.Comment
	switch eventType {
	case "page_created":
		message = fmt.Sprintf(pageCreateMessage, page.Title, page.Self, page.SpaceKey)
	case "comment_created":
		message = fmt.Sprintf(commentCreateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case "page_update":
		message = fmt.Sprintf(pageUpdateMessage, page.Title, page.Self, page.SpaceKey)
	case "comment_update":
		message = fmt.Sprintf(commentUpdateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case "page_removed":
		message = fmt.Sprintf(pageDeleteMessage, page.Title, page.Self, page.SpaceKey)
	case "comment_delete":
		message = fmt.Sprintf(commentDeleteMessage, comment.Parent.Title, comment.Parent.Self)
	default:
		return nil
	}

	post := &model.Post{
		UserId:  config.BotUserID,
		Type:    model.POST_DEFAULT,
		Message: message,
	}
	return post
}
