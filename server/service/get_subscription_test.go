package service

//
// import (
// 	"fmt"
// 	"net/http"
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
// func TestGetChannelSubscription(t *testing.T) {
// 	for name, val := range map[string]struct {
// 		channelID    string
// 		alias        string
// 		statusCode   int
// 		errorMessage string
// 	}{
// 		"get subscription success": {
// 			channelID:    "testtesttesttest",
// 			alias:        "test",
// 			statusCode:   http.StatusOK,
// 			errorMessage: "",
// 		},
// 		"subscription not found for alias": {
// 			channelID:    "testtesttesttest",
// 			alias:        "test4",
// 			statusCode:   http.StatusBadRequest,
// 			errorMessage: fmt.Sprintf(subscriptionNotFound, "test4"),
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
// 				"test1": {
// 					Alias:     "test1",
// 					BaseURL:   "https://test1.com",
// 					SpaceKey:  "TS1",
// 					ChannelID: "testtesttesttest",
// 					Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
// 				},
// 			}
// 			monkey.Patch(GetChannelSubscriptions, func(channelID string) (map[string]serializer.Subscription, string, error) {
// 				return channelSubscriptions, "testSub", nil
// 			})
// 			mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
// 			subscription, errCode, err := GetChannelSubscription(val.channelID, val.alias)
// 			assert.Equal(t, val.statusCode, errCode)
// 			if err != nil {
// 				assert.Equal(t, val.errorMessage, err.Error())
// 				return
// 			}
// 			assert.NotNil(t, subscription)
// 			assert.Equal(t, subscription.Alias, val.alias)
// 		})
// 	}
// }
