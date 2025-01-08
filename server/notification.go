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

func getNotification(plugin *Plugin) *notification {
	return &notification{
		plugin,
	}
}

func (n *notification) SendConfluenceNotifications(event serializer.ConfluenceEventV2, eventType, botUserID, userID string) {
	url := event.GetURL()
	var spaceKey string
	var pageID string

	if strings.Contains(eventType, Comment) {
		spaceKey = event.(*ConfluenceServerEvent).GetCommentSpaceKey()
		pageID = event.(*ConfluenceServerEvent).GetCommentContainerID()
	}
	if strings.Contains(eventType, Page) {
		spaceKey = event.(*ConfluenceServerEvent).GetPageSpaceKey()
		pageID = event.GetPageID()
	}
	if strings.Contains(eventType, Space) {
		spaceKey = event.GetSpaceKey()
		if spaceKey != "" {
			pageID = defaultPageID
		} else {
			return
		}
	}

	post := event.GetNotificationPost(eventType, url, botUserID)
	if post == nil || pageID == "" || url == "" || spaceKey == "" {
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
