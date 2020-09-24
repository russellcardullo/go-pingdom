package pingdom

import (
	"encoding/json"
	"fmt"
	"time"
)

// PingdomResponse represents a general response from the Pingdom API.
type PingdomResponse struct {
	Message string `json:"message"`
}

// PingdomError represents an error response from the Pingdom API.
type PingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

// CheckResponse represents the JSON response for a check from the Pingdom API.
type CheckResponse struct {
	ID                       int                 `json:"id"`
	Name                     string              `json:"name"`
	Resolution               int                 `json:"resolution,omitempty"`
	SendNotificationWhenDown int                 `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int                 `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool                `json:"notifywhenbackup,omitempty"`
	Created                  int64               `json:"created,omitempty"`
	Hostname                 string              `json:"hostname,omitempty"`
	Status                   string              `json:"status,omitempty"`
	LastErrorTime            int64               `json:"lasterrortime,omitempty"`
	LastTestTime             int64               `json:"lasttesttime,omitempty"`
	LastResponseTime         int64               `json:"lastresponsetime,omitempty"`
	Paused                   bool                `json:"paused,omitempty"`
	IntegrationIds           []int               `json:"integrationids,omitempty"`
	SeverityLevel            string              `json:"severity_level,omitempty"`
	Type                     CheckResponseType   `json:"type,omitempty"`
	Tags                     []CheckResponseTag  `json:"tags,omitempty"`
	UserIds                  []int               `json:"userids,omitempty"`
	Teams                    []CheckTeamResponse `json:"teams,omitempty"`
	ResponseTimeThreshold    int                 `json:"responsetime_threshold,omitempty"`
	ProbeFilters             []string            `json:"probe_filters,omitempty"`
	IP6                      bool                `json:"ip6,omitempty"`

	// Legacy; this is not returned by the API, we backfill the value from the
	// Teams field.
	TeamIds []int
}

// TmsCheckResponse represents the JSON response for a TMS check from the Pingdom API.
type TmsCheckResponse struct {
	ID                       int                    `json:"id"`
	Name                     string                 `json:"name"`
	Steps                    []TmsStep              `json:"steps"`
	Active                   bool                   `json:"active,omitempty"`
	ContactIds               []int                  `json:"contact_ids,omitempty"`
	CustomMessage            string                 `json:"custom_message,omitempty"`
	IntegrationIds           []int                  `json:"integration_ids,omitempty"`
	Interval                 int                    `json:"interval,omitempty"`
	Metadata                 map[string]interface{} `json:"metadata,omitempty"`
	Region                   string                 `json:"region,omitempty"`
	SendNotificationWhenDown int                    `json:"send_notification_when_down,omitempty"`
	SeverityLevel            string                 `json:"severity_level,omitempty"`
	Tags                     []string               `json:"tags,omitempty"`
	TeamIds                  []int                  `json:"team_ids,omitempty"`
}

// TmsStatusChangeResponse represents the JSON response for a TMS status change from the Pingdom API.
type TmsStatusChangeResponse struct {
	Report TmsStatusChange `json:"report"`
}

// TmsPerformanceReportResponse represents the JSON response for a TMS performance report for a single transaction from the Pingdom API.
type TmsPerformanceReportResponse struct {
	Report TmsPerformanceReport `json:"report"`
}

// TmsStatusChange represents the JSON for a TMS status change report from the Pingdom API.
type TmsStatusChange struct {
	CheckId int        `json:"check_id"`
	Name    string     `json:"name"`
	States  []TmsState `json:"states"`
}

type TmsPerformanceReport struct {
	CheckId    int           `json:"check_id"`
	Name       string        `json:"name"`
	Intervals  []TmsInterval `json:"intervals"`
	Resolution Resolution    `json:"resolution"`
}

type TmsState struct {
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Error   int       `json:"id,omitempty"`
}

type TmsInterval struct {
	AverageResponse int             `json:"average_response,omitempty"`
	Downtime        int             `json:"downtime,omitempty"`
	From            time.Time       `json:"from"`
	Steps           []TmsStepStatus `json:"steps,omitempty"`
	Unmonitored     int             `json:"unmonitored,omitempty"`
	Uptime          int             `json:"uptime,omitempty"`
}

type TmsStepStatus struct {
	Step            TmsStep `json:"step,omitempty"`
	AverageResponse int     `json:"average_response,omitempty"`
}

// CheckTeamResponse is a Team returned inside of a Check instance. (We can't
// use TeamResponse because the ID returned here is an int, not a string).
type CheckTeamResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CheckResponseType is the type of the Pingdom check.
type CheckResponseType struct {
	Name string                    `json:"-"`
	HTTP *CheckResponseHTTPDetails `json:"http,omitempty"`
	TCP  *CheckResponseTCPDetails  `json:"tcp,omitempty"`
}

// CheckResponseTag is an optional tag that can be added to checks.
type CheckResponseTag struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Count interface{} `json:"count"`
}

// MaintenanceResponse represents the JSON response for a maintenance from the Pingdom API.
type MaintenanceResponse struct {
	ID             int                      `json:"id"`
	Description    string                   `json:"description"`
	From           int64                    `json:"from"`
	To             int64                    `json:"to"`
	RecurrenceType string                   `json:"recurrencetype"`
	RepeatEvery    int                      `json:"repeatevery"`
	EffectiveTo    int64                    `json:"effectiveto"`
	Checks         MaintenanceCheckResponse `json:"checks"`
}

// MaintenanceCheckResponse represents Check reply in json MaintenanceResponse.
type MaintenanceCheckResponse struct {
	Uptime []int `json:"uptime"`
	Tms    []int `json:"tms"`
}

