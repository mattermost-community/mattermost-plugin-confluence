package util

import (
	"fmt"
	url2 "net/url"
)

func GetKey(url, spaceKey string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", u.Hostname(), spaceKey), nil
}
