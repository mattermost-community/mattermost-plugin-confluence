package main

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

const (
	specifyAlias              = "Please specify an alias."
	subscriptionDeleteSuccess = "**%s** has been deleted."
	noChannelSubscription     = "No subscriptions found for this channel."
	helpText                  = "###### Mattermost Confluence Plugin - Slash Command Help\n\n" +
		"* `/confluence subscribe` - Subscribe the current channel to notifications from Confluence.\n" +
		"* `/confluence unsubscribe \"<alias>\"` - Unsubscribe the current channel from notifications associated with the given alias.\n" +
		"* `/confluence list` - List all subscriptions for the current channel.\n" +
		"* `/confluence edit \"<alias>\"` - Edit the subscription settings associated with the given alias."
	invalidCommand			  = "Invalid command parameters. Please use `/confluence help` for more information."
)

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	return mockAPI
}

func TestExecuteCommand(t *testing.T) {
	p := Plugin{}
	mockAPI := baseMock()

	for name, val := range map[string]struct {
		commandArgs       *model.CommandArgs
		patchFunctionCalls func()
		ephemeralMessage string
	}{
		"empty command ": {
			commandArgs:       &model.CommandArgs{Command: "/confluence", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: helpText,
		},
		"list command": {
			commandArgs:       &model.CommandArgs{Command: "/confluence list", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			patchFunctionCalls: func() {
				channelSubscriptions := map[string]serializer.Subscription{
					"test": {
						Alias:     "test",
						BaseURL:   "https://test.com",
						SpaceKey:  "TS",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
					"test1": {
						Alias:     "test1",
						BaseURL:   "https://test1.com",
						SpaceKey:  "TS1",
						ChannelID: "testtesttesttest",
						Events:    []string{serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
					},
				}
				monkey.Patch(service.GetChannelSubscriptions, func(channelID string) (map[string]serializer.Subscription, string, error) {
					return channelSubscriptions, "testSub", nil
				})
			},
			ephemeralMessage: "| Alias | Base Url | Space Key | Events|\n| :----|:--------| :--------| :-----|\n|test|https://test.com|TS|Comment Remove, Comment Update|\n|test1|https://test1.com|TS1|Comment Remove, Comment Update|",
		},
		"list command empty list": {
			commandArgs:       &model.CommandArgs{Command: "/confluence list", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			patchFunctionCalls: func() {
				monkey.Patch(service.GetChannelSubscriptions, func(channelID string) (map[string]serializer.Subscription, string, error) {
					return map[string]serializer.Subscription{}, "testSub", nil
				})
			},
			ephemeralMessage: noChannelSubscription,
		},
		"help command": {
			commandArgs:       &model.CommandArgs{Command: "/confluence help", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: helpText,
		},
		"unsubscribe command ": {
			commandArgs:       &model.CommandArgs{Command: "/confluence unsubscribe \"abc\"", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			patchFunctionCalls: func () {
				monkey.Patch(service.DeleteSubscription, func (channelID, alias string) error {
					return nil
				})
			},
			ephemeralMessage: fmt.Sprintf(subscriptionDeleteSuccess, "abc"),
		},
		"unsubscribe command no alias": {
			commandArgs:       &model.CommandArgs{Command: "/confluence unsubscribe", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: specifyAlias,
		},
		"invalid command": {
			commandArgs:       &model.CommandArgs{Command: "/confluence xyz", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: invalidCommand,
		},

	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
				post := args.Get(1).(*model.Post)
				assert.Equal(t, val.ephemeralMessage, post.Message)
			}).Once().Return(&model.Post{})
			if val.patchFunctionCalls != nil {
				val.patchFunctionCalls()
			}
			res, err := p.ExecuteCommand(&plugin.Context{}, val.commandArgs)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		})
	}
}
