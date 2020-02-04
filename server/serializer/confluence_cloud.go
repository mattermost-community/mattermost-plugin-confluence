package serializer

import (
	"encoding/json"
	"io"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
)

type ConfluenceCloudEvent struct {
	UserAccountID string   `json:"userAccountId"`
	AccountType   string   `json:"accountType"`
	UpdateTrigger string   `json:"updateTrigger"`
	Timestamp     int      `json:"timestamp"`
	Comment       *Comment `json:"comment"`
	Page          *Page    `json:"page"`
}

type Page struct {
	CreatorAccountID      string `json:"creatorAccountId"`
	SpaceKey              string `json:"spaceKey"`
	ModificationDate      int64  `json:"modificationDate"`
	LastModifierAccountID string `json:"lastModifierAccountId"`
	Self                  string `json:"self"`
	ID                    int    `json:"id"`
	Title                 string `json:"title"`
	CreationDate          int    `json:"creationDate"`
	ContentTypes          string `json:"contentType"`
	Version               int    `json:"version"`
}

type Comment struct {
	CreatorAccountID      string         `json:"creatorAccountId"`
	SpaceKey              string         `json:"spaceKey"`
	ModificationDate      int64          `json:"modificationDate"`
	LastModifierAccountID string         `json:"lastModifierAccountId"`
	Self                  string         `json:"self"`
	ID                    int            `json:"id"`
	CreationDate          int            `json:"creationDate"`
	ContentTypes          string         `json:"contentType"`
	Version               int            `json:"version"`
	Parent                *Page          `json:"parent"`
	InReplyTo             *ParentComment `json:"inReplyTo"`
}

type ParentComment struct {
	ID string `json:"id"`
}

func ConfluenceCloudEventFromJSON(data io.Reader) *ConfluenceCloudEvent {
	var confluenceServerEvent ConfluenceCloudEvent
	if err := json.NewDecoder(data).Decode(&confluenceServerEvent); err != nil {
		config.Mattermost.LogError("Unable to decode JSON for ConfluenceServerEvent.", "Error", err.Error())
	}
	return &confluenceServerEvent
}
