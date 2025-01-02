package types

type User struct {
	PluginVersion    string
	MattermostUserID ID `json:"mattermost_user_id"`
	InstanceURL      ID `json:"instance_url,omitempty"`
}

type ConfluenceUser struct {
	AccountID   string `json:"accountId,omitempty"`
	Name        string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type UserGroups struct {
	Groups []*UserGroup `json:"results,omitempty"`
}

type UserGroup struct {
	Name string `json:"name"`
}

type Connection struct {
	ConfluenceUser
	PluginVersion     string
	OAuth2Token       string `json:"token,omitempty"`
	DefaultProjectKey string `json:"default_project_key,omitempty"`
	IsAdmin           bool   `json:"is_admin,omitempty"`
	MattermostUserID  string `json:"mattermost_user_id,omitempty"`
}

func (c *Connection) ConfluenceAccountID() ID {
	if c.AccountID != "" {
		return ID(c.AccountID)
	}

	return ID(c.Name)
}

func NewUser(mattermostUserID ID) *User {
	return &User{
		MattermostUserID: mattermostUserID,
	}
}

func (user *User) AsConfigMap() map[string]interface{} {
	return map[string]interface{}{
		"mattermost_user_id": user.MattermostUserID.String(),
		"instance_url":       user.InstanceURL.String(),
	}
}
