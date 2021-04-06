package solarwinds

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	inviteUserResponseStr = `
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
	revokePendingInvitationResponseStr = `
{
  "data": {
    "deleteOrganizationInvitation": {
      "success": true,
      "code": "200",
      "message": "",
      "__typename": "MutationResponse"
    }
  }
}
`
	resendInvitationResponseStr = `
{
  "data": {
    "resendOrganizationInvitation": {
      "success": true,
      "code": "200",
      "message": "",
      "__typename": "MutationResponse"
    }
  }
}
`
	listInvitationResponseStr = `
{
  "data": {
    "user": {
      "id": "106586091288584192",
      "currentOrganization": {
        "id": "106269109693582336",
        "invitations": [
          {
            "email": "5et54o0OtS@foo.com",
            "role": "MEMBER",
            "date": "2021-03-25T02:36:48Z",
            "products": [
              {
                "name": "APPOPTICS",
                "role": "MEMBER",
                "access": true,
                "__typename": "ProductAccess"
              }
            ],
            "__typename": "OrganizationInvitation"
          },
          {
            "email": "0JTELJv5YA@foo.com",
            "role": "MEMBER",
            "date": "2021-03-25T02:37:25Z",
            "products": [
              {
                "name": "APPOPTICS",
                "role": "MEMBER",
                "access": true,
                "__typename": "ProductAccess"
              }
            ],
            "__typename": "OrganizationInvitation"
          }
        ],
        "__typename": "Organization"
      },
      "__typename": "AuthenticatedUser"
    }
  }
}
`
)

func TestInviteUser(t *testing.T) {
	setup()
	defer teardown()

	invitation := Invitation{
		Email: RandString(8) + "@foo.com",
		Role:  "Member",
		Products: []Product{
			{
				Name: "AppOptics",
				Role: "Admin",
			},
			{
				Name: "Loggly",
				Role: "User",
			},
		},
	}
	input := inviteUserVars{
		Input: invitation,
	}
	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, inviteUserOp, graphQLReq.OperationName)
		assert.Equal(t, inviteUserQuery, graphQLReq.Query)
		actualVars := inviteUserVars{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, input, actualVars)

		_, _ = fmt.Fprint(w, inviteUserResponseStr)
	})
	err := client.InvitationService.Create(invitation)
	assert.NoError(t, err)
}

func TestRevokePendingInvitation(t *testing.T) {
	setup()
	defer teardown()

	email := RandString(8) + "@foo.com"
	variables := revokeInvitationVars{
		Email: email,
	}
	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, revokeInvitationOp, graphQLReq.OperationName)
		assert.Equal(t, revokeInvitationQuery, graphQLReq.Query)
		actualVars := revokeInvitationVars{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, variables, actualVars)

		_, _ = fmt.Fprint(w, revokePendingInvitationResponseStr)
	})
	err := client.InvitationService.Revoke(email)
	assert.NoError(t, err)
}

func TestResendInvitation(t *testing.T) {
	setup()
	defer teardown()

	email := RandString(8) + "@foo.com"
	variables := resendInvitationVars{
		Email: email,
	}
	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, resendInvitationOp, graphQLReq.OperationName)
		assert.Equal(t, resendInvitationQuery, graphQLReq.Query)
		actualVars := resendInvitationVars{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, variables, actualVars)

		_, _ = fmt.Fprint(w, resendInvitationResponseStr)
	})
	err := client.InvitationService.Resend(email)
	assert.NoError(t, err)
}

func TestListInvitation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, listInvitationOp, graphQLReq.OperationName)
		assert.Equal(t, listInvitationQuery, graphQLReq.Query)

		_, _ = fmt.Fprint(w, listInvitationResponseStr)
	})
	invitationList, err := client.InvitationService.List()
	assert.NoError(t, err)
	invitations := invitationList.Organization.Invitations
	assert.Equal(t, len(invitations), 2)
}
