package pingdom

import (
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type UserService struct {
	client *Client
}

type ContactApi interface {
	ValidCreateContact() error
	PostContactParams() map[string]string
}

type UserApi interface {
	ValidCreate() error
	PostParams() map[string]string
	//PutParams() map[string]string
	//PutContactParams() map[string]string
}


func (cs *UserService) List(params ...map[string]string) ([]UsersResponse, error) {
	param := map[string]string{}
	if len(params) != 0 {
		for _, m := range params {
			for k, v := range m {
				param[k] = v
			}
		}
	}
	req, err := cs.client.NewRequest("GET", "/users", param)
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

func (cs *UserService) Create(user UserApi) (*UsersResponse, error) {
	if err := user.ValidCreate(); err != nil {
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

func (cs *UserService) CreateContact(userId int, contact Contact) (*CreateUserContactResponse, error) {
	if err := contact.ValidCreateContact(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/users/"+ strconv.Itoa(userId), contact.PostContactParams())
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

//func (cs *UserService) Update(id int, user UserApi) (*PingdomResponse, error) {
//
//}
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

//func (cs *UserService) UpdateContact(id int, user UserApi) (*PingdomResponse, error) {
//
//}
//
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