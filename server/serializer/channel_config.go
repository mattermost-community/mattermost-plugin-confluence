package serializer

import (
	url2 "net/url"

	"github.com/pkg/errors"
)

type Subscription struct {
	Alias     string   `json:"alias"`
	BaseURL   string   `json:"baseURL"`
	SpaceKey  string   `json:"spaceKey"`
	Events    []string `json:"events"`
	ChannelID string   `json:"channelID"`
}

func (s *Subscription) IsValid() error {
	if s.Alias == "" {
		return errors.New("alias can not be empty")
	}
	if s.BaseURL == "" {
		return errors.New("base url can not be empty")
	}
	if _, err := url2.Parse(s.BaseURL); err != nil {
		return errors.New("enter a valid url")
	}
	if s.SpaceKey == "" {
		return errors.New("space key can not be empty")
	}
	return nil
}
