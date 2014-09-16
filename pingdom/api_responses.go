package pingdom

import "fmt"

// PingdomResponse represents a general response from the Pingdom API
type PingdomResponse struct {
	Message string `json:"message"`
}

// private types used to unmarshall json responses from pingdom

type checkResponse struct {
	Check *Check `json:"check"`
}

type listChecksResponse struct {
	Checks []Check `json:"checks"`
}

type pingdomErrorResponse struct {
	Error *pingdomError `json:"error"`
}

type pingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

func (r *pingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}
