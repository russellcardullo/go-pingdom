package solarwinds

// Constant values used in GraphQL requests.
const (
	inviteUserOp           = "createOrganizationAdminMutation"
	inviteUserQuery        = "mutation createOrganizationAdminMutation($input: CreateOrganizationInvitationInput!) {\n  createOrganizationInvitation(input: $input) {\n    success\n    code\n    message\n    invitation {\n      email\n      role\n      __typename\n    }\n    __typename\n  }\n}\n"
	inviteUserResponseType = "createOrganizationInvitation"

	revokeInvitationOp           = "deleteOrganizationInvitationMutation"
	revokeInvitationQuery        = "mutation deleteOrganizationInvitationMutation($email: ID!) {\n  deleteOrganizationInvitation(email: $email) {\n    success\n    code\n    message\n    __typename\n  }\n}\n"
	revokeInvitationResponseType = "deleteOrganizationInvitation"

	resendInvitationOp           = "resendOrganizationInvitationMutation"
	resendInvitationQuery        = "mutation resendOrganizationInvitationMutation($email: ID!) {\n  resendOrganizationInvitation(email: $email) {\n    success\n    code\n    message\n    __typename\n  }\n}\n"
	resendInvitationResponseType = "resendOrganizationInvitation"

	listInvitationOp           = "getInvitationsQuery"
	listInvitationQuery        = "query getInvitationsQuery {\n  user {\n    id\n    currentOrganization {\n      id\n      invitations {\n        email\n        role\n        date\n        products {\n          name\n          role\n          access\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
	listInvitationResponseType = "user"
)

type Invitation struct {
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Products []Product `json:"products"`
}

type Product struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type InvitationList struct {
	OwnerUserId  string                      `json:"id"`
	Organization OrganizationWithInvitations `json:"currentOrganization"`
}

type OrganizationWithInvitations struct {
	Id          string       `json:"id"`
	Invitations []Invitation `json:"invitations"`
}

type InvitationService struct {
	client *Client
}

type inviteUserVars struct {
	Input Invitation `json:"input"`
}

type revokeInvitationVars struct {
	Email string `json:"email"`
}

type resendInvitationVars struct {
	Email string `json:"email"`
}

func (is *InvitationService) Create(user Invitation) error {
	req := GraphQLRequest{
		OperationName: inviteUserOp,
		Query:         inviteUserQuery,
		Variables: inviteUserVars{
			Input: user,
		},
		ResponseType: inviteUserResponseType,
	}
	_, err := is.client.MakeGraphQLRequest(&req)
	return err
}

func (is *InvitationService) Revoke(email string) error {
	req := GraphQLRequest{
		OperationName: revokeInvitationOp,
		Query:         revokeInvitationQuery,
		Variables: revokeInvitationVars{
			Email: email,
		},
		ResponseType: revokeInvitationResponseType,
	}
	_, err := is.client.MakeGraphQLRequest(&req)
	return err
}

func (is *InvitationService) Resend(email string) error {
	req := GraphQLRequest{
		OperationName: resendInvitationOp,
		Query:         resendInvitationQuery,
		Variables: resendInvitationVars{
			Email: email,
		},
		ResponseType: resendInvitationResponseType,
	}
	_, err := is.client.MakeGraphQLRequest(&req)
	return err
}

func (is *InvitationService) List() (*InvitationList, error) {
	req := GraphQLRequest{
		OperationName: listInvitationOp,
		Query:         listInvitationQuery,
		ResponseType:  listInvitationResponseType,
	}
	resp, err := is.client.MakeGraphQLRequest(&req)
	if err != nil {
		return nil, err
	}
	invitationList := InvitationList{}
	if err := Convert(&resp, &invitationList); err != nil {
		return nil, err
	}
	return &invitationList, nil
}
