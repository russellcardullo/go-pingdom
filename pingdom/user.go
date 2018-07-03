package pingdom

import (
	"io/ioutil"
	"encoding/json"
)

type UserService struct {
	client *Client
}


type UserApi interface {
	PutParams() map[string]string
	PutContactParams() map[string]string
	PostParams() map[string]string
	PostContactParams() map[string]string
	DeleteParams() map[string]string
	DeleteContactParams() map[string]string
	Valid() error
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

//func (cs *UserService) Read(id int) (*UsersResponse, error) {
//
//}
//
//func (cs *UserService) Create(user UserApi) (*UsersResponse, error) {
//
//}
//
//func (cs *UserService) Update(id int, user UserApi) (*PingdomResponse, error) {
//
//}
//
//func (cs *UserService) Delete(id int) (*PingdomResponse, error) {
//
//}
//
//func (cs *UserService) CreateContact(user UserApi) (*CreateUserContactResponse, error) {
//
//}
//
//func (cs *UserService) UpdateContact(id int, user UserApi) (*PingdomResponse, error) {
//
//}
//
//func (cs *UserService) DeleteContact(id int) (*PingdomResponse, error) {
//
//}