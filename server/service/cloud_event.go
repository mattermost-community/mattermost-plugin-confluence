package service

import (
	"fmt"
	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-server/model"
)

func SendCloudNotification (cloudEvent serializer.ConfluenceCloudEvent, event string) *model.AppError {
	message := ""
	page := cloudEvent.Page
	comment := cloudEvent.Comment
	switch event {
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

	// TODO : Update the send notification logic.
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   message,
	}
	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}
