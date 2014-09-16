package pingdom

import "fmt"

// PingdomResponse represents a general response from the Pingdom API
type PingdomResponse struct {
	Message string `json:"message"`
}

// PingdomError represents an error response from the Pingdom API
type PingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

// Return string representation of the PingdomError
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall json responses from pingdom

type checkJsonResponse struct {
	Check *Check `json:"check"`
}

type listChecksJsonResponse struct {
	Checks []Check `json:"checks"`
}

type errorJsonResponse struct {
	Error *PingdomError `json:"error"`
}
