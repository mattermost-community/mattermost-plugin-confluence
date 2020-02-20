package service

import (
	"github.com/thoas/go-funk"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

func SendConfluenceNotifications(event serializer.ConfluenceEvent, eventType string) {
	details := event.GetEventDetails()
	post := event.GetNotificationPost(eventType)

	if post == nil || details.PageID == "" || details.URL == "" || details.SpaceKey == "" {
		return
	}
	subscriptionChannelIDs := getNotificationChannelIDs(details, eventType)
	for _, channelID := range subscriptionChannelIDs {
		post.ChannelId = channelID
		if _, err := config.Mattermost.CreatePost(post); err != nil {
			config.Mattermost.LogError("Unable to create Post in Mattermost", "Error", err.Error())
		}
	}
}

func getNotificationChannelIDs(d *serializer.Event, eventType string) []string {
	urlSpaceKeySubscriptions, err := GetSubscriptionsByURLSpaceKey(d.URL, d.SpaceKey)
	if err != nil {
		config.Mattermost.LogError("Unable to get subscribed channels.", "Error", err.Error())
		return nil
	}
	urlPageIDSubscriptions, err := GetSubscriptionsByURLPageID(d.URL, d.PageID)
	if err != nil {
		config.Mattermost.LogError("Unable to get subscribed channels.", "Error", err.Error())
		return nil
	}

	urlSpaceKeySubscriptionChannelIDs := make([]string, 0)
	urlPageIDSubscriptionChannelIDs := make([]string, 0)
	for channelID, events := range urlPageIDSubscriptions {
		if funk.Contains(events, eventType) {
			urlPageIDSubscriptionChannelIDs = append(urlPageIDSubscriptionChannelIDs, channelID)
		}
	}
	for channelID, events := range urlSpaceKeySubscriptions {
		if funk.Contains(events, eventType) {
			urlSpaceKeySubscriptionChannelIDs = append(urlSpaceKeySubscriptionChannelIDs, channelID)
		}
	}

	channelIDs := append(urlSpaceKeySubscriptionChannelIDs, urlPageIDSubscriptionChannelIDs...)
	return util.Deduplicate(channelIDs)
}
