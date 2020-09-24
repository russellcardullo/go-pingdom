package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

// TmsCheckService provides an interface to Pingdom TMS checks.
type TmsCheckService struct {
	client *Client
}

// TmsCheck is an struct representing a TMS Check.
type TmsCheck struct {
	Name                     string            `json:"name"`
	Steps                    []TmsStep         `json:"steps"`
	Active                   bool              `json:"active,omitempty"`
	ContactIds               []int             `json:"contact_ids,omitempty"`
	CustomMessage            string            `json:"custom_message,omitempty"`
	IntegrationIds           []int             `json:"integration_ids,omitempty"`
	Interval                 int               `json:"interval,omitempty"`
	Metadata                 map[string]string `json:"metadata,omitempty"`
	Region                   string            `json:"region,omitempty"`
	SendNotificationWhenDown int               `json:"send_notification_when_down,omitempty"`
	SeverityLevel            string            `json:"severity_level,omitempty"`
	Tags                     string            `json:"tags,omitempty"`
	TeamIds                  []int             `json:"team_ids,omitempty"`
}

type TmsStep struct {
	Function string            `json:"fn,omitempty"`
	Args     map[string]string `json:"args,omitempty"`
}

const (
	DESC Order      = "desc"
	ASC  Order      = "asc"
	HOUR Resolution = "hour"
	DAY  Resolution = "day"
	WEEK Resolution = "week"
)

type Order string
type Resolution string

type TmsStatusReportListByIdRequest struct {
	From  *time.Time
	To    *time.Time
	Order Order
}
type TmsStatusReportListRequest struct {
	From      *time.Time
	To        *time.Time
	Order     Order
	Limit     *int
	Offset    *int
	OmitEmpty bool
}
type TmsPerformanceReportRequest struct {
	From          *time.Time
	To            *time.Time
	Order         Order
	IncludeUptime bool
	Resolution    Resolution
}

// Valid determines whether the TmsCheck contains valid fields. This can be
// used to guard against sending illegal values to the Pingdom API.
func (ts *TmsCheck) Valid() error {
	if ts.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	//if ts.Hostname == "" {
	//	return fmt.Errorf("Invalid value for `Hostname`.  Must contain non-empty string")
	//}
	//
	//if ts.Resolution != 1 && ts.Resolution != 5 && ts.Resolution != 15 &&
	//	ts.Resolution != 30 && ts.Resolution != 60 {
	//	return fmt.Errorf("invalid value %v for `Resolution`, allowed values are [1,5,15,30,60]", ts.Resolution)
	//}
	return nil
}
func (tr *TmsStatusReportListByIdRequest) Valid() error {
	if tr.To != nil && tr.From != nil && tr.To.Before(*tr.From) {
		return fmt.Errorf("from date should be earlier then to date")
	}

	switch tr.Order {
	case DESC, ASC, "":
	default:
		return fmt.Errorf("invalid order allowed values are: %s, %s", ASC, DESC)
	}

	return nil
}
func (tr *TmsStatusReportListByIdRequest) GetParams() map[string]string {
	m := map[string]string{}

	if tr.From != nil {
		m["from"] = tr.From.Format(time.RFC3339)
	}

	if tr.To != nil {
		m["to"] = tr.To.Format(time.RFC3339)
	}

	if len(tr.Order) > 0 {
		m["order"] = string(tr.Order)
	}

	return m
}
func (tr *TmsStatusReportListRequest) Valid() error {
	if tr.To != nil && tr.From != nil && tr.To.Before(*tr.From) {
		return fmt.Errorf("from date should be earlier then to date")
	}

	if tr.Offset != nil && *tr.Offset < 0 {
		return fmt.Errorf("offset should be greater equal 0")
	}
	if tr.Limit != nil && *tr.Limit <= 0 {
		return fmt.Errorf("limit should be greater 0")
	}

	switch tr.Order {
	case DESC, ASC, "":
	default:
		return fmt.Errorf("invalid order allowed values are: %s, %s", ASC, DESC)
	}
	return nil
}
func (tr *TmsStatusReportListRequest) GetParams() map[string]string {
	m := map[string]string{}

	if tr.From != nil {
		m["from"] = tr.From.Format(time.RFC3339)
	}

	if tr.To != nil {
		m["to"] = tr.To.Format(time.RFC3339)
	}

	if tr.Offset != nil {
		m["offset"] = strconv.Itoa(*tr.Offset)
	}

	if tr.Limit != nil {
		m["limit"] = strconv.Itoa(*tr.Limit)
	}

	if len(tr.Order) > 0 {
		m["order"] = string(tr.Order)
	}

	//default is false
	if tr.OmitEmpty {
		m["omit_empty"] = "true"
	}

	return m
}
func (tr *TmsPerformanceReportRequest) Valid() error {
	if tr.To != nil && tr.From != nil && tr.To.Before(*tr.From) {
		return fmt.Errorf("from date should be earlier then to date")
	}

	switch tr.Order {
	case DESC, ASC, "":
	default:
		return fmt.Errorf("invalid order allowed values are: %s, %s", ASC, DESC)
	}

	switch tr.Resolution {
	case HOUR, DAY, WEEK, "":
	default:
		return fmt.Errorf("invalid order allowed values are: %s, %s. %s", HOUR, DAY, WEEK)
	}
	return nil
}
func (tr *TmsPerformanceReportRequest) GetParams() map[string]string {
	m := map[string]string{}

	if tr.From != nil {
		m["from"] = tr.From.Format(time.RFC3339)
	}

	if tr.To != nil {
		m["to"] = tr.To.Format(time.RFC3339)
	}

	if len(tr.Order) > 0 {
		m["order"] = string(tr.Order)
	}

	//default is false
	if tr.IncludeUptime {
		m["include_uptime"] = "true"
	}

	if len(tr.Resolution) > 0 {
		m["resolution"] = string(tr.Resolution)
	}

	return m
}

