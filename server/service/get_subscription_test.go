package service

import (
	"fmt"
	"net/http"
	"testing"

	"bou.ke/monkey"

	"github.com/stretchr/testify/assert"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

func TestGetChannelSubscription(t *testing.T) {
	for name, val := range map[string]struct {
		channelID    string
		alias        string
		statusCode   int
		errorMessage string
	}{
		"get subscription success": {
			channelID:    "testtesttesttest",
			alias:        "test",
			statusCode:   http.StatusOK,
			errorMessage: "",
		},
		"subscription not found for alias": {
			channelID:    "testtesttesttest",
			alias:        "test4",
			statusCode:   http.StatusBadRequest,
			errorMessage: fmt.Sprintf(subscriptionNotFound, "test4"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			subscriptions := serializer.Subscriptions{
				ByChannelID: map[string]serializer.StringSubscription{
					"testtesttesttest": {
						"test": serializer.SpaceSubscription{
							SpaceKey: "TS",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
					"testtesttesttes1": {
						"test": serializer.PageSubscription{
							PageID: "1234",
							BaseSubscription: serializer.BaseSubscription{
								Alias:     "test",
								BaseURL:   "https://test.com",
								ChannelID: "testtesttesttest",
								Events:    []string{serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
							},
						},
					},
				},
				ByURLSpaceKey: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/TS": {
						"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				},
				ByURLPagID: map[string]serializer.StringArrayMap{
					"confluence_subs/test.com/1234": {
						"testtesttesttes1": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
					},
				},
			}
			monkey.Patch(GetSubscriptions, func() (serializer.Subscriptions, error) {
				return subscriptions, nil
			})
			subscription, errCode, err := GetChannelSubscription(val.channelID, val.alias)
			assert.Equal(t, val.statusCode, errCode)
			if err != nil {
				assert.Equal(t, val.errorMessage, err.Error())
				return
			}
			assert.NotNil(t, subscription)
			assert.Equal(t, subscription.(serializer.SpaceSubscription).Alias, val.alias)
		})
	}
}
