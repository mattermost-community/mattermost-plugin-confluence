package util

import (
	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-server/model"
)

func PostCommandResponse(context *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: context.ChannelId,
		Message:   text,
	}
	_ = config.Mattermost.SendEphemeralPost(context.UserId, post)
}
