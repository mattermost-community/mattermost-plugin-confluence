package serializer

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
	CreatorAccountID      string `json:"creatorAccountId"`
	SpaceKey              string `json:"spaceKey"`
	ModificationDate      int64  `json:"modificationDate"`
	LastModifierAccountID string `json:"lastModifierAccountId"`
	Self                  string `json:"self"`
	ID                    int    `json:"id"`
	CreationDate          int    `json:"creationDate"`
	ContentTypes          string `json:"contentType"`
	Version               int    `json:"version"`
	Parent                *Page  `json:"parent"`
}