// List returns a list of TMD checks from Pingdom.
// This returns type TmsCheckResponse rather than TMSCheck since the
// Pingdom API does not return a complete representation of a check.
func (cs *TmsCheckService) List(params ...map[string]string) ([]TmsCheckResponse, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/tms/check", param)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listTmsChecksJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.TmsChecks, err
}

// Create a new TMS check.
func (cs *TmsCheckService) Create(check TmsCheck) (*TmsCheckResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJsonRequest("POST", "/tms/check", nil, check)
	if err != nil {
		return nil, err
	}

	m := &TmsCheckResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// ReadCheck returns detailed information about a pingdom TMS check given its ID.
func (cs *TmsCheckService) Read(id int) (*TmsCheckResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &tmsCheckDetailsJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m.TmsCheck, err
}

// Update will update the TMS check represented by the given ID with the values
// in the given check.  You should submit the complete list of values in
// the given check parameter, not just those that have changed.
func (cs *TmsCheckService) Update(id int, tmsCheck TmsCheck) (*TmsCheckResponse, error) {
	if err := tmsCheck.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJsonRequest("PUT", "/tms/checks/"+strconv.Itoa(id), nil, tmsCheck)
	if err != nil {
		return nil, err
	}

	m := &TmsCheckResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// Delete will delete the TMS check for the given ID.
func (cs *TmsCheckService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/tms/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

//Returns a status change report for all transaction checks in the current organization
func (cs *TmsCheckService) StatusReportList(request TmsStatusReportListRequest) (*TmsStatusChangeResponse, error) {
	if err := request.Valid(); err != nil {
		return nil, err
	}
	req, err := cs.client.NewRequest("GET", "/tms/check/report/status", request.GetParams())
	if err != nil {
		return nil, err
	}
	m := &TmsStatusChangeResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

//Returns a status change report for a single transaction checks in the current organization
func (cs *TmsCheckService) StatusReportById(id int, request TmsStatusReportListByIdRequest) (*TmsStatusChangeResponse, error) {
	if err := request.Valid(); err != nil {
		return nil, err
	}
	req, err := cs.client.NewRequest("GET", fmt.Sprintf("/tms/check/%s/report/status", strconv.Itoa(id)), request.GetParams())
	if err != nil {
		return nil, err
	}
	m := &TmsStatusChangeResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

//Returns a performance report for a single transaction checks in the current organization
func (cs *TmsCheckService) PerformanceReport(id int, request TmsPerformanceReportRequest) (*TmsPerformanceReportResponse, error) {
	if err := request.Valid(); err != nil {
		return nil, err
	}
	req, err := cs.client.NewRequest("GET", fmt.Sprintf("/tms/check/%s/report/performance", strconv.Itoa(id)), request.GetParams())
	if err != nil {
		return nil, err
	}
	m := &TmsPerformanceReportResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
