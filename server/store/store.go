package store

import (
	"bytes"
	"crypto/md5" // #nosec G501
	"crypto/rand"
	"encoding/json"
	"fmt"
	url2 "net/url"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-plugin-confluence/server/util"
	"github.com/mattermost/mattermost-plugin-confluence/server/util/types"
)

const (
	prefixOneTimeSecret             = "ots_" // + unique key that will be deleted after the first verification
	ConfluenceSubscriptionKeyPrefix = "confluence_subs"
	expiryStoreTimeoutSeconds       = 15 * 60
	keyTokenSecret                  = "token_secret"
	keyRSAKey                       = "rsa_key"
	prefixUser                      = "user_"
	AdminMattermostUserID           = "admin"
)

var ErrNotFound = errors.New("not found")

func GetURLSpaceKeyCombinationKey(url, spaceKey string) string {
	u, _ := url2.Parse(url)
	return fmt.Sprintf("%s/%s/%s", ConfluenceSubscriptionKeyPrefix, u.Hostname(), spaceKey)
}

func GetURLPageIDCombinationKey(url, pageID string) string {
	u, _ := url2.Parse(url)
	return fmt.Sprintf("%s/%s/%s", ConfluenceSubscriptionKeyPrefix, u.Hostname(), pageID)
}

func GetSubscriptionKey() string {
	return util.GetKeyHash(ConfluenceSubscriptionKeyPrefix)
}

// from https://github.com/mattermost/mattermost-plugin-jira/blob/master/server/subscribe.go#L625
func AtomicModify(key string, modify func(initialValue []byte) ([]byte, error)) error {
	readModify := func() ([]byte, []byte, error) {
		initialBytes, appErr := config.Mattermost.KVGet(key)
		if appErr != nil {
			return nil, nil, errors.Wrap(appErr, "unable to read initial value")
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

func keyWithInstanceID(instanceID, key types.ID) string {
	h := md5.New() // #nosec G401
	fmt.Fprintf(h, "%s/%s", instanceID, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hashkey(prefix, key string) string {
	h := md5.New() // #nosec G401
	_, _ = h.Write([]byte(key))
	return fmt.Sprintf("%s%x", prefix, h.Sum(nil))
}

func get(key string, v interface{}) (returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr, "failed to get from store")
	}()

	data, appErr := config.Mattermost.KVGet(key)
	if appErr != nil {
		return appErr
	}
	if data == nil {
		return ErrNotFound
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func set(key string, v interface{}) (returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr, fmt.Sprintf("failed to set the value of given key: %s", key))
	}()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	appErr := config.Mattermost.KVSet(key, data)
	if appErr != nil {
		return appErr
	}
	return nil
}

func Load(key string) ([]byte, error) {
	data, appErr := config.Mattermost.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}
	if data == nil {
		return nil, errors.Wrap(ErrNotFound, key)
	}
	return data, nil
}

// revive:disable-next-line:exported
func StoreOAuth2State(state string) error {
	appErr := config.Mattermost.KVSetWithExpiry(
		hashkey(prefixOneTimeSecret, state), []byte(state), expiryStoreTimeoutSeconds)
	if appErr != nil {
		return errors.WithMessage(appErr, "failed to store state "+state)
	}
	return nil
}

func VerifyOAuth2State(state string) error {
	data, appErr := config.Mattermost.KVGet(hashkey(prefixOneTimeSecret, state))
	if appErr != nil {
		return errors.WithMessage(appErr, "failed to load state "+state)
	}

	if string(data) != state {
		return errors.New("invalid oauth state, please try again")
	}

	return nil
}

func EnsureAuthTokenEncryptionSecret() (secret []byte, returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr, "failed to ensure auth token secret")
	}()

	// nil, nil == NOT_FOUND, if we don't already have a key, try to generate one.
	secret, appErr := config.Mattermost.KVGet(keyTokenSecret)
	if appErr != nil {
		return nil, appErr
	}

	if len(secret) == 0 {
		newSecret := make([]byte, 32)
		_, err := rand.Reader.Read(newSecret)
		if err != nil {
			return nil, err
		}

		appErr = config.Mattermost.KVSet(keyTokenSecret, newSecret)
		if appErr != nil {
			return nil, appErr
		}
		secret = newSecret
		config.Mattermost.LogDebug("Stored: auth token secret")
	}

	// If we weren't able to save a new key above, another server must have beat us to it. Get the
	// key from the database, and if that fails, error out.
	if secret == nil {
		secret, appErr = config.Mattermost.KVGet(keyTokenSecret)
		if appErr != nil {
			return nil, appErr
		}
	}

	return secret, nil
}

