package solarwinds

const (
	listActiveUserOp           = "getUsersQuery"
	listActiveUserQuery        = "query getUsersQuery {\n  user {\n    id\n    currentOrganization {\n      id\n      members {\n        user {\n          id\n          firstName\n          lastName\n          email\n          lastLogin\n          __typename\n        }\n        role\n        products {\n          name\n          access\n          role\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
	listActiveUserResponseType = "user"

	getActiveUserOp           = "getEditUserQuery"
	getActiveUserQuery        = "query getEditUserQuery($userId: String!) {\n  user {\n    id\n    currentOrganization {\n      id\n      members(filter: {id: $userId}) {\n        id\n        user {\n          email\n          __typename\n        }\n        role\n        products {\n          name\n          role\n          access\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
	getActiveUserResponseType = "user"

	updateActiveUserOp           = "updateMemberRolesMutation"
	updateActiveUserQuery        = "mutation updateMemberRolesMutation($userId: ID!, $role: OrganizationRole!, $products: [ProductAccessInput!]) {\n  updateMemberRoles(userId: $userId, input: {role: $role, products: $products}) {\n    code\n    success\n    message\n    __typename\n  }\n}\n"
	updateActiveUserResponseType = "updateMemberRoles"
)

type UpdateActiveUserRequest struct {
	UserId   string    `json:"userId"`
	Role     string    `json:"role"`
	Products []Product `json:"products"`
	Email    string    `json:"-"`
}

type getActiveUserVars struct {
	UserId string `json:"userId"`
}

type ActiveUserList struct {
	OwnerUserId  string                  `json:"id"`
	Organization OrganizationWithMembers `json:"currentOrganization"`
}

type OrganizationWithMembers struct {
	Id      string               `json:"id"`
	Members []OrganizationMember `json:"members"`
}

type OrganizationMember struct {
	User     ActiveUser `json:"user"`
	Role     string     `json:"role"`
	Products []Product  `json:"products"`
}

type ActiveUser struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	LastLogin string `json:"lastLogin"`
}

type ActiveUserService struct {
	client *Client
}

func (us *ActiveUserService) List() (*ActiveUserList, error) {
	req := GraphQLRequest{
		OperationName: listActiveUserOp,
		Query:         listActiveUserQuery,
		ResponseType:  listActiveUserResponseType,
	}
	resp, err := us.client.MakeGraphQLRequest(&req)
	if err != nil {
		return nil, err
	}
	userList := ActiveUserList{}
	if err := Convert(&resp, &userList); err != nil {
		return nil, err
	}
	return &userList, nil
}

func (us *ActiveUserService) Get(userId string) (*ActiveUserList, error) {
	req := GraphQLRequest{
		OperationName: getActiveUserOp,
		Query:         getActiveUserQuery,
		Variables: getActiveUserVars{
			UserId: userId,
		},
		ResponseType: getActiveUserResponseType,
	}
	resp, err := us.client.MakeGraphQLRequest(&req)
	if err != nil {
		return nil, err
	}
	userList := ActiveUserList{}
	if err := Convert(&resp, &userList); err != nil {
		return nil, err
	}
	return &userList, nil
}

func (us *ActiveUserService) Update(update UpdateActiveUserRequest) error {
	req := GraphQLRequest{
		OperationName: updateActiveUserOp,
		Query:         updateActiveUserQuery,
		Variables:     update,
		ResponseType:  updateActiveUserResponseType,
	}
	_, err := us.client.MakeGraphQLRequest(&req)
	return err
}

func (us *ActiveUserService) GetByEmail(email string) (*OrganizationMember, error) {
	activeUserList, err := us.List()
	if err != nil {
		return nil, err
	}
	var targetUser *OrganizationMember
	for _, activeUser := range activeUserList.Organization.Members {
		if activeUser.User.Email == email {
			targetUser = &activeUser
			break
		}
	}
	return targetUser, nil
}
