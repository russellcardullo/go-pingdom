package pingdom

import (
	"errors"
	"fmt"
	"strconv"
)

// HttpCheck represents a Pingdom http check.
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
	UseLegacyNotifications   bool   `json:"use_legacy_notifications,omitempty"`
}

// PingCheck represents a Pingdom ping check
type PingCheck struct {
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
	UseLegacyNotifications   bool   `json:"use_legacy_notifications,omitempty"`
}

// Params returns a map of parameters for an HttpCheck that can be sent along
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
		"use_legacy_notifications": strconv.FormatBool(ck.UseLegacyNotifications),
		"type": "http",
	}
}

// Determine whether the HttpCheck contains valid fields.  This can be
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

// Params returns a map of parameters for a PingCheck that can be sent along
// with an HTTP POST or PUT request
func (ck *PingCheck) Params() map[string]string {
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
		"use_legacy_notifications": strconv.FormatBool(ck.UseLegacyNotifications),
		"type": "ping",
	}
}

// Determine whether the PingCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *PingCheck) Valid() error {
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
