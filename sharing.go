package dropbox

import (
	"encoding/json"
	"net/http"
	"time"
)

// Sharing client.
type Sharing struct {
	*Client
}

// NewSharing client.
func NewSharing(config *Config) *Sharing {
	return &Sharing{
		Client: &Client{
			Config: config,
		},
	}
}

// CreateSharedLinkInput request input.
type CreateSharedLinkInput struct {
	Path     string `json:"path"`
	ShortURL bool   `json:"short_url"`
}

// CreateSharedLinkOutput request output.
type CreateSharedLinkOutput struct {
	URL             string `json:"url"`
	Path            string `json:"path"`
	VisibilityModel struct {
		Tag VisibilityType `json:".tag"`
	} `json:"visibility"`
	Expires time.Time `json:"expires,omitempty"`
	Header  http.Header
}

// VisibilityType determines who can access the link.
type VisibilityType string

// Visibility types supported.
const (
	Public           VisibilityType = "public"
	TeamOnly                        = "team_only"
	Password                        = "password"
	TeamAndPassword                 = "team_and_password"
	SharedFolderOnly                = "shared_folder_only"
)

// CreateSharedLink returns a shared link.
func (c *Sharing) CreateSharedLink(in *CreateSharedLinkInput) (out *CreateSharedLinkOutput, err error) {
	body, hdr, err := c.call("/sharing/create_shared_link", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	if err == nil {
		out.Header = hdr
	}
	return
}

// SharedLinkSettings defines the requested settings for the newly created shared
// link.
type SharedLinkSettings struct {
	RequestedVisibility VisibilityType `json:"requested_visibility,omitempty"`
	LinkPassword        string         `json:"link_password,omitempty"`
	Expires             time.Time      `json:"expires,omitempty"`
}

// CreateSharedLinkWithSettingsInput request input.
type CreateSharedLinkWithSettingsInput struct {
	Path     string              `json:"path"`
	Settings *SharedLinkSettings `json:"settings,omitempty"`
}

// LinkPermissions define the permissions of a shared link.
type LinkPermissions struct {
	CanRevoke       bool `json:"can_revoke"`
	VisibilityModel struct {
		Tag VisibilityType `json:".tag"`
	} `json:"visibility"`
	RevokeFailureReason struct {
		Tag VisibilityType `json:".tag"`
	} `json:"revoke_failure_reason"`
}

// TeamMemberInfo defines information about a team member.
type TeamMemberInfo struct {
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team_info"`
	DisplayName string `json:"display_name"`
	MemberID    string `json:"member_id"`
}

// CreateSharedLinkWithSettingsOutput request output.
type CreateSharedLinkWithSettingsOutput struct {
	Metadata
	LinkPermissions *LinkPermissions `json:"link_permissions"`
	TeamMemberInfo  *TeamMemberInfo  `json:"team_member_info"`
	Header          http.Header
}

// CreateSharedLinkWithSettings returns a shared link.
func (c *Sharing) CreateSharedLinkWithSettings(in *CreateSharedLinkWithSettingsInput) (out *CreateSharedLinkWithSettingsOutput, err error) {
	body, hdr, err := c.call("/sharing/create_shared_link_with_settings", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	if err == nil {
		out.Header = hdr
	}
	return
}
