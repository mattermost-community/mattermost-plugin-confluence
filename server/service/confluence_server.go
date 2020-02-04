package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

func SendConfluenceServerNotifications(event *serializer.ConfluenceServerEvent) {
	post := generateConfluenceServerNotificationPost(event)
	if post == nil {
		return
	}
	SendConfluenceNotifications(post, event.BaseURL, event.Space.Key, event.Event)
}

func generateConfluenceServerNotificationPost(event *serializer.ConfluenceServerEvent) *model.Post {
	eventMap := map[string]string{
		"comment_removed": "[%s](%s) deleted a [comment](%s) in [Confluence](%s).",
		"comment_created": "[%s](%s) created  a new [comment](%s) in [Confluence](%s).",
		"comment_updated": "[%s](%s) updated a [comment](%s) in [Confluence](%s).",
		"page_created":    "[%s](%s) created a new [page](%s) in [Confluence](%s).",
		"page_removed":    "[%s](%s) removed a [page](%s) from trash in [Confluence](%s).",
		"page_restored":   "[%s](%s) restored a [page](%s) from trash in [Confluence](%s).",
		"page_trashed":    "[%s](%s) moved a [page](%s) to trash in [Confluence](%s).",
		"page_updated":    "[%s](%s) updated a [page](%s) in [Confluence](%s).",
	}

	// if event is not in the eventMap
	if _, ok := eventMap[event.Event]; !ok {
		return nil
	}

	pretext := fmt.Sprintf(eventMap[event.Event], event.User.FullName, event.User.URL, event.ContentURL, event.BaseURL)
	fallback := fmt.Sprintf(eventMap[event.Event], event.User.FullName, event.User.URL, event.ContentURL, event.BaseURL)

	pageTitle := ""
	pageLink := ""
	pageExcerpt := ""

	var fields []*model.SlackAttachmentField

	if event.Comment != nil {
		fields = append(fields, &model.SlackAttachmentField{
			Title: "Comment",
			Value: event.Comment.Excerpt,
			Short: false,
		})

		if event.Comment.ParentComment != nil {
			fields = append(fields, &model.SlackAttachmentField{
				Title: "In Reply To",
				Value: event.Comment.ParentComment.Excerpt,
				Short: false,
			})
		}
	}

	if event.Page != nil {
		pageTitle = event.Page.Title
		pageLink = event.Page.URL
		pageExcerpt = event.Page.Excerpt
	} else if event.Blog != nil {
		pageTitle = event.Blog.Title
		pageLink = event.Blog.URL
		pageExcerpt = event.Blog.Excerpt
	}

	// If event is an instance of Edited event
	if strings.TrimSpace(event.VersionComment) != "" {
		fields = append(fields, &model.SlackAttachmentField{
			Title: "Version Comment",
			Value: strings.TrimSpace(event.VersionComment),
			Short: true,
		})
	}

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Space",
		Value: fmt.Sprintf("[%s](%s)", event.Space.Name, event.Space.URL),
		Short: true,
	})

	fields = append(fields, &model.SlackAttachmentField{
		Title: "Time",
		Value: time.Unix(0, event.Timestamp*int64(1000*1000)).Format(time.RFC1123), // TODO: Use a util to parse time in a human readable format.
		Short: true,
	})

	attachment := &model.SlackAttachment{
		Fallback:  fallback,
		Pretext:   pretext,
		Title:     pageTitle,
		TitleLink: pageLink,
		Text:      pageExcerpt,
		Fields:    fields,
	}

	post := &model.Post{
		UserId: config.BotUserID,
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	return post
}