// ProbeResponse represents the JSON response for probes from the Pingdom API.
type ProbeResponse struct {
	ID         int    `json:"id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	IPv6       string `json:"ipv6"`
	CountryISO string `json:"countryiso"`
	Region     string `json:"region"`
}

// SummaryPerformanceResponse represents the JSON response for a summary performance from the Pingdom API.
type SummaryPerformanceResponse struct {
	Summary SummaryPerformanceMap `json:"summary"`
}

// SummaryPerformanceMap is the performance broken down over different time intervals.
type SummaryPerformanceMap struct {
	Hours []SummaryPerformanceSummary `json:"hours,omitempty"`
	Days  []SummaryPerformanceSummary `json:"days,omitempty"`
	Weeks []SummaryPerformanceSummary `json:"weeks,omitempty"`
}

// SummaryPerformanceSummary is the metrics for a performance summary.
type SummaryPerformanceSummary struct {
	AvgResponse int `json:"avgresponse"`
	Downtime    int `json:"downtime"`
	StartTime   int `json:"starttime"`
	Unmonitored int `json:"unmonitored"`
	Uptime      int `json:"uptime"`
}

// ResultsResponse represents the JSON response for detailed check results from the Pingdom API.
type ResultsResponse struct {
	ActiveProbes []int    `json:"activeprobes"`
	Results      []Result `json:"results"`
}

// Result reprensents the JSON response for a detailed check result.
type Result struct {
	ProbeID        int    `json:"probeid"`
	Time           int    `json:"time"`
	Status         string `json:"status"`
	ResponseTime   int    `json:"responsetime"`
	StatusDesc     string `json:"statusdesc"`
	StatusDescLong string `json:"statusdesclong"`
}

// UserSmsResponse represents the JSON response for a user SMS contact.
type UserSmsResponse struct {
	Id          int    `json:"id"`
	Severity    string `json:"severity"`
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Provider    string `json:"provider"`
}

// UserEmailResponse represents the JSON response for a user email contact.
type UserEmailResponse struct {
	Id       int    `json:"id"`
	Severity string `json:"severity"`
	Address  string `json:"address"`
}

// CreateUserContactResponse represents the JSON response for a user contact.
type CreateUserContactResponse struct {
	Id int `json:"id"`
}

// UsersResponse represents the JSON response for a Pingom User.
type UsersResponse struct {
	Id       int                 `json:"id"`
	Paused   string              `json:"paused,omitempty"`
	Username string              `json:"name,omitempty"`
	Sms      []UserSmsResponse   `json:"sms,omitempty"`
	Email    []UserEmailResponse `json:"email,omitempty"`
}

// UnmarshalJSON converts a byte array into a CheckResponseType.
func (c *CheckResponseType) UnmarshalJSON(b []byte) error {
	var raw interface{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		c.Name = v
	case map[string]interface{}:
		if len(v) != 1 {
			return fmt.Errorf("Check detailed response `check.type` contains more than one object: %+v", v)
		}
		for k := range v {
			c.Name = k
		}

		// Allow continue use json.Unmarshall using a type != Unmarshaller
		// This avoid enter in a infinite loop
		type t CheckResponseType
		var rawCheckDetails t

		err := json.Unmarshal(b, &rawCheckDetails)
		if err != nil {
			return err
		}
		c.HTTP = rawCheckDetails.HTTP
		c.TCP = rawCheckDetails.TCP
	}
	return nil
}

// CheckResponseHTTPDetails represents the details specific to HTTP checks.
type CheckResponseHTTPDetails struct {
	Url               string            `json:"url,omitempty"`
	Encryption        bool              `json:"encryption,omitempty"`
	Port              int               `json:"port,omitempty"`
	Username          string            `json:"username,omitempty"`
	Password          string            `json:"password,omitempty"`
	ShouldContain     string            `json:"shouldcontain,omitempty"`
	ShouldNotContain  string            `json:"shouldnotcontain,omitempty"`
	PostData          string            `json:"postdata,omitempty"`
	RequestHeaders    map[string]string `json:"requestheaders,omitempty"`
	VerifyCertificate bool              `json:"verify_certificate,omitempty"`
	SSLDownDaysBefore int               `json:"ssl_down_days_before,omitempty"`
}

// CheckResponseTCPDetails represents the details specific to TCP checks.
type CheckResponseTCPDetails struct {
	Port           int    `json:"port,omitempty"`
	StringToSend   string `json:"stringtosend,omitempty"`
	StringToExpect string `json:"stringtoexpect,omitempty"`
}

// Return string representation of the PingdomError.
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall JSON responses from Pingdom.

type listChecksJSONResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type listTmsChecksJSONResponse struct {
	TmsChecks []TmsCheckResponse `json:"checks"`
	Limit     int                `json:"limit,omitempty"`
	Offset    int                `json:"offset,omitempty"`
}

type listMaintenanceJSONResponse struct {
	Maintenances []MaintenanceResponse `json:"maintenance"`
}

type listProbesJSONResponse struct {
	Probes []ProbeResponse `json:"probes"`
}

type checkDetailsJSONResponse struct {
	Check *CheckResponse `json:"check"`
}

type tmsCheckDetailsJSONResponse struct {
	TmsCheck *TmsCheckResponse `json:"check"`
}

type maintenanceDetailsJSONResponse struct {
	Maintenance *MaintenanceResponse `json:"maintenance"`
}

type createUserContactJSONResponse struct {
	Contact *CreateUserContactResponse `json:"contact_target"`
}

type createUserJSONResponse struct {
	User *UsersResponse `json:"user"`
}

type listUsersJSONResponse struct {
	Users []UsersResponse `json:"users"`
}

type errorJSONResponse struct {
	Error *PingdomError `json:"error"`
}
