package pingdom

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// HttpCheck represents a Pingdom http check.
type HttpCheck struct {
	Name                     string            `json:"name"`
	Hostname                 string            `json:"hostname,omitempty"`
	Resolution               int               `json:"resolution,omitempty"`
	Paused                   bool              `json:"paused,omitempty"`
	SendToAndroid            bool              `json:"sendtoandroid,omitempty"`
	SendToEmail              bool              `json:"sendtoemail,omitempty"`
	SendToIPhone             bool              `json:"sendtoiphone,omitempty"`
	SendToSms                bool              `json:"sendtosms,omitempty"`
	SendToTwitter            bool              `json:"sendtotwitter,omitempty"`
	SendNotificationWhenDown int               `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int               `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool              `json:"notifywhenbackup,omitempty"`
	UseLegacyNotifications   bool              `json:"use_legacy_notifications,omitempty"`
	Url                      string            `json:"url,omitempty"`
	Encryption               bool              `json:"encryption,omitempty"`
	Port                     int               `json:"port,omitempty"`
	Username                 string            `json:"username,omitempty"`
	Password                 string            `json:"password,omitempty"`
	ShouldContain            string            `json:"shouldcontain,omitempty"`
	ShouldNotContain         string            `json:"shouldnotcontain,omitempty"`
	PostData                 string            `json:"postdata,omitempty"`
	RequestHeaders           map[string]string `json:"requestheaders,omitempty"`
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
	m := map[string]string{
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
		"url":        ck.Url,
		"encryption": strconv.FormatBool(ck.Encryption),
		"postdata":   ck.PostData,
		"type":       "http",
	}

	// Ignore port is not defined
	if ck.Port != 0 {
		m["port"] = strconv.Itoa(ck.Port)
	}

	// ShouldContain and ShouldNotContain are mutually exclusive.
	// But we must define one so they can be emptied if required.
	if ck.ShouldContain != "" {
		m["shouldcontain"] = ck.ShouldContain
	} else {
		m["shouldnotcontain"] = ck.ShouldNotContain
	}

	// Convert auth
	if ck.Username != "" {
		m["auth"] = fmt.Sprintf("%s:%s", ck.Username, ck.Password)
	}

	// Convert headers
	var headers []string
	for k := range ck.RequestHeaders {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for i, k := range headers {
		m[fmt.Sprintf("requestheader%d", i)] = fmt.Sprintf("%s:%s", k, ck.RequestHeaders[k])
	}

	return m
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

	if ck.ShouldContain != "" && ck.ShouldNotContain != "" {
		err := fmt.Sprintf("`ShouldContain` and `ShouldNotContain` must not be declared at the same time")
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
