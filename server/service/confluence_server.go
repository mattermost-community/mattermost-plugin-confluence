package service

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

const (
	confluenceServerPageCreatedMessage            = "%s published a new page in %s."
	confluenceServerPageCreatedWithoutBodyMessage = "%s published a new page %s in %s."
	confluenceServerPageUpdatedMessage            = "%s updated %s in %s."
	confluenceServerPageTrashedMessage            = "%s trashed %s in %s."
	confluenceServerPageRestoredMessage           = "%s restored %s in %s."
	confluenceServerPageRemovedMessage            = "%s removed **%s** in %s."

	confluenceServerCommentCreatedMessage      = "%s commented on %s in %s."
	confluenceServerEmptyCommentCreatedMessage = "%s [commented](%s) on %s in %s."
	confluenceServerCommentReplyCreatedMessage = "%s replied to a comment on %s in %s."
	confluenceServerCommentUpdatedMessage      = "%s updated a comment on %s in %s."
	confluenceServerEmptyCommentUpdatedMessage = "%s updated a [comment](%s) on %s in %s."
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
		if strings.TrimSpace(event.Page.Excerpt) != "" {
			attachment = &model.SlackAttachment{
				Fallback:  message,
				Pretext:   message,
				Title:     event.Page.Title,
				TitleLink: event.Page.TinyURL,
				Text:      fmt.Sprintf("%s\n\n[**View in Confluence**](%s)", strings.TrimSpace(event.Page.Excerpt), event.Page.TinyURL),
			}
		} else {
			post.Message = fmt.Sprintf(confluenceServerPageCreatedWithoutBodyMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(true), event.GetSpaceDisplayName(true))
		}

	case serializer.PageUpdatedEvent:
		message := fmt.Sprintf(confluenceServerPageUpdatedMessage, event.GetUserDisplayName(true), event.GetPageDisplayName(true), event.GetSpaceDisplayName(true))
		if strings.TrimSpace(event.VersionComment) != "" {
			attachment = &model.SlackAttachment{
				Fallback: message,
				Pretext:  message,
				Text:     fmt.Sprintf("**What’s Changed?**\n> %s\n\n[**View in Confluence**](%s)", strings.TrimSpace(event.VersionComment), event.Page.TinyURL),
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
		message := fmt.Sprintf(confluenceServerCommentCreatedMessage, event.GetUserDisplayName(true), event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))

		text := ""
		if strings.TrimSpace(event.Comment.Excerpt) != "" {
			text += fmt.Sprintf("**%s wrote:**\n> %s\n\n", event.GetUserFirstName(), strings.TrimSpace(event.Comment.Excerpt))
		}
		if event.Comment.ParentComment != nil && strings.TrimSpace(event.Comment.ParentComment.Excerpt) != "" {
			message = fmt.Sprintf(confluenceServerCommentReplyCreatedMessage, event.GetUserDisplayName(true), event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))
			text += fmt.Sprintf("**In Reply to:**\n> %s\n", strings.TrimSpace(event.Comment.ParentComment.Excerpt))
		}

		if text != "" {
			attachment = &model.SlackAttachment{
				Fallback: message,
				Pretext:  message,
				Text:     fmt.Sprintf("%s\n\n[**View in Confluence**](%s)", text, event.Comment.URL),
			}
		} else {
			post.Message = fmt.Sprintf(confluenceServerEmptyCommentCreatedMessage, event.GetUserDisplayName(true), event.Comment.URL, event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))
		}

	case serializer.CommentUpdatedEvent:
		message := fmt.Sprintf(confluenceServerCommentUpdatedMessage, event.GetUserDisplayName(true), event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))
		if strings.TrimSpace(event.Comment.Excerpt) != "" {
			attachment = &model.SlackAttachment{
				Fallback: message,
				Pretext:  message,
				Text:     fmt.Sprintf("**Updated Comment:**\n> %s\n\n[**View in Confluence**](%s)", strings.TrimSpace(event.Comment.Excerpt), event.Comment.URL),
			}
		} else {
			post.Message = fmt.Sprintf(confluenceServerEmptyCommentUpdatedMessage, event.GetUserDisplayName(true), event.Comment.URL, event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))
		}

	case serializer.CommentRemovedEvent:
		// No link since the comment was removed.
		message := fmt.Sprintf(confluenceServerCommentRemovedMessage, event.GetUserDisplayName(true), event.GetCommentPageOrBlogDisplayName(true), event.GetSpaceDisplayName(true))
		if strings.TrimSpace(event.Comment.Excerpt) != "" {
			attachment = &model.SlackAttachment{
				Fallback: message,
				Pretext:  message,
				Text:     fmt.Sprintf("**Deleted Comment:**\n> %s", strings.TrimSpace(event.Comment.Excerpt)),
			}
		} else {
			post.Message = message
		}

	default:
		return nil
	}

	if attachment != nil {
		model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	}

	return post
}
