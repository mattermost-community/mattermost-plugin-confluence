package serializer

import (
	"github.com/pkg/errors"
	url2 "net/url"
)

type Subscription struct {
	Alias     string   `json:"alias"`
	BaseURL   string   `json:"baseURL"`
	SpaceKey  string   `json:"spaceKey"`
	Events    []string `json:"events"`
	ChannelID string   `json:"channelID"`
}

func (s *Subscription) IsValid () error {
	if s.Alias == "" {
		return errors.New("Alias can not be empty")
	}
	if s.BaseURL == "" {
		return  errors.New("Base url can not be empty.")
	}
	if _, err := url2.Parse(s.BaseURL); err != nil {
		return  errors.New("Enter a valid url.")
	}
	if s.SpaceKey == "" {
		return  errors.New("Space key can not be empty.")
	}
	return nil
}
