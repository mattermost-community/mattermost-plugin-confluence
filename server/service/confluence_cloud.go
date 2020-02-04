package service

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-server/model"
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
		message = fmt.Sprintf("Page [%s](%s) created in space **%s**", page.Title, page.Self, page.SpaceKey)
	case "comment_created":
		message = fmt.Sprintf("[Comment](%s) created in page [%s](%s)", comment.Self, comment.Parent.Title, comment.Parent.Self)
	case "page_update":
		message = fmt.Sprintf("Page [%s](%s) updated in space **%s**", page.Title, page.Self, page.SpaceKey)
	case "comment_update":
		message = fmt.Sprintf("[Comment](%s) updated in page [%s](%s)", comment.Self, comment.Parent.Title, comment.Parent.Self)
	case "page_removed":
		message = fmt.Sprintf("Page [%s](%s) created in space **%s**", page.Title, page.Self, page.SpaceKey)
	case "comment_delete":
		message = fmt.Sprintf("Comment deleted in page [%s](%s)", comment.Parent.Title, comment.Parent.Self)
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
