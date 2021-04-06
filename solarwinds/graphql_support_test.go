package solarwinds

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewGraphQLResponse(t *testing.T) {
	responseStr := `
{
  "data": {
    "createOrganizationInvitation": {
      "success": true,
      "code": "200",
      "message": "",
      "invitation": {
        "email": "vB0XMNWacL@foo.com",
        "role": "MEMBER",
        "__typename": "OrganizationInvitation"
      },
      "__typename": "CreateOrganizationInvitationResponse"
    }
  }
}
`
	resp, err := NewGraphQLResponse(strings.NewReader(responseStr), "createOrganizationInvitation")
	assert.NoError(t, err)
	assert.True(t, resp.isSuccess())
}
