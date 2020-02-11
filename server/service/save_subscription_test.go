package service

import (
	"net/http"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

func TestSaveSubscription(t *testing.T) {
	for name, val := range map[string]struct {
		subscription serializer.Subscription
		statusCode   int
		errorMessage string
		channelSubscriptions map[string]serializer.Subscription
		urlSpaceKeyCombinationSubscriptions map[string][]string
	}{
		"alias already exist": {
			subscription: serializer.Subscription{
				Alias:     "test",
				BaseURL:   "https://test.com",
				SpaceKey:  "TS",
				ChannelID: "testtesttesttest",
				Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: aliasAlreadyExist,
			channelSubscriptions : map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttesttest",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				},
			},
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"url space key combination already exist": {
			subscription: serializer.Subscription{
				Alias:     "test1",
				BaseURL:   "https://test.com",
				SpaceKey:  "TS",
				ChannelID: "testtesttesttest",
				Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: urlSpaceKeyAlreadyExist,
			channelSubscriptions : map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttesttest",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				},
			},
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"subscription unique base url": {
			subscription: serializer.Subscription{
				Alias:     "test1",
				BaseURL:   "https://test1.com",
				SpaceKey:  "TS",
				ChannelID: "testtesttest1234",
				Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			statusCode:   http.StatusOK,
			errorMessage: "",
			channelSubscriptions : map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttest1234",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				},
			},
			urlSpaceKeyCombinationSubscriptions : map[string][]string{},
		},
		"subscription unique space key": {
			subscription: serializer.Subscription{
				Alias:     "test1",
				BaseURL:   "https://test.com",
				SpaceKey:  "TSST",
				ChannelID: "testtesttest1234",
				Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			statusCode:   http.StatusOK,
			errorMessage: "",
			channelSubscriptions : map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttest1234",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				},
			},
			urlSpaceKeyCombinationSubscriptions : map[string][]string{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()
			monkey.Patch(GetChannelSubscriptions, func(channelID string) (map[string]serializer.Subscription, string, error) {
				return val.channelSubscriptions, "testSub", nil
			})
			monkey.Patch(GetURLSpaceKeyCombinationSubscriptions, func(baseURL, spaceKey string) (map[string][]string, string, error) {
				return val.urlSpaceKeyCombinationSubscriptions, "testSub", nil
			})
			mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			errCode, err := SaveNewSubscription(val.subscription)
			assert.Equal(t, val.statusCode, errCode)
			if err != nil {
				assert.Equal(t, val.errorMessage, err.Error())
			}
		})
	}
}
