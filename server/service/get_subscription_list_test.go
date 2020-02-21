package service

import (
	"testing"

	"bou.ke/monkey"

	"github.com/stretchr/testify/assert"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

func TestGetSubscriptionsByChannelID(t *testing.T) {
	for name, val := range map[string]struct {
		channelID string
		expected serializer.StringSubscription
	}{
		"single subscription": {
			channelID: "testtesttesttest",
			expected: serializer.StringSubscription{
				"test":  serializer.SpaceSubscription{
					SpaceKey:  "TS",
					BaseSubscription: serializer.BaseSubscription{
						Alias:     "test",
						BaseURL:   "https://test.com",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
			},
		},
		"multiple subscription": {
			channelID: "testtesttesttes1",
			expected: serializer.StringSubscription{
				"test":  serializer.PageSubscription{
					PageID:  "1234",
					BaseSubscription: serializer.BaseSubscription{
						Alias:     "test",
						BaseURL:   "https://test.com",
						ChannelID: "testtesttesttes1",
						Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
				},
				"tes1":  serializer.SpaceSubscription{
					SpaceKey:  "TS",
					BaseSubscription: serializer.BaseSubscription{
						Alias:     "test",
						BaseURL:   "https://test.com",
						ChannelID: "testtesttesttes1",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
			},
		},
		"no subscription": {
			channelID: "testtsettest1234",
			expected: serializer.StringSubscription(nil),
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			subscriptions := serializer.Subscriptions{
				ByChannelID: map[string]serializer.StringSubscription {
					"testtesttesttest" : {
						"test":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
					"testtesttesttes1" : {
						"test":  serializer.PageSubscription{
							PageID:  "1234",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"tes1":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
				},
				ByURLSpaceKey: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/TS": {
						"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
						"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
				ByURLPagID: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/1234": {
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
				},
			}
			monkey.Patch(GetSubscriptions, func()(serializer.Subscriptions, error) {
				return subscriptions, nil
			})
			sub, err := GetSubscriptionsByChannelID(val.channelID)
			assert.Nil(t, err)
			assert.Equal(t, val.expected, sub)
		})
	}
}

func TestGetSubscriptionsByURLPageID(t *testing.T) {
	for name, val := range map[string]struct {
		url      string
		pageID   string
		expected serializer.StringArrayMap
	}{
		"single subscription": {
			url:    "https://test.com",
			pageID: "1234",
			expected: serializer.StringArrayMap{
				"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"multiple subscription": {
			url:    "https://test.com",
			pageID: "12345",
			expected: serializer.StringArrayMap{
				"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
				"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"no subscription": {
			url:      "https://test.com",
			pageID:   "123456",
			expected: serializer.StringArrayMap(nil),
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			subscriptions := serializer.Subscriptions{
				ByChannelID: map[string]serializer.StringSubscription {
					"testtesttesttest" : {
						"test":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"abc":  serializer.PageSubscription{
							PageID:  "12345",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "abc",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
					"testtesttesttes1" : {
						"test":  serializer.PageSubscription{
							PageID:  "1234",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"xyz":  serializer.PageSubscription{
							PageID:  "12345",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "xyz",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"tes1":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "tes1",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
				},
				ByURLSpaceKey: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/TS": {
						"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
						"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
				ByURLPagID: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/1234": {
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
					"confluence_subs/test.com/12345": {
						"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
				},
			}
			monkey.Patch(GetSubscriptions, func()(serializer.Subscriptions, error) {
				return subscriptions, nil
			})
			sub, err := GetSubscriptionsByURLPageID(val.url, val.pageID)
			assert.Nil(t, err)
			assert.Equal(t, val.expected, sub)
		})
	}
}

func TestGetSubscriptionsByURLSpaceKey(t *testing.T) {
	for name, val := range map[string]struct {
		url      string
		spaceKey   string
		expected serializer.StringArrayMap
	}{
		"single subscription": {
			url:    "https://test.com",
			spaceKey: "TS1",
			expected: serializer.StringArrayMap{
				"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"multiple subscription": {
			url:    "https://test.com",
			spaceKey: "TS",
			expected: serializer.StringArrayMap{
				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
		},
		"no subscription": {
			url:      "https://test.com",
			spaceKey:   "ggh",
			expected: serializer.StringArrayMap(nil),
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			subscriptions := serializer.Subscriptions{
				ByChannelID: map[string]serializer.StringSubscription {
					"testtesttesttest" : {
						"test":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"abc":  serializer.PageSubscription{
							PageID:  "12345",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "abc",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
					"testtesttesttes1" : {
						"test":  serializer.PageSubscription{
							PageID:  "1234",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"xyz":  serializer.PageSubscription{
							PageID:  "12345",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "xyz",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"tes1":  serializer.SpaceSubscription{
							SpaceKey:  "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "tes1",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
						"tesa":  serializer.SpaceSubscription{
							SpaceKey:  "TS1",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "tesa",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttes1",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
				},
				ByURLSpaceKey: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/TS": {
						"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
						"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
					"confluence_subs/test.com/TS1": {
						"testtesttesttes1": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
				ByURLPagID: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/1234": {
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
					"confluence_subs/test.com/12345": {
						"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
				},
			}
			monkey.Patch(GetSubscriptions, func()(serializer.Subscriptions, error) {
				return subscriptions, nil
			})
			sub, err := GetSubscriptionsByURLSpaceKey(val.url, val.spaceKey)
			assert.Nil(t, err)
			assert.Equal(t, val.expected, sub)
		})
	}
}
