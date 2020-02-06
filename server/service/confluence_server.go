package service

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const (
	confluenceServerPageCreatedMessage  = "%s published a new page in %s."
	confluenceServerPageUpdatedMessage  = "%s updated %s in %s."
	confluenceServerPageTrashedMessage  = "%s moved page %s to trash in %s."
	confluenceServerPageRestoredMessage = "%s restored a page %s from trash in %s."
	confluenceServerPageRemovedMessage  = "%s purged a page **%s** from trash in %s."

	confluenceServerCommentCreatedMessage      = "%s commented on %s in %s."
	confluenceServerCommentReplyCreatedMessage = "%s replied to a comment on %s in %s."
	confluenceServerCommentUpdatedMessage      = "%s updated a comment on %s in %s."
	confluenceServerCommentRemovedMessage      = "%s removed a comment on %s in %s."
)

func SendConfluenceServerNotifications(event *serializer.ConfluenceServerEvent) {
	post := generateConfluenceServerNotificationPost(event)
	if post == nil {
		return
	}
	SendConfluenceNotifications(post, event.BaseURL, event.Space.Key, event.Event)
}

func generateConfluenceServerNotificationPost(event *serializer.ConfluenceServerEvent) *model.Post {
	var attachment *model.SlackAttachment
	post := &model.Post{
		UserId: config.BotUserID,
	}

	switch event.Event {
	case serializer.PageCreatedEvent:
		message := fmt.Sprintf(confluenceServerPageCreatedMessage, event.GetUserDisplayName(true), event.GetSpaceDisplayName(true))
		attachment = &model.SlackAttachment{
			Fallback:  message,
			Pretext:   message,
			Title:     event.Page.Title,
			TitleLink: event.Page.URL,
			Text:      event.Page.Excerpt,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Link",
					Value: fmt.Sprintf("[View in Confluence](%s)", event.Page.URL),
					Short: false,
				},
			},
		}

	case serializer.PageUpdatedEvent:
		message := fmt.Sprintf(confluenceServerPageUpdatedMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(true), event.GetSpaceDisplayName(true))
		if strings.TrimSpace(event.VersionComment) != "" {
			attachment = &model.SlackAttachment{
				Fallback: message,
				Pretext:  message,
				Text:     event.VersionComment,
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Link",
						Value: fmt.Sprintf("[View in Confluence](%s)", event.Page.URL),
						Short: false,
					},
				},
			}
		} else {
			post.Message = message
		}

	case serializer.PageTrashedEvent:
		post.Message = fmt.Sprintf(confluenceServerPageTrashedMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(true), event.GetSpaceDisplayName(true))

	case serializer.PageRestoredEvent:
		post.Message = fmt.Sprintf(confluenceServerPageRestoredMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(true), event.GetSpaceDisplayName(true))

	case serializer.PageRemovedEvent:
		// No link for page since the page was removed
		post.Message = fmt.Sprintf(confluenceServerPageRemovedMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(false), event.GetSpaceDisplayName(true))

	case serializer.CommentCreatedEvent:
		commentedOn := event.GetCommentPageOrBlogDisplayName()
		message := fmt.Sprintf(confluenceServerCommentCreatedMessage, event.GetUserDisplayName(true), commentedOn, event.GetSpaceDisplayName(true))
		var fields []*model.SlackAttachmentField

		if event.Comment.ParentComment != nil {
			message = fmt.Sprintf(confluenceServerCommentReplyCreatedMessage, event.GetUserDisplayName(true), commentedOn, event.GetSpaceDisplayName(true))
			fields = append(fields, &model.SlackAttachmentField{
				Title: "In Reply to",
				Value: event.Comment.ParentComment.Excerpt,
				Short: false,
			})
		}
		fields = append(fields, &model.SlackAttachmentField{
			Title: "Link",
			Value: fmt.Sprintf("[View in Confluence](%s)", event.Comment.URL),
			Short: false,
		})
		attachment = &model.SlackAttachment{
			Fallback: message,
			Pretext:  message,
			Text:     event.Comment.Excerpt,
			Fields:   fields,
		}

	case serializer.CommentUpdatedEvent:
		commentedOn := event.GetCommentPageOrBlogDisplayName()
		message := fmt.Sprintf(confluenceServerCommentUpdatedMessage, event.GetUserDisplayName(true), commentedOn, event.GetSpaceDisplayName(true))
		attachment = &model.SlackAttachment{
			Fallback: message,
			Pretext:  message,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Updated Comment",
					Value: event.Comment.Excerpt,
					Short: false,
				},
				{
					Title: "Link",
					Value: fmt.Sprintf("[View in Confluence](%s)", event.Comment.URL),
					Short: false,
				},
			},
		}

	case serializer.CommentRemovedEvent:
		commentedOn := event.GetCommentPageOrBlogDisplayName()
		message := fmt.Sprintf(confluenceServerCommentRemovedMessage, event.GetUserDisplayName(true), commentedOn, event.GetSpaceDisplayName(true))
		attachment = &model.SlackAttachment{
			Fallback: message,
			Pretext:  message,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Comment",
					Value: event.Comment.Excerpt,
					Short: false,
				},
			},
		}

	default:
		return nil
	}

	if attachment != nil {
		model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	}

	return post
}
