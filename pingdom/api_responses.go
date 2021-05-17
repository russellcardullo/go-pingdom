package pingdom

import (
	"encoding/json"
	"fmt"
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
	IPv6                     bool                `json:"ipv6,omitempty"`

	// Legacy; this is not returned by the API, we backfill the value from the
	// Teams field.
	TeamIds []int
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
	DNS  *CheckResponseDNSDetails  `json:"dns,omitempty"`
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

// TeamResponse represents the JSON response for alerting teams from the Pingdom API.
type TeamResponse struct {
	ID      int                  `json:"id"`
	Name    string               `json:"name,omitempty"`
	Members []TeamMemberResponse `json:"members,omitempty"`
}

// TeamMemberResponse represents the JSON response for contacts in alerting teams from the Pingdom API.
type TeamMemberResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// TeamDeleteResponse represents the JSON response for delete team from the Pingdom API.
type TeamDeleteResponse struct {
	Message string `json:"message"`
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
		c.DNS = rawCheckDetails.DNS
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

// CheckResponseDNSDetails represents the details specific to DNS checks.
type CheckResponseDNSDetails struct {
	ExpectedIP string `json:"expectedip,omitempty"`
	NameServer string `json:"nameserver,omitempty"`
}

// Return string representation of the PingdomError.
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall JSON responses from Pingdom.

type listChecksJSONResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type listMaintenanceJSONResponse struct {
	Maintenances []MaintenanceResponse `json:"maintenance"`
}

type listProbesJSONResponse struct {
	Probes []ProbeResponse `json:"probes"`
}

type listTeamsJSONResponse struct {
	Teams []TeamResponse `json:"teams"`
}

type teamDetailsJSONResponse struct {
	Team *TeamResponse `json:"team"`
}

type contactDetailsJSONResponse struct {
	Contact *Contact `json:"contact"`
}

type checkDetailsJSONResponse struct {
	Check *CheckResponse `json:"check"`
}

type maintenanceDetailsJSONResponse struct {
	Maintenance *MaintenanceResponse `json:"maintenance"`
}

type createContactJSONResponse struct {
	Contact *Contact `json:"contact"`
}

type listContactsJSONResponse struct {
	Contacts []Contact `json:"contacts"`
}

// TMSCheckResponse represents the  JSON response for a TMS Check from the Pingdom API.
type TMSCheckResponse struct {
	ID                int      `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	Type              string   `json:"type,omitempty"`
	Active            bool     `json:"active,omitempty"`
	Status            string   `json:"status,omitempty"`
	Interval          int      `json:"interval,omitempty"`
	Region            string   `json:"region,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	LastDowntimeStart int64    `json:"last_downtime_start,omitempty"`
	LastDowntimeEnd   int64    `json:"last_downtime_end,omitempty"`
	CreatedAt         int64    `json:"created_at,omitempty"`
	ModifiedAt        int64    `json:"modified_at,omitempty"`
}

// TMSCheckDetailResponse represents the  JSON response for a TMS Check from the Pingdom API.
type TMSCheckDetailResponse struct {
	TMSCheck
	ID                int    `json:"id,omitempty"`
	Type              string `json:"type,omitempty"`
	LastDowntimeStart int64  `json:"last_downtime_start,omitempty"`
	LastDowntimeEnd   int64  `json:"last_downtime_end,omitempty"`
	CreatedAt         int64  `json:"created_at,omitempty"`
	ModifiedAt        int64  `json:"modified_at,omitempty"`
	Status            string `json:"status,omitempty"`
}
type TMSCheckStatusReportResponse struct {
	CheckID int              `json:"check_id,omitempty"`
	Name    string           `json:"name,omitempty"`
	States  []TMSCheckStatus `json:"states,omitempty"`
}

type TMSCheckStatus struct {
	ErrorInStep int    `json:"error_in_step,omitempty"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Message     string `json:"message,omitempty"`
	Status      string `json:"status,omitempty"`
}

type TMSCheckPerformanceReportResponse struct {
	CheckID    int                `json:"check_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Resolution string             `json:"resolution,omitempty"`
	Intervals  []TMSCheckInterval `json:"intervals,omitempty"`
}

type TMSCheckInterval struct {
	AverageResponse int64                `json:"average_response,omitempty"`
	Downtime        int64                `json:"downtime,omitempty"`
	From            string               `json:"from,omitempty"`
	Steps           []TMSCheckStepReport `json:"steps,omitempty"`
	Unmonitored     int64                `json:"unmonitored,omitempty"`
	Uptime          int64                `json:"uptime,omitempty"`
}

type TMSCheckStepReport struct {
	AverageResponse int64        `json:"average_response,omitempty"`
	Step            TMSCheckStep `json:"step,omitempty"`
}

type listTMSChecksJSONResponse struct {
	TMSChecks []TMSCheckResponse `json:"checks"`
}

type tmsChecksDetailJSONResponse struct {
	TMSCheck *TMSCheckDetailResponse `json:"check"`
}

type tmsChecksStatusReportJSONResponse struct {
	Report *TMSCheckStatusReportResponse `json:"report"`
}

type tmsChecksStatusReportsJSONResponse struct {
	Reports []TMSCheckStatusReportResponse `json:"report"`
}

type tmsChecksPerformanceReportJSONResponse struct {
	Report *TMSCheckPerformanceReportResponse `json:"report"`
}

type errorJSONResponse struct {
	Error *PingdomError `json:"error"`
}
