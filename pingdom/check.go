package pingdom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

// CheckService provides an interface to Pingdom checks
type CheckService struct {
	client *Client
}

type Check interface {
	Params() map[string]string
	Valid() error
}

// Check represents a Pingdom http check.
type HttpCheck struct {
	Name                     string `json:"name"`
	Hostname                 string `json:"hostname,omitempty"`
	Resolution               int    `json:"resolution,omitempty"`
	Paused                   bool   `json:"paused,omitempty"`
	SendToAndroid            bool   `json:"sendtoandroid,omitempty"`
	SendToEmail              bool   `json:"sendtoemail,omitempty"`
	SendToIPhone             bool   `json:"sendtoiphone,omitempty"`
	SendToSms                bool   `json:"sendtosms,omitempty"`
	SendToTwitter            bool   `json:"sendtotwitter,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
}

// Return a list of checks from Pingdom.
func (cs *CheckService) List() ([]CheckResponse, error) {
	req, err := cs.client.NewRequest("GET", "/api/2.0/checks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listChecksJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Checks, err
}

// Create a new check.  This function will validate the given check param
// to ensure that it contains correct values before submitting the request
// Returns a Check object representing the response from Pingdom.  Note
// that Pingdom does not return a full check object so in the returned
// object you should only use the ID field.
func (cs *CheckService) Create(check Check) (*CheckResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/api/2.0/checks", check.Params())
	if err != nil {
		return nil, err
	}

	m := &checkDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Check, err
}

// ReadCheck returns detailed information about a pingdom check given its ID.
func (cs *CheckService) Read(id int) (*CheckResponse, error) {
	req, err := cs.client.NewRequest("GET", "/api/2.0/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &checkDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m.Check, err
}

// UpdateCheck will update the check represented by the given ID with the values
// in the given check.  You should submit the complete list of values in
// the given check parameter, not just those that have changed.
func (cs *CheckService) Update(id int, check Check) (*PingdomResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	params := check.Params()
	delete(params, "type")
	req, err := cs.client.NewRequest("PUT", "/api/2.0/checks/"+strconv.Itoa(id), params)
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

// DeleteCheck will delete the check for the given ID.
func (cs *CheckService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/api/2.0/checks/"+strconv.Itoa(id), nil)
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

// Params returns a map of parameters for a Check that can be sent along
// with an HTTP POST or PUT request
func (ck *HttpCheck) Params() map[string]string {
	return map[string]string{
		"name":                     ck.Name,
		"host":                     ck.Hostname,
		"resolution":               strconv.Itoa(ck.Resolution),
		"paused":                   strconv.FormatBool(ck.Paused),
		"sendtoemail":              strconv.FormatBool(ck.SendToEmail),
		"sendtosms":                strconv.FormatBool(ck.SendToSms),
		"sendtotwitter":            strconv.FormatBool(ck.SendToTwitter),
		"sendtoiphone":             strconv.FormatBool(ck.SendToIPhone),
		"sendtoandroid":            strconv.FormatBool(ck.SendToAndroid),
		"sendnotificationwhendown": strconv.Itoa(ck.SendNotificationWhenDown),
		"notifyagainevery":         strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup":         strconv.FormatBool(ck.NotifyWhenBackup),
		"type":                     "http",
	}
}

// Determine whether the Check contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *HttpCheck) Valid() error {
	if ck.Name == "" {
		return errors.New("Invalid value for `Name`.  Must contain non-empty string")
	}

	if ck.Hostname == "" {
		return errors.New("Invalid value for `Hostname`.  Must contain non-empty string")
	}

	if ck.Resolution != 1 && ck.Resolution != 5 && ck.Resolution != 15 &&
		ck.Resolution != 30 && ck.Resolution != 60 {
		err := fmt.Sprintf("Invalid value %v for `Resolution`.  Allowed values are [1,5,15,30,60].", ck.Resolution)
		return errors.New(err)
	}
	return nil
}
