package pingdom

import (
	"encoding/json"
	"fmt"
)

// Team represents a Pingdom Team Data.
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	MemberIDs []int  `json:"member_ids,omitempty"`
}

// RenderForJSONAPI returns the JSON formatted version of this object that may be submitted to Pingdom
func (t *Team) RenderForJSONAPI() string {
	b := map[string]interface{}{
		"name":       t.Name,
		"member_ids": t.MemberIDs,
	}
	jsonBody, _ := json.Marshal(b)
	return string(jsonBody)
}

// Valid Determines whether the Team contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API.
func (t *Team) Valid() error {
	if t.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	return nil
}
