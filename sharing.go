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
	Expires             time.Time      `json:"-"`
}

// CreateSharedLinkWithSettingsInput request input.
type CreateSharedLinkWithSettingsInput struct {
	Path     string              `json:"path"`
	Settings *SharedLinkSettings `json:"settings,omitempty"`
}

// SharedLinkAccessFailureReason determines the access failure reason.
type SharedLinkAccessFailureReason string

// SharedLinkAccessFailureReason types supported.
const (
	LoginRequired       SharedLinkAccessFailureReason = "login_required"
	EmailVerifyRequired                               = "email_verify_required"
	PasswordRequired                                  = "password_required"
	ReasonTeamOnly                                    = "team_only"
	ReasonOwnerOnly                                   = "owner_only"
)

// LinkPermissions defines the permissions of a shared link.
type LinkPermissions struct {
	CanRevoke               bool `json:"can_revoke"`
	ResolvedVisibilityModel struct {
		Tag VisibilityType `json:".tag"`
	} `json:"resolved_visibility"`
	RequestedVisibilityModel struct {
		Tag VisibilityType `json:".tag"`
	} `json:"requested_visibility"`
	RevokeFailureReasonModel struct {
		Tag SharedLinkAccessFailureReason `json:".tag"`
	} `json:"revoke_failure_reason"`
}

// Team provides information about a team.
type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TeamMemberInfo provides information about a team member.
type TeamMemberInfo struct {
	Team        *Team  `json:"team_info"`
	DisplayName string `json:"display_name"`
	MemberID    string `json:"member_id"`
}

// CreateSharedLinkWithSettingsOutput request output.
type CreateSharedLinkWithSettingsOutput struct {
	Metadata
	URL             string           `json:"url"`
	LinkPermissions *LinkPermissions `json:"link_permissions"`
	TeamMemberInfo  *TeamMemberInfo  `json:"team_member_info"`
	Header          http.Header
}

// CreateSharedLinkWithSettings returns a shared link.
func (c *Sharing) CreateSharedLinkWithSettings(in *CreateSharedLinkWithSettingsInput) (out *CreateSharedLinkWithSettingsOutput, err error) {
	type settings struct {
		RequestedVisibility VisibilityType `json:"requested_visibility,omitempty"`
		LinkPassword        string         `json:"link_password,omitempty"`
		Expires             string         `json:"expires,omitempty"`
	}
	var s *settings
	if in.Settings != nil {
		s = &settings{
			RequestedVisibility: in.Settings.RequestedVisibility,
			LinkPassword:        in.Settings.LinkPassword,
		}
		if !in.Settings.Expires.IsZero() {
			s.Expires = in.Settings.Expires.Format("2006-01-02T15:04:05Z")
		}
	}
	in2 := &struct {
		Path     string    `json:"path"`
		Settings *settings `json:"settings,omitempty"`
	}{
		in.Path,
		s,
	}
	body, hdr, err := c.call("/sharing/create_shared_link_with_settings", in2)
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

// RevokeSharedLinkInput request input.
type RevokeSharedLinkInput struct {
	URL string `json:"url"`
}

// RevokeSharedLink revokes a shared link.
func (c *Sharing) RevokeSharedLink(in *RevokeSharedLinkInput) error {
	body, _, err := c.call("/sharing/revoke_shared_link", in)
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

// ListSharedLinksInput request input.
type ListSharedLinksInput struct {
	Path       string `json:"path,omitempty"`
	Cursor     string `json:"cursor,omitempty"`
	DirectOnly bool   `json:"direct_only,omitempty"`
}

// SharedLinkMetadata provides information about a shared link.
type SharedLinkMetadata struct {
	Tag                  string           `json:".tag"`
	URL                  string           `json:"url"`
	Name                 string           `json:"name"`
	LinkPermissions      *LinkPermissions `json:"link_permissions"`
	PathLower            string           `json:"path_lower"`
	ClientModified       time.Time        `json:"client_modified"`
	ServerModified       time.Time        `json:"server_modified"`
	Rev                  string           `json:"rev"`
	Size                 uint64           `json:"size"`
	ID                   string           `json:"id"`
	TeamMemberInfo       *TeamMemberInfo  `json:"team_member_info"`
	ContentOwnerTeamInfo *Team            `json:"content_owner_team_info"`
}

// ListSharedLinksOutput request output.
type ListSharedLinksOutput struct {
	Cursor  string `json:"cusror"`
	HasMore bool   `json:"has_more"`
	Links   []*SharedLinkMetadata
	Header  http.Header
}

// ListSharedLinks gets the list of shared links of the user.
func (c *Sharing) ListSharedLinks(in *ListSharedLinksInput) (*ListSharedLinksOutput, error) {
	body, hdr, err := c.call("/sharing/list_shared_links", in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var out *ListSharedLinksOutput
	err = json.NewDecoder(body).Decode(&out)
	if err == nil {
		out.Header = hdr
	}
	return out, err
}
