package acceptance

import (
	"github.com/nordcloud/go-pingdom/solarwinds"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	solarwindsClient        *solarwinds.Client
	runSolarwindsAcceptance bool
)

func init() {
	if os.Getenv("SOLARWINDS_ACCEPTANCE") == "1" {
		runSolarwindsAcceptance = true
		config := solarwinds.ClientConfig{
			Username: os.Getenv("SOLARWINDS_USER"),
			Password: os.Getenv("SOLARWINDS_PASSWD"),
		}
		solarwindsClient, _ = solarwinds.NewClient(config)
		err := solarwindsClient.Init()
		if err != nil {
			panic(err)
		}
	}
}

func TestInvitations(t *testing.T) {
	if !runSolarwindsAcceptance {
		t.Skip()
	}
	email := solarwinds.RandString(10) + "@foo.com"
	invitationService := solarwindsClient.InvitationService
	err := invitationService.Create(solarwinds.Invitation{
		Email: email,
		Role:  "MEMBER",
		Products: []solarwinds.Product{
			{
				Name: "APPOPTICS",
				Role: "MEMBER",
			},
		},
	})
	assert.NoError(t, err)

	invitationList, err := invitationService.List()
	assert.NoError(t, err)
	assert.True(t, len(invitationList.Organization.Invitations) > 0)

	err = invitationService.Resend(email)
	assert.NoError(t, err)

	err = invitationService.Revoke(email)
	assert.NoError(t, err)

	err = invitationService.Resend(email)
	assert.Error(t, err)
}

func TestActiveUsers(t *testing.T) {
	if !runSolarwindsAcceptance {
		t.Skip()
	}

	userService := solarwindsClient.ActiveUserService
	currentUserEmail := os.Getenv("SOLARWINDS_USER")

	userList, err := userService.List()
	assert.NoError(t, err)
	var currentMember *solarwinds.OrganizationMember
	for _, member := range userList.Organization.Members {
		if currentUserEmail == member.User.Email {
			currentMember = &member
			break
		}
	}
	if currentMember == nil {
		t.Errorf("current member is nil")
	} else {
		singleUser, err := userService.Get(currentMember.User.Id)
		assert.NoError(t, err)
		assert.Equal(t, currentMember.User.Email, singleUser.Organization.Members[0].User.Email)

		containsRole := func(member *solarwinds.OrganizationMember, app string, role string) bool {
			for _, product := range member.Products {
				if product.Name == app && product.Role == role {
					return true
				}
			}
			return false
		}
		updateAddRole := solarwinds.UpdateActiveUserRequest{
			UserId: currentMember.User.Id,
			Role:   currentMember.Role,
			Products: []solarwinds.Product{
				{
					Name: "LOGGLY",
					Role: "MEMBER",
				},
			},
		}
		assert.True(t, containsRole(currentMember, "LOGGLY", "NO_ACCESS"))
		err = userService.Update(updateAddRole)
		assert.NoError(t, err)

		singleUser, err = userService.Get(currentMember.User.Id)
		assert.NoError(t, err)
		assert.True(t, containsRole(&singleUser.Organization.Members[0], "LOGGLY", "MEMBER"))

		updateRevokeRole := solarwinds.UpdateActiveUserRequest{
			UserId: currentMember.User.Id,
			Role:   currentMember.Role,
			Products: []solarwinds.Product{
				{
					Name: "LOGGLY",
					Role: "NO_ACCESS",
				},
			},
		}
		err = userService.Update(updateRevokeRole)
		assert.NoError(t, err)
		singleUser, _ = userService.Get(currentMember.User.Id)
		assert.True(t, containsRole(&singleUser.Organization.Members[0], "LOGGLY", "NO_ACCESS"))
	}
}

func TestUsers(t *testing.T) {
	if !runSolarwindsAcceptance {
		t.Skip()
	}
	email := solarwinds.RandString(10) + "@foo.com"
	userToCreate := solarwinds.User{
		Email: email,
		Role:  "MEMBER",
		Products: []solarwinds.Product{
			{
				Name: "APPOPTICS",
				Role: "MEMBER",
			},
		},
	}
	userService := solarwindsClient.UserService
	err := userService.Create(userToCreate)
	assert.NoError(t, err)

	user, err := userService.Retrieve(email)
	assert.NoError(t, err)
	assert.Equal(t, userToCreate, *user)

	userUpdate := userToCreate
	userUpdate.Products = []solarwinds.Product{
		{
			Name: "APPOPTICS",
			Role: "ADMIN",
		},
	}
	err = userService.Update(userUpdate)
	assert.NoError(t, err)

	userAfterUpdate, err := userService.Retrieve(email)
	assert.NoError(t, err)
	assert.Equal(t, userUpdate, *userAfterUpdate)

	err = userService.Delete(email)
	assert.NoError(t, err)

	userAfterDelete, err := userService.Retrieve(email)
	assert.NoError(t, err)
	assert.Nil(t, userAfterDelete)
}
