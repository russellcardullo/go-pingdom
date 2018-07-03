package pingdom

import (
	"fmt"
)

type UserSms struct {
	Severity string `json:"severity"`
	CountryCode string `json:"country_code"`
	Number string `json:"number"`
	Provider string `json:"provider"`
}

type UserEmail struct {
	Severity string `json:"severity"`
	Address string `json:"address"`
}

type Contact struct {
	Severity string `json:"severitylevel"`
	CountryCode string `json:"countrycode"`
	Number string `json:"number"`
	Provider string `json:"provider"`
	Email string `json:"email"`
}

// MaintenanceWindow represents a Pingdom Maintenance Window.
type User struct {
	Paused         int64  `json:"paused,omitempty"`
	Username       string `json:"name,omitempty"`
	Sms			   []UserSmsResponse `json:"sms,omitempty"`
	Email 		   []UserEmailResponse `json:"email,omitempty"`
}

func (u *User) ValidCreate() error {
	if u.Username == "" {
		return fmt.Errorf("Invalid value for `Username`.  Must contain non-empty string")
	}

	return nil
}

func (c *Contact) ValidCreateContact() error {
	if c.Email == "" && c.Number == "" {
		return fmt.Errorf("you must provide either an Email or a Phone Number to create a contact target")
}

	if c.Number != "" && c.CountryCode == "" {
		return fmt.Errorf("you must provide a Country Code if providing a phone number")
	}

	return nil
}

func (u *User) PostParams() map[string]string {
	m := map[string]string{
		"name": u.Username,
	}

	return m
}

func (c *Contact) PostContactParams() map[string]string {
	m := map[string]string{}

	// Ignore if not defined
	if c.Email != "" {
		m["email"] = c.Email
	}

	if c.Number != "" {
		m["number"] = c.Number
	}

	if c.CountryCode != "" {
		m["countrycode"] = c.CountryCode
	}

	if c.Severity != "" {
		m["severitylevel"] = c.Severity
	}

	if c.Provider != "" {
		m["provider"] = c.Provider
	}

	return m
}

//func (u *User) PutParams() map[string]string {
//
//}
//
//func (u *User) PutContactParams() map[string]string {
//
//}
//
//func (u *User) DeleteParams() map[string]string {
//
//}
//
//func (u *User) DeleteContactParams() map[string]string {
//
//}
//
