package serializer

import (
	"encoding/json"
	"io"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
)

const (
	ConfluenceContentTypePage     = "page"
	ConfluenceContentTypeBlogPost = "blogpost"
	ConfluenceContentTypeComment  = "comment"
)

type ConfluenceServerUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

type ConfluenceServerSpace struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	URL         string `json:"url"`
	Status      string `json:"status"`
}

type ConfluenceServerPageAncestor struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type ConfluenceServerPage struct {
	IsHomePage    bool                           `json:"is_home_page"`
	CreatedAt     int64                          `json:"created_at"`
	Title         string                         `json:"title"`
	Content       string                         `json:"content"`
	HTMLContent   string                         `json:"html_content"`
	IsDeleted     bool                           `json:"is_deleted"`
	ContentType   string                         `json:"content_type"`
	UpdatedAt     int64                          `json:"updated_at"`
	IsDraft       bool                           `json:"is_draft"`
	ID            string                         `json:"id"`
	Ancestors     []ConfluenceServerPageAncestor `json:"ancestors"`
	IsRootLevel   bool                           `json:"is_root_level"`
	ContentID     int                            `json:"content_id"`
	IsIndexable   bool                           `json:"is_indexable"`
	Version       int                            `json:"version"`
	CreatedBy     ConfluenceServerUser           `json:"created_by"`
	URL           string                         `json:"url"`
	Labels        []string                       `json:"labels"`
	TinyURL       string                         `json:"tiny_url"`
	UpdatedBy     ConfluenceServerUser           `json:"updated_by"`
	EditURL       string                         `json:"edit_url"`
	IsUnpublished bool                           `json:"is_unpublished"`
	Position      interface{}                    `json:"position"`
	Excerpt       string                         `json:"excerpt"`
	IsCurrent     bool                           `json:"is_current"`
}

type ConfluenceServerParentComment struct {
	CreatedAt   int64                `json:"created_at"`
	Title       *string              `json:"title"`
	Version     int                  `json:"version"`
	CreatedBy   ConfluenceServerUser `json:"created_by"`
	URL         string               `json:"url"`
	Content     string               `json:"content"`
	Labels      []string             `json:"labels"`
	HTMLContent string               `json:"html_content"`
	ContentType string               `json:"content_type"`
	UpdatedAt   int64                `json:"updated_at"`
	UpdatedBy   ConfluenceServerUser `json:"updated_by"`
	ID          string               `json:"id"`
	Excerpt     string               `json:"excerpt"`
}

type ConfluenceServerComment struct {
	ParentComment    *ConfluenceServerParentComment `json:"parent"`
	DisplayTitle     string                         `json:"display_title"`
	IsInlineComment  bool                           `json:"is_inline_comment"`
	CreatedAt        int64                          `json:"created_at"`
	ThreadChangeDate int64                          `json:"thread_change_date"`
	DescendantsCount int                            `json:"descendants_count"`
	Title            *string                        `json:"title"`
	Version          int                            `json:"version"`
	CreatedBy        ConfluenceServerUser           `json:"created_by"`
	URL              string                         `json:"url"`
	Content          string                         `json:"content"`
	Labels           []string                       `json:"labels"`
	HTMLContent      string                         `json:"html_content"`
	Depth            int                            `json:"depth"`
	ContentType      string                         `json:"content_type"`
	UpdatedAt        int64                          `json:"updated_at"`
	UpdatedBy        ConfluenceServerUser           `json:"updated_by"`
	ID               string                         `json:"id"`
	Excerpt          string                         `json:"excerpt"`
	Status           string                         `json:"status"`
}

type ConfluenceServerBlogPost struct {
	CreatedAt   int64                `json:"created_at"`
	Title       string               `json:"title"`
	Version     int                  `json:"version"`
	CreatedBy   ConfluenceServerUser `json:"created_by"`
	URL         string               `json:"url"`
	Content     string               `json:"content"`
	Labels      []string             `json:"labels"`
	HTMLContent string               `json:"html_content"`
	ContentType string               `json:"content_type"`
	UpdatedAt   int64                `json:"updated_at"`
	UpdatedBy   ConfluenceServerUser `json:"updated_by"`
	ID          string               `json:"id"`
	Excerpt     string               `json:"excerpt"`
}

type ConfluenceServerEvent struct {
	VersionComment string                    `json:"version_comment"`
	IsMinorEdit    bool                      `json:"is_minor_edit"`
	Creator        ConfluenceServerUser      `json:"creator"`
	ContentType    string                    `json:"content_type"`
	BaseURL        string                    `json:"base_url"`
	ContentURL     string                    `json:"content_url"`
	ContainerType  string                    `json:"container_type"`
	Comment        *ConfluenceServerComment  `json:"comment"`
	Page           *ConfluenceServerPage     `json:"page"`
	Blog           *ConfluenceServerBlogPost `json:"blog"`
	Event          string                    `json:"event"`
	Excerpt        string                    `json:"excerpt"`
	User           ConfluenceServerUser      `json:"user"`
	Space          ConfluenceServerSpace     `json:"space"`
	Timestamp      int64                     `json:"timestamp"`
}

func ConfluenceServerEventFromJSON(data io.Reader) *ConfluenceServerEvent {
	var me ConfluenceServerEvent
	if err := json.NewDecoder(data).Decode(&me); err != nil {
		config.Mattermost.LogError("Unable to decode JSON for ConfluenceServerEvent.", "Error", err.Error())
	}
	return &me
}
