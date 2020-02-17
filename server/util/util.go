package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
)

// GetKeyHash can be used to create a hash from a string
func GetKeyHash(key string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
		end := Min(len(s), indexes[i+1][0])

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
			cleanedArgs[count] = arg
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
	return "/api/v1/atlassian-connect.json?secret=" + url.QueryEscape(config.GetConfig().Secret)
}
