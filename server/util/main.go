package util

import (
	"errors"
	"fmt"
	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	url2 "net/url"
	"regexp"
	"strings"
)

func GetURLSpaceKeyCombinationKey(url, spaceKey string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", config.ConfluenceSubscriptionKeyPrefix, u.Hostname(), spaceKey), nil
}

func GetChannelSubscriptionKey(channelID string) string {
	return fmt.Sprintf("%s/%s", config.ConfluenceSubscriptionKeyPrefix, channelID)
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
