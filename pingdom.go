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

// Client represents a client to the Pingdom API.  This package also
// provides a NewClient function for convenience to initialize a client
// with default parameters.
type Client struct {
	User     string
	Password string
	APIKey   string
	BaseURL  *url.URL
	client   *http.Client
}

// PingdomResponse represents a general response from the Pingdom API
type PingdomResponse struct {
	Message string `json:"message"`
}

// private types used to unmarshall json responses from pingdom

type checkResponse struct {
	Check Check `json:"check"`
}

type listChecksResponse struct {
	Checks []Check `json:"checks"`
}

type pingdomErrorResponse struct {
	Error pingdomError `json:"error"`
}

type pingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

func (r *pingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// NewClient returns a Pingdom client with a default base URL and HTTP client
func NewClient(user string, password string, key string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{user, password, key, baseURL, http.DefaultClient}
	return c
}

// NewRequest makes a new HTTP Request.  The method param should be an HTTP method in
// all caps such as GET, POST, PUT, DELETE.  The rsc param should correspond with
// a restful resource.  Params can be passed in as a map of strings
// Usually users of the client can use one of the convenience methods such as
// ListChecks, etc but this method is provided to allow for making other
// API calls that might not be built in.
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

// Takes an HTTP response and determines whether it was successful.
// Returns nil if the HTTP status code is within the 2xx range.  Returns
// an error otherwise.
func validateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	m := &pingdomErrorResponse{}
	err := json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return err
	}

	return &m.Error
}

// Return a list of checks from Pingdom.
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

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listChecksResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m.Checks, err
}

// Create a new check.  This function will validate the given check param
// to ensure that it contains correct values before submitting the request
// Returns a Check object representing the response from Pingdom.  Note
// that Pingdom does not return a full check object so in the returned
// object you should only use the ID field.
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

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &checkResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return &m.Check, err

}

// ReadCheck returns detailed information about a pingdom check given its ID.
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

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &checkResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return &m.Check, err
}

// UpdateCheck will update the check represented by the given ID with the values
// in the given check.  You should submit the complete list of values in
// the given check parameter, not just those that have changed.
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

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &PingdomResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m, err
}

// DeleteCheck will delete the check for the given ID.
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

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &PingdomResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	return m, err
}
