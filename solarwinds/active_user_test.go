package solarwinds

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	listActiveUserResponseStr = `
{
  "data": {
    "user": {
      "id": "106586091288584192",
      "currentOrganization": {
        "id": "106269109693582336",
        "members": [
          {
            "user": {
              "id": "23285292452068352",
              "firstName": "IT",
              "lastName": "Nordcloud",
              "email": "foo@nordcloud.com",
              "lastLogin": "2021-03-23T07:17:48Z",
              "__typename": "User"
            },
            "role": "ADMIN",
            "products": [
              {
                "name": "APPOPTICS",
                "access": false,
                "role": "NO_ACCESS",
                "__typename": "ProductAccess"
              },
              {
                "name": "LOGGLY",
                "access": false,
                "role": "NO_ACCESS",
                "__typename": "ProductAccess"
              },
              {
                "name": "PINGDOM",
                "access": true,
                "role": "ADMIN",
                "__typename": "ProductAccess"
              }
            ],
            "__typename": "OrganizationMember"
          },
          {
            "user": {
              "id": "74914272581727232",
              "firstName": "Nordcloud",
              "lastName": "MC-Tooling",
              "email": "bar@nordcloud.com",
              "lastLogin": "2021-03-24T23:04:56Z",
              "__typename": "User"
            },
            "role": "ADMIN",
            "products": [
              {
                "name": "APPOPTICS",
                "access": false,
                "role": "NO_ACCESS",
                "__typename": "ProductAccess"
              },
              {
                "name": "LOGGLY",
                "access": false,
                "role": "NO_ACCESS",
                "__typename": "ProductAccess"
              },
              {
                "name": "PINGDOM",
                "access": false,
                "role": "NO_ACCESS",
                "__typename": "ProductAccess"
              }
            ],
            "__typename": "OrganizationMember"
          }
        ],
        "__typename": "Organization"
      },
      "__typename": "AuthenticatedUser"
    }
  }
}
`
	getActiveUserResponseStr = `
{
  "data": {
    "user": {
      "id": "106586091288584192",
      "currentOrganization": {
        "id": "106269109693582336",
        "members": [
          {
            "id": "106586091288584192",
            "user": {
              "email": "foo@nordcloud.com",
              "__typename": "User"
            },
            "role": "ADMIN",
            "products": [
              {
                "name": "APPOPTICS",
                "role": "MEMBER",
                "access": true,
                "__typename": "ProductAccess"
              },
              {
                "name": "LOGGLY",
                "role": "NO_ACCESS",
                "access": false,
                "__typename": "ProductAccess"
              },
              {
                "name": "PINGDOM",
                "role": "ADMIN",
                "access": true,
                "__typename": "ProductAccess"
              }
            ],
            "__typename": "OrganizationMember"
          }
        ],
        "__typename": "Organization"
      },
      "__typename": "AuthenticatedUser"
    }
  }
}
`
	updateActiveUserResponseStr = `
{
  "data": {
    "updateMemberRoles": {
      "code": "200",
      "success": true,
      "message": "",
      "__typename": "UpdateMemberRolesResponse"
    }
  }
}
`
)

func TestListActiveUsers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, listActiveUserOp, graphQLReq.OperationName)
		assert.Equal(t, listActiveUserQuery, graphQLReq.Query)

		_, _ = fmt.Fprint(w, listActiveUserResponseStr)
	})
	userList, err := client.ActiveUserService.List()
	assert.NoError(t, err)
	members := userList.Organization.Members
	assert.Equal(t, "106586091288584192", userList.OwnerUserId)
	assert.Equal(t, len(members), 2)

	user, err := client.ActiveUserService.GetByEmail(members[1].User.Email)
	assert.NoError(t, err)
	assert.Equal(t, members[1], *user)
}

func TestGetActiveUser(t *testing.T) {
	setup()
	defer teardown()

	userId := "106586091288584192"
	input := getActiveUserVars{
		UserId: userId,
	}
	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, getActiveUserOp, graphQLReq.OperationName)
		assert.Equal(t, getActiveUserQuery, graphQLReq.Query)
		actualVars := getActiveUserVars{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, input, actualVars)
		_, _ = fmt.Fprint(w, getActiveUserResponseStr)
	})
	userList, err := client.ActiveUserService.Get("106586091288584192")
	assert.NoError(t, err)
	members := userList.Organization.Members
	assert.Equal(t, len(members), 1)
	member := members[0]
	assert.Equal(t, "foo@nordcloud.com", member.User.Email)
}

func TestUpdateActiveUser(t *testing.T) {
	setup()
	defer teardown()

	update := UpdateActiveUserRequest{
		UserId: "106586091288584192",
		Role:   "ADMIN",
		Products: []Product{
			{
				Name: "APPOPTICS",
				Role: "MEMBER",
			},
		},
	}
	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, updateActiveUserOp, graphQLReq.OperationName)
		assert.Equal(t, updateActiveUserQuery, graphQLReq.Query)
		actualVars := UpdateActiveUserRequest{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, update, actualVars)
		_, _ = fmt.Fprint(w, updateActiveUserResponseStr)
	})
	err := client.ActiveUserService.Update(update)
	assert.NoError(t, err)
}
