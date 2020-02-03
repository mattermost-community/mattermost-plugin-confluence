package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/model"
)

var CloudEvent = &Endpoint{
	RequiresAuth: false,
	Path:         "/cloud/{event:[A-Za-z0-9_]+}",
	Method:       http.MethodPost,
	Execute:      parseCloudEvent,
}

func parseCloudEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	body := json.NewDecoder(r.Body)
	cloudEvent := serializer.ConfluenceCloudEvent{}
	if err := body.Decode(&cloudEvent); err != nil {
		config.Mattermost.LogError("Error decoding request body.", "Error", err.Error())
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}

	switch params["event"] {
	case "page_created":
		pageCreateEvent(cloudEvent.Page)
	case "comment_created":
		commentCreateEvent(cloudEvent.Comment)
	case "page_update":
		pageUpdateEvent(cloudEvent.Page)
	case "comment_update":
		commentUpdate(cloudEvent.Comment)
	case "page_delete":
		pageDeleteEvent(cloudEvent.Page)
	case "comment_delete":
		commentDeleteEvent(cloudEvent.Comment)
	}
}

func pageCreateEvent(page *serializer.Page) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Page](%s) created in space **%s**", page.Self, page.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", page.Title, page.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, page.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}

func commentCreateEvent(comment *serializer.Comment) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Comment](%s) created in space **%s**", comment.Self, comment.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", comment.Parent.Title, comment.Parent.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, comment.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}

func pageDeleteEvent(page *serializer.Page) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Page](%s) deleted in space **%s**", page.Self, page.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", page.Title, page.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, page.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}

func pageUpdateEvent(page *serializer.Page) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Page](%s) updated in space **%s**", page.Self, page.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", page.Title, page.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, page.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}

func commentUpdate(comment *serializer.Comment) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Comment](%s) updated in space **%s**", comment.Self, comment.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", comment.Parent.Title, comment.Parent.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, comment.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}

func commentDeleteEvent(comment *serializer.Comment) *model.AppError {
	post := &model.Post{
		ChannelId: "617td783y38y8egdymy4w6qisw",
		UserId:    config.BotUserID,
		Type:      model.POST_DEFAULT,
		Message:   fmt.Sprintf("[Comment](%s) deleted in space **%s**", comment.Self, comment.SpaceKey),
	}

	slackAttachmentFields := []*model.SlackAttachmentField{
		{
			Title: "Page",
			Value: fmt.Sprintf("[%s](%s)", comment.Parent.Title, comment.Parent.Self),
			Short: true,
		},
		{
			Title: "Time",
			Value: time.Unix(0, comment.ModificationDate*int64(1000*1000)).Format(time.RFC1123),
			Short: true,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{
		{
			Fields: slackAttachmentFields,
		},
	})

	_, appErr := config.Mattermost.CreatePost(post)
	return appErr
}
