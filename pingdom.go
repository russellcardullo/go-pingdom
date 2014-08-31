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

type Message struct {
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

func (pc *Client) ListChecks() ([]Check, error) {
	req, _ := http.NewRequest("GET", pc.BaseURL.String()+"/api/2.0/checks", nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.Key)
	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 { // OK
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		//fmt.Println("Body:", bodyString)
		m := &ListChecksResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		if err != nil {
			return nil, err
		}
		return m.Checks, err
	}
	return nil, err
}

func (pc *Client) CreateCheck(check HttpCheck) (*Check, error) {
	baseUrl, err := url.Parse(pc.BaseURL.String() + "/api/2.0/checks")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("name", check.Name)
	params.Set("host", check.Host)
	params.Set("type", "http")

	baseUrl.RawQuery = params.Encode()

	req, _ := http.NewRequest("POST", baseUrl.String(), nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.Key)
	//fmt.Println("Req:", req)

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		//fmt.Println("Body:", bodyString)

		m := &Message{}
		err := json.Unmarshal([]byte(bodyString), &m)
		if err != nil {
			return nil, err
		}
		return &m.Check, err
	}
	return nil, err
}

func (pc *Client) ReadCheck(id string) (*Check, error) {
	baseUrl, err := url.Parse(pc.BaseURL.String() + "/api/2.0/checks/" + id)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", baseUrl.String(), nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.Key)
	//fmt.Println("Req:", req)

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		//fmt.Println("Body:", bodyString)

		m := &Message{}
		err := json.Unmarshal([]byte(bodyString), &m)
		if err != nil {
			return nil, err
		}
		return &m.Check, err
	}
	return nil, err
}

func (pc *Client) UpdateCheck(id string, name string, host string) (*PingdomResponse, error) {
	baseUrl, err := url.Parse(pc.BaseURL.String() + "/api/2.0/checks/" + id)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("host", host)
	baseUrl.RawQuery = params.Encode()

	req, _ := http.NewRequest("PUT", baseUrl.String(), nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.Key)
	//fmt.Println("Req:", req)

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		//fmt.Println("Body:", bodyString)

		m := &PingdomResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		if err != nil {
			return nil, err
		}
		return m, err
	}
	return nil, err
}

func (pc *Client) DeleteCheck(id string) (*PingdomResponse, error) {
	baseUrl, err := url.Parse(pc.BaseURL.String() + "/api/2.0/checks/" + id)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("DELETE", baseUrl.String(), nil)
	req.SetBasicAuth(pc.User, pc.Password)
	req.Header.Add("App-Key", pc.Key)
	//fmt.Println("Req:", req)

	resp, err := pc.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		//fmt.Println("Body:", bodyString)

		m := &PingdomResponse{}
		err := json.Unmarshal([]byte(bodyString), &m)
		if err != nil {
			return nil, err
		}
		return m, err
	}
	return nil, err
}
