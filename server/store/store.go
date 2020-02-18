package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	url2 "net/url"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

const ConfluenceSubscriptionKeyPrefix = "confluence_subscription121"

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

func GetURLSpaceKeyCombinationKey(url, spaceKey string) string {
	u, _ := url2.Parse(url)
	return fmt.Sprintf("%s/%s/%s", ConfluenceSubscriptionKeyPrefix, u.Hostname(), spaceKey)
}

func GetURLPageIDCombinationKey(url, pageID string) string {
	u, _ := url2.Parse(url)
	return fmt.Sprintf("%s/%s/%s", ConfluenceSubscriptionKeyPrefix, u.Hostname(), pageID)
}

func GetChannelSubscriptionKey(channelID string) string {
	return fmt.Sprintf("%s/%s", ConfluenceSubscriptionKeyPrefix, channelID)
}

func GetSubscriptionKey() string {
	return "abcde"
}

func AtomicModify(key string, modify func(initialValue []byte) ([]byte, error)) error {
	readModify := func() ([]byte, []byte, error) {
		initialBytes, appErr := config.Mattermost.KVGet(key)
		if appErr != nil {
			return nil, nil, errors.Wrap(appErr, "unable to read inital value")
		}

		modifiedBytes, err := modify(initialBytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "modification error")
		}

		return initialBytes, modifiedBytes, nil
	}

	var (
		retryLimit     = 5
		retryWait      = 30 * time.Millisecond
		success        = false
		currentAttempt = 0
	)
	for !success {
		initialBytes, newValue, err := readModify()

		if err != nil {
			return err
		}

		var setError *model.AppError
		success, setError = config.Mattermost.KVCompareAndSet(key, initialBytes, newValue)
		if setError != nil {
			return errors.Wrap(setError, "problem writing value")
		}

		if currentAttempt == 0 && bytes.Equal(initialBytes, newValue) {
			return nil
		}

		currentAttempt++
		if currentAttempt >= retryLimit {
			return errors.New("reached write attempt limit")
		}

		time.Sleep(retryWait)
	}

	return nil
}
