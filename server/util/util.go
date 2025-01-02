package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
)

const (
	Script = "script"
	Style  = "style"
)

// GetKeyHash can be used to create a hash from a string
func GetKeyHash(key string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func SplitArgs(s string) ([]string, error) {
	indexes := regexp.MustCompile("\"").FindAllStringIndex(s, -1)
	if len(indexes)%2 != 0 {
		return []string{}, errors.New("quotes not closed")
	}

	indexes = append([][]int{{0, 0}}, indexes...)

	if indexes[len(indexes)-1][1] < len(s) {
		indexes = append(indexes, [][]int{{len(s), 0}}...)
	}

	var args []string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i][1]
		end := min(len(s), indexes[i+1][0])

		if i%2 == 0 {
			args = append(args, strings.Split(strings.Trim(s[start:end], " "), " ")...)
		} else {
			args = append(args, s[start:end])
		}
	}

	cleanedArgs := make([]string, len(args))
	count := 0

	for _, arg := range args {
		if arg != "" {
			cleanedArgs[count] = strings.TrimSpace(arg)
			count++
		}
	}

	return cleanedArgs[0:count], nil
}

func GetPluginKey() string {
	var regexpNonAlnum = regexp.MustCompile("[^a-zA-Z0-9]+")
	return "mattermost_" + regexpNonAlnum.ReplaceAllString(GetSiteURL(), "_")
}

func GetPluginURLPath() string {
	return "/plugins/" + config.PluginName + "/api/v1"
}

func GetPluginURL() string {
	return strings.TrimRight(GetSiteURL(), "/") + GetPluginURLPath()
}

func GetSiteURL() string {
	ptr := config.Mattermost.GetConfig().ServiceSettings.SiteURL
	if ptr == nil {
		return ""
	}
	return *ptr
}

func GetAtlassianConnectURLPath() string {
	return "/atlassian-connect.json?secret=" + url.QueryEscape(config.GetConfig().Secret)
}

func GetConfluenceServerWebhookURLPath() string {
	return "/server/webhook?secret=" + url.QueryEscape(config.GetConfig().Secret)
}

func IsSystemAdmin(userID string) bool {
	user, appErr := config.Mattermost.GetUser(userID)
	if appErr != nil {
		return false
	}
	return user.IsInRole(model.SystemAdminRoleId)
}

func Deduplicate(a []string) []string {
	check := make(map[string]int)
	result := make([]string, 0)
	for _, val := range a {
		check[val] = 1
	}

	for key := range check {
		result = append(result, key)
	}

	return result
}
