package service

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const (
	pageCreateMessage    = "A new page [%s](%s) was created in the **%s** space."
	commentCreateMessage = "A new [comment](%s) was posted on the [%s](%s) page."
	pageUpdateMessage    = "Page [%s](%s) was updated in the **%s** space."
	commentUpdateMessage = "A [comment](%s) was updated on the [%s](%s) page."
	pageDeleteMessage    = "Page **%s** was removed from the **%s** space."
	commentDeleteMessage = "A comment was deleted from the [%s](%s) page."
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
	case serializer.PageCreatedEvent:
		message = fmt.Sprintf(pageCreateMessage, page.Title, page.Self, page.SpaceKey)
	case serializer.CommentCreatedEvent:
		message = fmt.Sprintf(commentCreateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case serializer.PageUpdatedEvent:
		message = fmt.Sprintf(pageUpdateMessage, page.Title, page.Self, page.SpaceKey)
	case serializer.CommentUpdatedEvent:
		message = fmt.Sprintf(commentUpdateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case serializer.PageRemovedEvent:
		message = fmt.Sprintf(pageDeleteMessage, page.Title, page.SpaceKey)
	case serializer.CommentRemovedEvent:
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
