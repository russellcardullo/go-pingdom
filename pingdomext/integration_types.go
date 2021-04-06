package pingdomext

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// IntegrationProvider represents a Pingdom integration provider.
type IntegrationProvider struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WebHookIntegration represents a Pingdom WebHook integration.
type WebHookIntegration struct {
	Active     bool         `json:"active"`
	ProviderID int          `json:"provider_id"`
	UserData   *WebHookData `json:"user_data"`
}

// WebHookData represents a WebHook data in the WebHook integration.
type WebHookData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

/*
// LibratoIntegration represents a Pingdom Librato integration.
type LibratoIntegration struct {
	Active     bool         `json:"active"`
	ProviderID int          `json:"provider_id"`
	UserData   *LibratoData `json:"user_data"`
}

// LibratoData represents a Librato data in the Librato integration.
type LibratoData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ApiToken string `json:"apiToken"`
}
*/

// PostParams returns a map of parameters for an WebHook integration that can be sent along.
func (wi *WebHookIntegration) PostParams() map[string]string {
	dataJSON, err := json.Marshal(wi.UserData)
	fmt.Println(err)
	m := map[string]string{
		"active":      strconv.FormatBool(wi.Active),
		"provider_id": strconv.Itoa(wi.ProviderID),
		"data_json":   string(dataJSON),
	}
	return m
}

// Valid determines whether the WebHook integration contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API.
func (wi *WebHookIntegration) Valid() error {
	if wi.ProviderID != 1 && wi.ProviderID != 2 {
		return fmt.Errorf("Invalid value for `provider`.  Must contain available provider id")
	}
	if wi.UserData.Name == "" {
		return fmt.Errorf("Invalid value for `name`.  Must contain non-empty string")
	}
	if wi.UserData.URL == "" {
		return fmt.Errorf("Invalid value for `url`.  Must contain non-empty string")
	}
	return nil
}

/*

// PostParams returns a map of parameters for an Librato integration that can be sent along.
func (li *LibratoIntegration) PostParams() map[string]string {
	dataJSON, err := toJsonNoEscape(li.UserData)
	fmt.Println(err)
	m := map[string]string{
		"active":      strconv.FormatBool(li.Active),
		"provider_id": strconv.Itoa(li.ProviderID),
		"data_json":   string(dataJSON),
	}
	return m
}

// Valid determines whether the Librato integration contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API.
func (li *LibratoIntegration) Valid() error {
	if li.ProviderID != 1 && li.ProviderID != 2 {
		return fmt.Errorf("Invalid value for `provider`.  Must contain available provider")
	}
	if li.UserData.Name == "" {
		return fmt.Errorf("Invalid value for `name`.  Must contain non-empty string")
	}
	if li.UserData.ApiToken == "" {
		return fmt.Errorf("Invalid value for `api token`.  Must contain non-empty string")
	}
	if li.UserData.Email == "" {
		return fmt.Errorf("Invalid value for `email`.  Must contain non-empty string")
	}
	return nil
}

*/
