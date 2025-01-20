package main

import (
	"strings"

	"github.com/thoas/go-funk"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
)

const defaultPageID = "-1"

type notification struct {
	*Plugin
}

func (p *Plugin) getNotification() *notification {
	return &notification{
		p,
	}
}

func (n *notification) SendConfluenceNotifications(event serializer.ConfluenceEventV2, eventType, botUserID, userID string) {
	url := event.GetURL()
	if url == "" {
		return
	}

	spaceKey, pageID := n.extractSpaceKeyAndPageID(event, eventType)
	if spaceKey == "" || pageID == "" {
		return
	}

	post := event.GetNotificationPost(eventType, url, botUserID)
	if post == nil {
		return
	}

	subscriptionChannelIDs := n.getNotificationChannelIDs(url, spaceKey, pageID, eventType, userID)
	for _, channelID := range subscriptionChannelIDs {
		post.ChannelId = channelID
		if _, err := n.API.CreatePost(post); err != nil {
			n.API.LogError("Unable to create Post in Mattermost", "Error", err.Error())
		}
	}
}

func (n *notification) extractSpaceKeyAndPageID(event serializer.ConfluenceEventV2, eventType string) (string, string) {
	var spaceKey, pageID string

	switch {
	case strings.Contains(eventType, Comment):
		if e, ok := event.(*ConfluenceServerEvent); ok {
			spaceKey = e.GetCommentSpaceKey()
			pageID = e.GetCommentContainerID()
		}
	case strings.Contains(eventType, Page):
		if e, ok := event.(*ConfluenceServerEvent); ok {
			spaceKey = e.GetPageSpaceKey()
			pageID = event.GetPageID()
		}
	case strings.Contains(eventType, Space):
		spaceKey = event.GetSpaceKey()
		if spaceKey != "" {
			pageID = defaultPageID
		}
	}

	return spaceKey, pageID
}

func (n *notification) getNotificationChannelIDs(url, spaceKey, pageID, eventType, userID string) []string {
	urlSpaceKeySubscriptions, err := service.GetSubscriptionsByURLSpaceKey(url, spaceKey)
	if err != nil {
		n.API.LogError("Unable to get subscribed channels for spaceKey.", "Error", err.Error())
		return nil
	}
	urlPageIDSubscriptions, err := service.GetSubscriptionsByURLPageID(url, pageID)
	if err != nil {
		n.API.LogError("Unable to get subscribed channels for pageID.", "Error", err.Error())
		return nil
	}

	urlPageIDSubscriptionChannelIDs := GetURLSubscriptionChannelIDs(urlPageIDSubscriptions, eventType, userID)
	urlSpaceKeySubscriptionChannelIDs := GetURLSubscriptionChannelIDs(urlSpaceKeySubscriptions, eventType, userID)

	return util.Deduplicate(append(urlSpaceKeySubscriptionChannelIDs, urlPageIDSubscriptionChannelIDs...))
}

func GetURLSubscriptionChannelIDs(urlSubscriptions serializer.StringArrayMap, eventType, userID string) []string {
	var urlSubscriptionChannelIDs []string

	for channelID, events := range urlSubscriptions {
		if funk.Contains(events, eventType) {
			urlSubscriptionChannelIDs = append(urlSubscriptionChannelIDs, channelID)
		}
	}

	return urlSubscriptionChannelIDs
}

func GetURLSubscriptionUserIDs(urlSubscriptions serializer.StringArrayMap, eventType string) []string {
	var userIDs []string

	for id, events := range urlSubscriptions {
		if funk.Contains(events, eventType) {
			userIDs = append(userIDs, id)
		}
	}

	return userIDs
}
