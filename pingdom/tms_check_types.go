package pingdom

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type TMSCheck struct {
	Name                     string            `json:"name,omitempty"`
	Steps                    []TMSCheckStep    `json:"steps,omitempty"`
	Active                   bool              `json:"active"`
	ContactIDs               []int             `json:"contact_ids,omitempty"`
	CustomMessage            string            `json:"custom_message,omitempty"`
	IntegrationIDs           []int             `json:"integration_ids,omitempty"`
	Interval                 int64             `json:"interval,omitempty"`
	Metadata                 *TMSCheckMetaData `json:"metadata,omitempty"`
	Region                   string            `json:"region,omitempty"`
	SendNotificationWhenDown int               `json:"send_notification_when_down,omitempty"`
	SeverityLevel            string            `json:"severity_level,omitempty"`
	Tags                     []string          `json:"tags,omitempty"`
	TeamIDs                  []int             `json:"team_ids,omitempty"`
}

type TMSCheckStep struct {
	Args map[string]string `json:"args,omitempty"`
	Fn   string            `json:"fn,omitempty"`
}

type TMSCheckMetaData struct {
	Authentications    interface{} `json:"authentications,omitempty"`
	DisableWebSecurity bool        `json:"disableWebSecurity,omitempty"`
	Height             int         `json:"height,omitempty"`
	Width              int         `json:"width,omitempty"`
}

// RenderForJSONAPI returns the JSON formatted version of this object that may be submitted to Pingdom
func (t *TMSCheck) RenderForJSONAPI() string {
	jsonBody, _ := json.Marshal(t)
	return string(jsonBody)
}

// Valid Determines whether the TMSCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API.
func (t *TMSCheck) Valid() error {

	if t.Name == "" {
		return fmt.Errorf("Invalid value for `Name`. Must contain non-empty string.")
	}

	if t.Steps == nil {
		return fmt.Errorf("Invalid value for `Steps`. Must contain non-empty value.")
	}

	if len(t.Steps) == 0 {
		return fmt.Errorf("Invalid value for `Steps`. Must contain non-empty value.")
	}

	if t.Interval != 0 && t.Interval != 5 && t.Interval != 10 && t.Interval != 20 && t.Interval != 60 && t.Interval != 720 && t.Interval != 1440 {
		return fmt.Errorf("Invalid value for `Interval`. Please provide one of the following valid values instead: [5 10 20 60 720 1440].")
	}

	if t.SeverityLevel != "" && t.SeverityLevel != "high" && t.SeverityLevel != "low" {
		return fmt.Errorf("Invalid value for `SeverityLevel`. Please provide one of the following valid values instead: [high,low].")
	}

	if t.Tags != nil {
		for _, tag := range t.Tags {
			reg := regexp.MustCompile(`[0-9A-Za-z_-]+`)
			match := reg.FindString(tag)
			if match != tag {
				return fmt.Errorf("Invalid value for `Tags`. The tag name may contain the characters 'A-Z', 'a-z', '0-9', '_' and '-'.")

			}

		}
	}

	return nil
}
