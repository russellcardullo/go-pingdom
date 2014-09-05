package pingdom

import (
	"errors"
	"fmt"
	"strconv"
)

// Check represents a Pingdom Check
type Check struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	Resolution               int    `json:"resolution,omitempty"`
	SendToAndroid            bool   `json:"sendtoandroid,omitempty"`
	SendToEmail              bool   `json:"sendtoemail,omitempty"`
	SendToIPhone             bool   `json:"sendtoiphone,omitempty"`
	SendToSms                bool   `json:"sendtosms,omitempty"`
	SendToTwitter            bool   `json:"sendtotwitter,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
	Created                  int64  `json:"created,omitempty"`
	Hostname                 string `json:"hostname,omitempty"`
	Status                   string `json:"status,omitempty"`
	LastErrorTime            int64  `json:"lasterrortime,omitempty"`
	LastTestTime             int64  `json:"lasttesttime,omitempty"`
	LastResponseTime         int64  `json:"lastresponsetime,omitempty"`
	Paused                   bool   `json:"paused,omitempty"`
	ContactIds               []int  `json:"contactids,omitempty"`
}

// Params returns a map of parameters for a Check that can be sent along
// with an HTTP POST or PUT request
func (ck *Check) Params() map[string]string {
	return map[string]string{
		"name":                     ck.Name,
		"host":                     ck.Hostname,
		"paused":                   strconv.FormatBool(ck.Paused),
		"resolution":               strconv.Itoa(ck.Resolution),
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
func (ck *Check) Valid() error {
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
