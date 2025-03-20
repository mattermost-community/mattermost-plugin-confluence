package main

import (
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
)

const (
	helpText = "###### Mattermost Confluence Plugin - Slash Command Help\n\n" +
		"* `/confluence connect` - Connect your Mattermost user to Confluence.\n" +
		"* `/confluence disconnect` - Disconnect your Mattermost user from Confluence.\n" +
		"* `/confluence subscribe` - Subscribe the current channel to notifications from Confluence.\n" +
		"* `/confluence unsubscribe \"<name>\"` - Unsubscribe the current channel from notifications associated with the given subscription name.\n" +
		"* `/confluence list` - List all subscriptions for the current channel.\n" +
		"* `/confluence edit \"<name>\"` - Edit the subscription settings associated with the given subscription name.\n"
)

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	return mockAPI
}

func TestExecuteCommand(t *testing.T) {
	p := Plugin{}

	// TODO: Add the testcases for unsubscribe and other commands
	for name, val := range map[string]struct {
		commandArgs      *model.CommandArgs
		ephemeralMessage string
		isAdmin          bool
		patchAPICalls    func()
	}{
		"empty command ": {
			commandArgs:      &model.CommandArgs{Command: "/confluence", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: helpText + sysAdminHelpText,
			isAdmin:          true,
		},
		"help command": {
			commandArgs:      &model.CommandArgs{Command: "/confluence help", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: helpText + sysAdminHelpText,
			isAdmin:          true,
		},
		"invalid command": {
			commandArgs:      &model.CommandArgs{Command: "/confluence xyz", UserId: "abcdabcdabcdabcd", ChannelId: "testtesttesttest"},
			ephemeralMessage: invalidCommand,
			isAdmin:          true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()

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
