package main

import "github.com/mattermost/mattermost-plugin-confluence/server/serializer"

// Client is the combined interface for all upstream APIs and convenience methods.
type Client interface {
	RESTService
}

// RESTService is the low-level interface for invoking the upstream service.
// Endpoint can be a "short" API URL path, including the version desired, like "v3/user",
// or a fully-qualified URL, with a non-empty scheme.
type RESTService interface {
	GetSelf() (*ConfluenceUser, int, error)
	GetSpaceData(string) (*SpaceResponse, int, error)
	GetUserGroups(*Connection) ([]*UserGroup, int, error)
	GetSpaces() ([]*Spaces, int, error)
	CreatePage(spaceKey string, pageDetails *serializer.PageDetails) (*CreatePageResponse, int, error)
}
