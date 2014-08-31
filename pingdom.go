package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.pingdom.com/"
)

type Client struct {
	User     string
	Password string
	Key      string
	BaseURL  *url.URL
	client   *http.Client
}

type HttpCheck struct {
	Name string
	Host string
}

type Check struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	Resolution               int    `json:"resolution,omitempty"`
	SendToEmail              bool   `json:"sendtoemail,omitempty"`
	SendToTwitter            bool   `json:"sendtotwitter,omitempty"`
	SendToIPhone             bool   `json:"sendtoiphone,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
	Created                  int64  `json:"created,omitempty"`
	Hostname                 string `json:"hostname,omitempty"`
	Status                   string `json:"status,omitempty"`
	LastErrorTime            int64  `json:"lasterrortime,omitempty"`
	LastTestTime             int64  `json:"lasttesttime,omitempty"`
	LastResponseTime         int64  `json:"lastresponsetime,omitempty"`
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
	req.Header.Add("App-Key", pc.Key)
	return req, err
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

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		m := &ListChecksResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		return m.Checks, err
	}
	return nil, err
}

func (pc *Client) CreateCheck(check HttpCheck) (*Check, error) {
	params := map[string]string{
		"name": check.Name,
		"host": check.Host,
		"type": "http",
	}
	req, err := pc.NewRequest("POST", "/api/2.0/checks", params)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		m := &CheckResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		return &m.Check, err
	}
	return nil, err
}

func (pc *Client) ReadCheck(id string) (*Check, error) {
	req, err := pc.NewRequest("GET", "/api/2.0/checks/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		m := &CheckResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		return &m.Check, err
	}
	return nil, err
}

func (pc *Client) UpdateCheck(id string, name string, host string) (*PingdomResponse, error) {
	params := map[string]string{
		"name": name,
		"host": host,
	}
	req, err := pc.NewRequest("PUT", "/api/2.0/checks/"+id, params)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		m := &PingdomResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		return m, err
	}
	return nil, err
}

func (pc *Client) DeleteCheck(id string) (*PingdomResponse, error) {
	req, err := pc.NewRequest("DELETE", "/api/2.0/checks/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		m := &PingdomResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		return m, err
	}
	return nil, err
}
