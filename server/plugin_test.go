package main

import (
	"fmt"
	"testing"

	"bou.ke/monkey"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-plugin-confluence/server/service"
)

const (
	userID    = "abcdabcdabcdabcd"
	channelID = "testChannelID"
)

func TestExecuteCommand(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	for name, val := range map[string]struct {
		commandArgs      *model.CommandArgs
		ephemeralMessage string
		isAdmin          bool
		patchAPICalls    func()
	}{
		"empty command ": {
			commandArgs:      &model.CommandArgs{Command: "/confluence", UserId: userID, ChannelId: channelID},
			ephemeralMessage: commonHelpText,
			isAdmin:          false,
		},
		"help command": {
			commandArgs:      &model.CommandArgs{Command: "/confluence help", UserId: userID, ChannelId: channelID},
			ephemeralMessage: commonHelpText,
			isAdmin:          false,
		},
		"unsubscribe command ": {
			commandArgs:      &model.CommandArgs{Command: "/confluence unsubscribe \"abc\"", UserId: userID, ChannelId: channelID},
			ephemeralMessage: fmt.Sprintf(subscriptionNotFound, "abc"),
			isAdmin:          false,
			patchAPICalls: func() {
				monkey.Patch(p.DeleteSubscription, func(channelID, alias, userID string) error {
					return nil
				})
			},
		},
		"unsubscribe command no alias": {
			commandArgs:      &model.CommandArgs{Command: "/confluence unsubscribe", UserId: userID, ChannelId: channelID},
			ephemeralMessage: specifyAlias,
			isAdmin:          false,
		},
		"invalid command": {
			commandArgs:      &model.CommandArgs{Command: "/confluence xyz", UserId: userID, ChannelId: channelID},
			ephemeralMessage: invalidCommand,
			isAdmin:          false,
		},
		"admin restricted": {
			commandArgs:      &model.CommandArgs{Command: "/confluence install \"server\" \"https:\\confluence.test.com\"", UserId: userID, ChannelId: channelID},
			ephemeralMessage: installOnlySystemAdmin,
			isAdmin:          false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()

			monkey.Patch(service.GetSubscriptions, func() (serializer.Subscriptions, error) {
				return subscriptions, nil
			})

			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
				post := args.Get(1).(*model.Post)
				assert.Equal(t, val.ephemeralMessage, post.Message)
			}).Once().Return(&model.Post{})

			roles := "system_user"
			if val.isAdmin {
				roles += " system_admin"
			}
			mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Id: "123", Roles: roles}, nil)
			if val.patchAPICalls != nil {
				val.patchAPICalls()
			}

			res, err := p.ExecuteCommand(&plugin.Context{}, val.commandArgs)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		})
	}
}
