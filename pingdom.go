package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultBaseURL = "https://api.pingdom.com/"
)

type Client struct {
	User     string
	Password string
	APIKey   string
	BaseURL  *url.URL
	client   *http.Client
}

type CheckResponse struct {
	Check Check `json:"check"`
}

type ListChecksResponse struct {
	Checks []Check `json:"checks"`
}

type PingdomResponse struct {
	Message string `json:"message"`
}

type PingdomErrorResponse struct {
	Error PingdomError `json:"error"`
}

type PingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

func NewClient(user string, password string, key string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{user, password, key, baseURL, http.DefaultClient}
	return c
}

func (pc *Client) NewRequest(method string, rsc string, params map[string]string) (*http.Request, error) {
	baseUrl, err := url.Parse(pc.BaseURL.String() + rsc)
	if err != nil {
		return nil, err
	}

	if params != nil {
		ps := url.Values{}
		for k, v := range params {
			ps.Set(k, v)
		}
		baseUrl.RawQuery = ps.Encode()
	}

	req, err := http.NewRequest(method, baseUrl.String(), nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.APIKey)
	return req, err
}

func ValidateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	m := &PingdomErrorResponse{}
	err := json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return err
	}

	return &m.Error
}

func (pc *Client) ListChecks() ([]Check, error) {
	req, err := pc.NewRequest("GET", "/api/2.0/checks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ValidateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &ListChecksResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m.Checks, err
}

func (pc *Client) CreateCheck(check *Check) (*Check, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	req, err := pc.NewRequest("POST", "/api/2.0/checks", check.Params())
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ValidateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &CheckResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return &m.Check, err

}

func (pc *Client) ReadCheck(id int) (*Check, error) {
	req, err := pc.NewRequest("GET", "/api/2.0/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ValidateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &CheckResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return &m.Check, err
}

func (pc *Client) UpdateCheck(id int, check *Check) (*PingdomResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	params := check.Params()
	delete(params, "type")
	req, err := pc.NewRequest("PUT", "/api/2.0/checks/"+strconv.Itoa(id), params)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ValidateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &PingdomResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m, err
}

func (pc *Client) DeleteCheck(id int) (*PingdomResponse, error) {
	req, err := pc.NewRequest("DELETE", "/api/2.0/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ValidateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &PingdomResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m, err
}
