package service

import (
	"encoding/json"
	"testing"

	"bou.ke/monkey"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/plugin/plugintest"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	return mockAPI
}

func TestGetChannelSubscriptions(t *testing.T) {
	for name, val := range map[string]struct {
		Subscriptions map[string]serializer.Subscription
		RunAssertions func(t *testing.T, s map[string]serializer.Subscription)
	}{
		"single subscription": {
			Subscriptions: map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttesttest",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				},
			},
			RunAssertions: func(t *testing.T, subscriptions map[string]serializer.Subscription) {
				expected := map[string]serializer.Subscription{
					"test": {
						Alias:     "test",
						BaseURL:   "https://test.com",
						SpaceKey:  "TS",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				}
				assert.Equal(t, expected, subscriptions)
			},
		},
		"multiple subscription": {
			Subscriptions: map[string]serializer.Subscription{
				"test": {
					Alias:     "test",
					BaseURL:   "https://test.com",
					SpaceKey:  "TS",
					ChannelID: "testtesttesttest",
					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
				},
				"test1": {
					Alias:     "test1",
					BaseURL:   "https://test1.com",
					SpaceKey:  "TS1",
					ChannelID: "testtesttesttest",
					Events:    []string{serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
				},
			},
			RunAssertions: func(t *testing.T, sub map[string]serializer.Subscription) {
				expected := map[string]serializer.Subscription{
					"test": {
						Alias:     "test",
						BaseURL:   "https://test.com",
						SpaceKey:  "TS",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
					},
					"test1": {
						Alias:     "test1",
						BaseURL:   "https://test1.com",
						SpaceKey:  "TS1",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
					},
				}
				assert.Equal(t, expected, sub)
			},
		},
		"no subscription": {
			Subscriptions: map[string]serializer.Subscription{},
			RunAssertions: func(t *testing.T, sub map[string]serializer.Subscription) {
				expected := map[string]serializer.Subscription{}
				assert.Equal(t, expected, sub)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()
			subscriptionBytes, err := json.Marshal(val.Subscriptions)
			assert.Nil(t, err)
			monkey.Patch(store.GetChannelSubscriptionKey, func(channelID string) string {
				return "testSubscriptionKey"
			})
			monkey.Patch(util.GetKeyHash, func(channelID string) string {
				return "testKey"
			})
			mockAPI.On("KVGet", "testKey").Return(subscriptionBytes, nil)
			sub, key, err := GetChannelSubscriptions("testtesttesttest")
			assert.Nil(t, err)
			assert.NotNil(t, sub)
			assert.NotNil(t, key)
			val.RunAssertions(t, sub)
		})
	}
}

func TestGetURLSpaceKeyCombinationSubscriptions(t *testing.T) {
	for name, val := range map[string]struct {
		Subscriptions map[string][]string
		RunAssertions func(t *testing.T, s map[string][]string)
	}{
		"multiple subscription": {
			Subscriptions: map[string][]string{
				"testtesttesttest": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
				"testtesttest1234": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent, serializer.CommentCreatedEvent},
			},
			RunAssertions: func(t *testing.T, s map[string][]string) {
				expected := map[string][]string{
					"testtesttesttest": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
					"testtesttest1234": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent, serializer.CommentCreatedEvent},
				}
				assert.Equal(t, expected, s)
			},
		},
		"single subscription": {
			Subscriptions: map[string][]string{
				"testtesttesttest": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
			},
			RunAssertions: func(t *testing.T, s map[string][]string) {
				expected := map[string][]string{
					"testtesttesttest": {serializer.CommentUpdatedEvent, serializer.PageRemovedEvent},
				}
				assert.Equal(t, expected, s)
			},
		},
		"no subscription": {
			Subscriptions: map[string][]string{},
			RunAssertions: func(t *testing.T, s map[string][]string) {
				expected := map[string][]string{}
				assert.Equal(t, expected, s)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()
			subscriptionBytes, err := json.Marshal(val.Subscriptions)
			assert.Nil(t, err)
			monkey.Patch(store.GetURLSpaceKeyCombinationKey, func(url, spaceKey string) (string, error) {
				return "testSubscriptionKey", nil
			})
			monkey.Patch(util.GetKeyHash, func(channelID string) string {
				return "testKey"
			})
			mockAPI.On("KVGet", "testKey").Return(subscriptionBytes, nil)
			sub, key, err := GetURLSpaceKeyCombinationSubscriptions("https://test.com", "TS")
			assert.Nil(t, err)
			assert.NotNil(t, sub)
			assert.NotNil(t, key)
			val.RunAssertions(t, sub)
		})
	}
}
