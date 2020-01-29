package store

import (
	"encoding/json"
	"fmt"
	url2 "net/url"

	"github.com/Brightscout/mattermost-plugin-confluence/server/util"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/pkg/errors"
)

const ConfluenceSubscriptionKeyPrefix = "confluence_subscription"

func Set(key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if appErr := config.Mattermost.KVSet(util.GetKeyHash(key), bytes); appErr != nil {
		return errors.New(appErr.Error())
	}
	return nil
}

func Get(key string, data interface{}) error {
	bytes, appErr := config.Mattermost.KVGet(util.GetKeyHash(key))
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	if bytes == nil {
		return nil
	}

	err := json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	return nil
}

func GetURLSpaceKeyCombinationKey(url, spaceKey string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", ConfluenceSubscriptionKeyPrefix, u.Hostname(), spaceKey), nil
}

func GetChannelSubscriptionKey(channelID string) string {
	return fmt.Sprintf("%s/%s", ConfluenceSubscriptionKeyPrefix, channelID)
}
