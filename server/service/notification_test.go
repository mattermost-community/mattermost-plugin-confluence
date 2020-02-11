package service

import (
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/mock"

	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
)

func TestNotifications(t *testing.T) {
	for name, val := range map[string]struct {
		baseURL string
		spaceKey string
		event string
		urlSpaceKeyCombinationSubscriptions map[string][]string
		runAssertions func(t *testing.T, a  *plugintest.API)
	}{
		"create event": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.CommentCreatedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
				"testtesttest1234": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			runAssertions: func(t *testing.T, a  *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 1)
			},
		},
		"page event": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.PageRemovedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.PageRemovedEvent},
			},
			runAssertions: func(t *testing.T, a *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 1)
			},
		},
		"single notification": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.CommentCreatedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			runAssertions: func(t *testing.T, a  *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 1)
			},
		},
		"multiple notification": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.CommentUpdatedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
				"testtesttest1234": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
				"testtesttest1235": {serializer.PageRemovedEvent, serializer.PageCreatedEvent, serializer.CommentUpdatedEvent},
				"testtesttest1236": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			runAssertions: func(t *testing.T, a  *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 4)
			},
		},
		"no notification": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.PageRemovedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			runAssertions: func(t *testing.T, a  *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 0)
			},
		},
		"multiple subscription single notification": {
			baseURL: "https://test.com",
			spaceKey: "TEST",
			event: serializer.CommentCreatedEvent,
			urlSpaceKeyCombinationSubscriptions : map[string][]string{
				"testtesttesttest": {serializer.CommentCreatedEvent, serializer.CommentUpdatedEvent},
				"testtesttest1234": {serializer.CommentRemovedEvent, serializer.CommentUpdatedEvent},
			},
			runAssertions: func(t *testing.T, a  *plugintest.API) {
				a.AssertNumberOfCalls(t, "CreatePost", 1)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			mockAPI := baseMock()
			mockAPI.On("LogError",
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string"),
				mock.AnythingOfTypeArgument("string")).Return(nil)
			mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
			monkey.Patch(GetURLSpaceKeyCombinationSubscriptions, func(baseURL, spaceKey string) (map[string][]string, string, error) {
				return val.urlSpaceKeyCombinationSubscriptions, "testSub", nil
			})
			SendConfluenceNotifications(&model.Post{}, "https://test.com", "TEST", val.event)
			val.runAssertions(t, mockAPI)
		})
	}
}
