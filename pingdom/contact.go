package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// ContactService provides an interface to Pingdom contacts.
type ContactService struct {
	client *Client
}

// ContactAPI is an interface representing a Pingdom Contact.
type ContactAPI interface {
	RenderForJSONAPI() string
	ValidContact() error
}

// List returns a list of all contacts and their contact details.
func (cs *ContactService) List() ([]Contact, error) {

	req, err := cs.client.NewRequest("GET", "/alerting/contacts", nil)
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

	u := &listContactsJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &u)

	return u.Contacts, err
}

// Read return a contact object from Pingdom.
func (cs *ContactService) Read(contactID int) (*Contact, error) {
	req, err := cs.client.NewRequest("GET", "/alerting/contacts/"+strconv.Itoa(contactID), nil)
	if err != nil {
		return nil, err
	}

	c := &contactDetailsJSONResponse{}
	_, err = cs.client.Do(req, c)
	if err != nil {
		return nil, err
	}

	return c.Contact, nil
}

// Create adds a new contact.
func (cs *ContactService) Create(contact ContactAPI) (*Contact, error) {
	if err := contact.ValidContact(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJSONRequest("POST", "/alerting/contacts", contact.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	m := &createContactJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return m.Contact, err
}

// Update a contact's core properties not contact targets.
func (cs *ContactService) Update(id int, contact ContactAPI) (*PingdomResponse, error) {
	if err := contact.ValidContact(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJSONRequest("PUT", "/alerting/contacts/"+strconv.Itoa(id), contact.RenderForJSONAPI())
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

// Delete removes a contact from Pingdom.
func (cs *ContactService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/alerting/contacts/"+strconv.Itoa(id), nil)
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
