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

// CheckResponse represents the json response for a check from the Pingdom API
type CheckResponse struct {
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

// Return string representation of the PingdomError
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall json responses from pingdom

type listChecksJsonResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type checkDetailsJsonResponse struct {
	Check *CheckResponse `json:"check"`
}

type errorJsonResponse struct {
	Error *PingdomError `json:"error"`
}
