package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/pkg/errors"
)

func getKeyHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func Set(key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if appErr := config.Mattermost.KVSet(getKeyHash(key), bytes); appErr != nil {
		return errors.New(appErr.Error())
	}
	return nil
}

func Get(key string, data interface{}) error {
	bytes, appErr := config.Mattermost.KVGet(getKeyHash(key))
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
