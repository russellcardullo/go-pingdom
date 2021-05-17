package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type TMSCheckService struct {
	client *Client
}

// TMSCheckSAPI is an interface representing a Pingdom team.
type TMSCheckSAPI interface {
	RenderForJSONAPI() string
	Valid() error
}

// List return a list of TMS checks from Pingdom.
func (cs *TMSCheckService) List(params ...map[string]string) ([]TMSCheckResponse, error) {

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

	t := &listTMSChecksJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &t)

	return t.TMSChecks, err
}

func (cs *TMSCheckService) Read(id int) (*TMSCheckDetailResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms/check/"+strconv.Itoa(id), nil)
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

	t := &tmsChecksDetailJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &t)

	return t.TMSCheck, err
}

func (cs *TMSCheckService) Create(tmsCheck *TMSCheck) (*TMSCheckDetailResponse, error) {
	if err := tmsCheck.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJSONRequest("POST", "/tms/check", tmsCheck.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	t := &tmsChecksDetailJSONResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t.TMSCheck, err
}

func (cs *TMSCheckService) Update(id int, tmsCheck *TMSCheck) (*TMSCheckDetailResponse, error) {
	if err := tmsCheck.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJSONRequest("PUT", "/tms/check/"+strconv.Itoa(id), tmsCheck.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	t := &tmsChecksDetailJSONResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t.TMSCheck, err

}

func (cs *TMSCheckService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/tms/check/"+strconv.Itoa(id), nil)
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

func (cs *TMSCheckService) GetStatusReport(id int, params map[string]string) (*TMSCheckStatusReportResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms/check/"+strconv.Itoa(id)+"/report/status", params)
	if err != nil {
		return nil, err
	}

	m := &tmsChecksStatusReportJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Report, err
}

func (cs *TMSCheckService) ListStatusReports(params map[string]string) ([]TMSCheckStatusReportResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms/check/report/status", params)
	if err != nil {
		return nil, err
	}

	m := &tmsChecksStatusReportsJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Reports, err
}

func (cs *TMSCheckService) GetPerfomanceReport(id int, params map[string]string) (*TMSCheckPerformanceReportResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms/check/"+strconv.Itoa(id)+"/report/performance", params)
	if err != nil {
		return nil, err
	}

	m := &tmsChecksPerformanceReportJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Report, err
}
