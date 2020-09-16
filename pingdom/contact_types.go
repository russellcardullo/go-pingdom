package pingdom

import (
	"encoding/json"
	"fmt"
)

// UserSms represents the sms contact object for a User.
type UserSms struct {
	Severity    string `json:"severity"`
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Provider    string `json:"provider"`
}

// UserEmail represents the email contact object for a User.
type UserEmail struct {
	Severity string `json:"severity"`
	Address  string `json:"address"`
}

// NotificationTargets represents different ways a contact could be notified of alerts
type NotificationTargets struct {
	SMS   []SMSNotification   `json:"sms,omitempty"`
	Email []EmailNotification `json:"email,omitempty"`
	APNS  []APNSNotification  `json:"apns,omitempty"`
	AGCM  []AGCMNotification  `json:"agcm,omitempty"`
}

// SMSNotification represents a text message notification
type SMSNotification struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Provider    string `json:"provider"`
	Severity    string `json:"severity"`
}

// EmailNotification represents an email address notification
type EmailNotification struct {
	Address  string `json:"address"`
	Severity string `json:"severity"`
}

// APNSNotification represents an APNS device notification
type APNSNotification struct {
	Device   string `json:"apns_device"`
	Name     string `json:"device_name"`
	Severity string `json:"severity"`
}

// AGCMNotification represents an AGCM notification
type AGCMNotification struct {
	AGCMID   string `json:"agcm_id"`
	Severity string `json:"severity"`
}

// ContactTeam represents an alerting team from the view of a Contact
type ContactTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Contact represents a Pingdom Contact.
type Contact struct {
	ID                  int                 `json:"id"`
	Name                string              `json:"name"`
	NotificationTargets NotificationTargets `json:"notification_targets"`
	Owner               bool                `json:"owner"`
	Paused              bool                `json:"paused"`
	Teams               []ContactTeam       `json:"teams"`
	Type                string              `json:"type"`
}

// ValidContact determines whether a Contact contains valid fields.
func (c *Contact) ValidContact() error {
	if c.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	return nil
}

// RenderForJSONAPI returns the JSON formatted version of this object that may be submitted to Pingdom
func (c *Contact) RenderForJSONAPI() string {
	u := map[string]interface{}{
		"name":                 c.Name,
		"notification_targets": c.NotificationTargets,
		"paused":               c.Paused,
	}
	jsonBody, _ := json.Marshal(u)
	return string(jsonBody)
}
