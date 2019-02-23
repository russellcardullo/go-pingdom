package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// UserService provides an interface to Pingdom users.
type UserService struct {
	client *Client
}

// User is an interface representing a Pingdom team.
type UserApi interface {
	ValidUser() error
	PostParams() map[string]string
	PutParams() map[string]string
}

// List returns a list of all users and their contact details.
func (cs *UserService) List() ([]UsersResponse, error) {

	req, err := cs.client.NewRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	u := &listUsersJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &u)

	return u.Users, err
}

// Read return a user object from Pingdom.
func (cs *UserService) Read(userId int) (*UsersResponse, error) {
	users, err := cs.List()
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].Id == userId {
			return &users[i], nil
		}
	}

	return nil, fmt.Errorf("UserId: " + strconv.Itoa(userId) + " not found")
}

// Create adds a new user.
func (cs *UserService) Create(user UserApi) (*UsersResponse, error) {
	if err := user.ValidUser(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/users", user.PostParams())
	if err != nil {
		return nil, err
	}

	m := &createUserJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.User, err
}

// CreateContact adds a contact target to an existing user.
func (cs *UserService) CreateContact(userId int, contact Contact) (*CreateUserContactResponse, error) {
	if err := contact.ValidContact(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/users/"+strconv.Itoa(userId), contact.PostContactParams())
	if err != nil {
		return nil, err
	}

	m := &createUserContactJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Contact, err
}

// Update a user's core properties not contact targets.
func (cs *UserService) Update(id int, user UserApi) (*PingdomResponse, error) {
	if err := user.ValidUser(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("PUT", "/users/"+strconv.Itoa(id), user.PutParams())
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// UpdateContact updates a contact by id, will change an email to sms or sms to email
// if you provide an id for the other.
func (cs *UserService) UpdateContact(userId int, contactId int, contact Contact) (*PingdomResponse, error) {
	if err := contact.ValidContact(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("PUT", "/users/"+strconv.Itoa(userId)+"/"+strconv.Itoa(contactId), contact.PutContactParams())
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// Delete removes a user from Pingdom.
func (cs *UserService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/users/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// DeleteContact deletes a contact target from a user, either an email or sms property of a user.
func (cs *UserService) DeleteContact(userId int, contactId int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/users/"+strconv.Itoa(userId)+"/"+strconv.Itoa(contactId), nil)
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}