// revive:disable-next-line:exported
func StoreConnection(instanceID, mattermostUserID types.ID, connection *types.Connection, pluginVersion string) (returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr,
			fmt.Sprintf("failed to store connection, mattermostUserID:%s, Confluence user:%s", mattermostUserID, connection.DisplayName))
	}()

	connection.PluginVersion = pluginVersion

	err := set(keyWithInstanceID(instanceID, mattermostUserID), connection)
	if err != nil {
		return err
	}

	err = set(keyWithInstanceID(instanceID, connection.ConfluenceAccountID()), mattermostUserID)
	if err != nil {
		return err
	}

	// Also store AccountID -> mattermostUserID because Confluence Cloud is deprecating the name field
	// https://developer.atlassian.com/cloud/Confluence/platform/api-changes-for-user-privacy-announcement/
	err = set(keyWithInstanceID(instanceID, connection.ConfluenceAccountID()), mattermostUserID)
	if err != nil {
		return err
	}

	config.Mattermost.LogDebug("Stored: connection, keys:\n\t%s (%s): %+v\n\t%s (%s): %s",
		keyWithInstanceID(instanceID, mattermostUserID), mattermostUserID, connection,
		keyWithInstanceID(instanceID, connection.ConfluenceAccountID()), connection.ConfluenceAccountID(), mattermostUserID)

	return nil
}

func LoadConnection(instanceID, mattermostUserID types.ID, pluginVersion string) (*types.Connection, error) {
	c := &types.Connection{}
	err := get(keyWithInstanceID(instanceID, mattermostUserID), c)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to load connection for Mattermost user ID:%q, Confluence:%q", mattermostUserID, instanceID)
	}
	c.PluginVersion = pluginVersion
	return c, nil
}

func DeleteConnection(instanceID, mattermostUserID types.ID, pluginVersion string) (returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr,
			fmt.Sprintf("failed to delete user, mattermostUserId:%s", mattermostUserID))
	}()

	c, err := LoadConnection(instanceID, mattermostUserID, pluginVersion)
	if err != nil {
		return err
	}

	// Check for whether the admin token stored for each confluenceURL is of the current user or not. If it is then delete that admin connection also
	if c.IsAdmin {
		adminConnection, lErr := LoadConnection(instanceID, AdminMattermostUserID, pluginVersion)
		if lErr != nil {
			return lErr
		}

		// Check if both the tokens are same or not
		if c.OAuth2Token == adminConnection.OAuth2Token {
			if err = DeleteConnectionFromKVStore(instanceID, AdminMattermostUserID, c); err != nil {
				return err
			}
		}
	}

	err = DeleteConnectionFromKVStore(instanceID, mattermostUserID, c)
	if err != nil {
		return err
	}

	return nil
}

func DeleteConnectionFromKVStore(instanceID, mattermostUserID types.ID, c *types.Connection) error {
	appErr := config.Mattermost.KVDelete(keyWithInstanceID(instanceID, mattermostUserID))
	if appErr != nil {
		return appErr
	}

	appErr = config.Mattermost.KVDelete(keyWithInstanceID(instanceID, c.ConfluenceAccountID()))
	if appErr != nil {
		return appErr
	}

	config.Mattermost.LogDebug("Deleted: user, keys: %s(%s), %s(%s)",
		mattermostUserID, keyWithInstanceID(instanceID, mattermostUserID),
		c.ConfluenceAccountID(), keyWithInstanceID(instanceID, c.ConfluenceAccountID()))
	return nil
}

func LoadUser(mattermostUserID types.ID) (*types.User, error) {
	user := types.NewUser(mattermostUserID)
	key := hashkey(prefixUser, mattermostUserID.String())
	err := get(key, user)
	if err != nil {
		return nil, errors.WithMessage(err,
			fmt.Sprintf("failed to load confluence user for mattermostUserId:%s", mattermostUserID))
	}
	return user, nil
}

// revive:disable-next-line:exported
func StoreUser(user *types.User, pluginVersion string) (returnErr error) {
	defer func() {
		if returnErr == nil {
			return
		}
		returnErr = errors.WithMessage(returnErr,
			fmt.Sprintf("failed to store user, mattermostUserId:%s", user.MattermostUserID))
	}()

	user.PluginVersion = pluginVersion

	key := hashkey(prefixUser, user.MattermostUserID.String())
	err := set(key, user)
	if err != nil {
		return err
	}

	config.Mattermost.LogDebug("Stored: user %s key:%s: connected to:%q", user.MattermostUserID, key, user.InstanceURL)
	return nil
}
