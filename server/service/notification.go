package service

import (
	"github.com/mattermost/mattermost-server/model"
)

func SendConfluenceNotifications(post *model.Post, baseURL, spaceKey, event string) {
	// subscriptions, _, err := GetURLSpaceKeyCombinationSubscriptions(baseURL, spaceKey)
	// if err != nil {
	// 	config.Mattermost.LogError("Unable to get subscribed channels.", "Error", err.Error())
	// 	return
	// }
	//
	// for channelID, events := range subscriptions {
	// 	if funk.Contains(events, event) {
	// 		post.ChannelId = channelID
	// 		if _, err := config.Mattermost.CreatePost(post); err != nil {
	// 			config.Mattermost.LogError("Unable to create Post in Mattermost", "Error", err.Error())
	// 		}
	// 	}
	// }
}
