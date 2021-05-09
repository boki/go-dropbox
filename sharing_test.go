package dropbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharing_CreateSharedLink(t *testing.T) {
	c := client()
	out, err := c.Sharing.CreateSharedLink(&CreateSharedLinkInput{
		Path:     "/hello.txt",
		ShortURL: true,
	})

	assert.NoError(t, err, "error sharing file")
	assert.Equal(t, "/hello.txt", out.Path)
}

func TestSharing_RevokeSharedLink(t *testing.T) {
	c := client()
	out, err := c.Sharing.CreateSharedLink(&CreateSharedLinkInput{
		Path: "/hello.txt",
	})
	assert.NoError(t, err, "error sharing file")

	err = c.Sharing.RevokeSharedLink(&RevokeSharedLinkInput{URL: out.URL})
	assert.NoError(t, err, "error revoking shared file")
}
