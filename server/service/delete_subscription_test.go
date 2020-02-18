package service

// import (
// 	"fmt"
// 	"testing"
//
// 	"bou.ke/monkey"
// 	"github.com/stretchr/testify/mock"
//
// 	"github.com/stretchr/testify/assert"
//
// 	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
// )
//
// func TestDeleteSubscription(t *testing.T) {
// 	for name, val := range map[string]struct {
// 		channelID string
// 		alias     string
// 		apiCalls  func(t *testing.T, channelID, alias string)
// 	}{
// 		"subscription delete success": {
// 			channelID: "testtestesttest",
// 			alias:     "test",
// 			apiCalls: func(t *testing.T, channelID, alias string) {
// 				err := DeleteSubscription(channelID, alias)
// 				assert.Nil(t, err)
// 			},
// 		},
// 		"subscription not found": {
// 			channelID: "testtestesttest",
// 			alias:     "test1",
// 			apiCalls: func(t *testing.T, channelID, alias string) {
// 				err := DeleteSubscription(channelID, alias)
// 				assert.NotNil(t, err)
// 				assert.Equal(t, fmt.Sprintf(subscriptionNotFound, alias), err.Error())
// 			},
// 		},
// 	} {
// 		t.Run(name, func(t *testing.T) {
// 			defer monkey.UnpatchAll()
// 			mockAPI := baseMock()
// 			channelSubscriptions := map[string]serializer.Subscription{
// 				"test": {
// 					Alias:     "test",
// 					BaseURL:   "https://test.com",
// 					SpaceKey:  "TS",
// 					ChannelID: "testtesttesttest",
// 					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
// 				},
// 			}
// 			urlSpaceKeyCombinationSubscriptions := map[string][]string{
// 				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
// 			}
// 			monkey.Patch(GetChannelSubscriptions, func(channelID string) (map[string]serializer.Subscription, string, error) {
// 				return channelSubscriptions, "testSub", nil
// 			})
// 			monkey.Patch(GetURLSpaceKeyCombinationSubscriptions, func(baseURL, spaceKey string) (map[string][]string, string, error) {
// 				return urlSpaceKeyCombinationSubscriptions, "testSub", nil
// 			})
// 			mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
// 			val.apiCalls(t, val.channelID, val.alias)
// 		})
// 	}
// }
