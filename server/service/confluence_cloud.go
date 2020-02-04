package service

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const (
	confluenceCloudPageCreateMessage    = "A new page [%s](%s) was created in the **%s** space."
	confluenceCloudPageUpdateMessage    = "Page [%s](%s) was updated in the **%s** space."
	confluenceCloudPageDeleteMessage    = "Page **%s** was removed from the **%s** space."
	confluenceCloudCommentCreateMessage = "A new [comment](%s) was posted on the [%s](%s) page."
	confluenceCloudCommentUpdateMessage = "A [comment](%s) was updated on the [%s](%s) page."
	confluenceCloudCommentDeleteMessage = "A comment was deleted from the [%s](%s) page."
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
		message = fmt.Sprintf(confluenceCloudPageCreateMessage, page.Title, page.Self, page.SpaceKey)
	case serializer.CommentCreatedEvent:
		message = fmt.Sprintf(confluenceCloudCommentCreateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case serializer.PageUpdatedEvent:
		message = fmt.Sprintf(confluenceCloudPageUpdateMessage, page.Title, page.Self, page.SpaceKey)
	case serializer.CommentUpdatedEvent:
		message = fmt.Sprintf(confluenceCloudCommentUpdateMessage, comment.Self, comment.Parent.Title, comment.Parent.Self)
	case serializer.PageRemovedEvent:
		message = fmt.Sprintf(confluenceCloudPageDeleteMessage, page.Title, page.SpaceKey)
	case serializer.CommentRemovedEvent:
		message = fmt.Sprintf(confluenceCloudCommentDeleteMessage, comment.Parent.Title, comment.Parent.Self)
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
