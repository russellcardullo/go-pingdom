package solarwinds

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// TODO: These values are from the mock response strings. Ideally they should be randomly generated on each test run.
const (
	activeUserEmail   = "foo@nordcloud.com"
	activeUserId      = "23285292452068352"
	nonExistUserEmail = "other@nordcloud.com"
	pendingUserEmail  = "5et54o0OtS@foo.com"
)

func TestRetrieveUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)

		var responseStr string
		switch graphQLReq.OperationName {
		case listActiveUserOp:
			responseStr = listActiveUserResponseStr
		case listInvitationOp:
			responseStr = listInvitationResponseStr
		default:
			panic("not supposed to reach here")
		}
		_, _ = fmt.Fprint(w, responseStr)
	})
	userService := client.UserService
	user, err := userService.Retrieve(activeUserEmail)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = userService.Retrieve(nonExistUserEmail)
	assert.NoError(t, err)
	assert.Nil(t, user)

	user, err = userService.Retrieve(pendingUserEmail)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestCreateUser(t *testing.T) {
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
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)
		assert.Equal(t, inviteUserOp, graphQLReq.OperationName)
		assert.Equal(t, inviteUserQuery, graphQLReq.Query)
		actualVars := inviteUserVars{}
		_ = Convert(&graphQLReq.Variables, &actualVars)
		assert.Equal(t, input, actualVars)

		_, _ = fmt.Fprint(w, inviteUserResponseStr)
	})
	err := client.UserService.Create(invitation)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	setup()
	defer teardown()

	update := User{
		Role: "ADMIN",
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

		switch graphQLReq.OperationName {
		case listActiveUserOp:
			_, _ = fmt.Fprint(w, listActiveUserResponseStr)
		case updateActiveUserOp:
			assert.Equal(t, updateActiveUserQuery, graphQLReq.Query)
			actualVars := UpdateActiveUserRequest{}
			_ = Convert(&graphQLReq.Variables, &actualVars)
			assert.Equal(t, activeUserId, actualVars.UserId)
			_, _ = fmt.Fprint(w, updateActiveUserResponseStr)
		case revokeInvitationOp:
			_, _ = fmt.Fprint(w, revokePendingInvitationResponseStr)
		case inviteUserOp:
			_, _ = fmt.Fprint(w, inviteUserResponseStr)
		case listInvitationOp:
			_, _ = fmt.Fprint(w, listInvitationResponseStr)
		default:
			t.Errorf("should not have op: %v", graphQLReq.OperationName)
		}
	})

	userService := client.UserService
	update.Email = activeUserEmail
	err := userService.Update(update)
	assert.NoError(t, err)

	update.Email = nonExistUserEmail
	err = userService.Update(update)
	assert.Error(t, err)

	update.Email = pendingUserEmail
	err = userService.Update(update)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(graphQLEndpoint, func(w http.ResponseWriter, r *http.Request) {
		graphQLReq := GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&graphQLReq)

		switch graphQLReq.OperationName {
		case revokeInvitationOp:
			assert.Equal(t, revokeInvitationQuery, graphQLReq.Query)
			actualVars := revokeInvitationVars{}
			_ = Convert(&graphQLReq.Variables, &actualVars)
			if actualVars.Email == pendingUserEmail {
				_, _ = fmt.Fprint(w, revokePendingInvitationResponseStr)
			} else {
				_, _ = fmt.Fprint(w, `
{
  "data": {
    "deleteOrganizationInvitation": {
      "success": false,
      "code": "500",
      "message": "user not exist",
      "__typename": "MutationResponse"
    }
  }
}
`)
			}
		case listActiveUserOp:
			_, _ = fmt.Fprint(w, listActiveUserResponseStr)
		default:
			t.Errorf("should not have op: %v", graphQLReq.OperationName)
		}
	})

	userService := client.UserService
	err := userService.Delete(pendingUserEmail)
	assert.NoError(t, err)

	err = userService.Delete(activeUserEmail)
	assert.Error(t, err)
	assert.Equal(t, ErrCodeDeleteActiveUserException, err.(*ClientError).StatusCode)

	err = userService.Delete(nonExistUserEmail)
	assert.Error(t, err)
}
