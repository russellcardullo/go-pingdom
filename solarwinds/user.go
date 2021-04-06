package solarwinds

import (
	"fmt"
	"log"
)

type User = Invitation

type UserService struct {
	ActiveUserService *ActiveUserService
	InvitationService *InvitationService
}

// Create will create a new invitation for the user. It is not possible to add user without going through invitation.
func (us *UserService) Create(user User) error {
	return us.InvitationService.Create(user)
}

// Update will first try to update an active user with the given email. If no such user exist, will see if there
// is an invitation with the email, if yes, will revoke the invitation and send a new one. Otherwise, error is returned.
func (us *UserService) Update(update User) error {
	activeUser, _ := us.ActiveUserService.GetByEmail(update.Email)
	if activeUser != nil {
		activeUserUpdate := UpdateActiveUserRequest{
			UserId:   activeUser.User.Id,
			Role:     update.Role,
			Products: update.Products,
		}
		return us.ActiveUserService.Update(activeUserUpdate)
	}

	log.Printf("Will revoke the invitation and send a new one for user: %v", update.Email)
	invitationService := us.InvitationService
	invitationList, err := invitationService.List()
	if err != nil {
		return err
	}
	var invitationFound bool
	for _, invitation := range invitationList.Organization.Invitations {
		if invitation.Email == update.Email {
			invitationFound = true
			if err = invitationService.Revoke(update.Email); err != nil {
				return err
			}
		}
	}
	if !invitationFound {
		return fmt.Errorf("there is no invitation with email: %v", update.Email)
	}
	if err = invitationService.Create(Invitation{
		Email:    update.Email,
		Role:     update.Role,
		Products: update.Products,
	}); err != nil {
		return err
	}
	return nil
}

// Delete will only be effective if it is an invitation. There is no way to delete an active user in Pingdom.
func (us *UserService) Delete(email string) error {
	activeUser, _ := us.ActiveUserService.GetByEmail(email)
	if activeUser != nil {
		return NewErrorAttemptDeleteActiveUser(email)
	}
	err := us.InvitationService.Revoke(email)
	if err != nil {
		return NewNetworkError(err)
	}
	return err
}

// Retrieve return the user information, either it is an invitation or an active user.
func (us *UserService) Retrieve(email string) (*User, error) {
	activeUser, err := us.ActiveUserService.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if activeUser != nil {
		user := User{
			Email:    email,
			Role:     activeUser.Role,
			Products: activeUser.Products,
		}
		return &user, nil
	}

	log.Printf("user %v is not found in active user list, will look up in invitations", email)
	invitationList, err := us.InvitationService.List()
	if err != nil {
		return nil, err
	}
	var targetInvitation *Invitation
	for _, invitation := range invitationList.Organization.Invitations {
		if invitation.Email == email {
			targetInvitation = &invitation
			break
		}
	}
	return targetInvitation, nil
}
