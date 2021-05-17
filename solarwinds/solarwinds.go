package solarwinds

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultBaseURL              = "https://my.solarwinds.cloud"
	graphQLEndpoint             = "/common/graphql"
	headerNameSetCookie         = "Set-Cookie"
	cookieNameSwicus            = "swicus"
	cookieNameSwiSettings       = "swi-settings"
	headerNameCSRFToken         = "X-CSRF-Token"
	EnvSolarwindsUser           = "SOLARWINDS_USER"
	EnvSolarwindsPassword       = "SOLARWINDS_PASSWD"
	EnvSolarwindsOrganizationId = "SOLARWINDS_ORG_ID"
)

type Client struct {
	csrfToken         string
	swiSettings       string
	email             string
	password          string
	organizationId    string
	client            *http.Client
	baseURL           string
	InvitationService *InvitationService
	ActiveUserService *ActiveUserService
	UserService       *UserService
}

type ClientConfig struct {
	Username       string // Typically this is an email
	Password       string
	OrganizationId string
	BaseURL        string // For UT
}

type loginPayload struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	LoginQueryParams string `json:"loginQueryParams"`
}

type loginResult struct {
	Swicus      string
	RedirectURL string
}

// Does not involve any network interactions.
func NewClient(config ClientConfig) (*Client, error) {
	var baseURLToUse *url.URL
	var err error
	if config.BaseURL == "" {
		baseURLToUse, err = url.Parse(defaultBaseURL)
	} else {
		baseURLToUse, err = url.Parse(config.BaseURL)
	}
	if err != nil {
		return nil, err
	}

	username := config.Username
	if username == "" {
		username = os.Getenv(EnvSolarwindsUser)
	}

	password := config.Password
	if password == "" {
		password = os.Getenv(EnvSolarwindsPassword)
	}

	organizationId := config.OrganizationId
	if organizationId == "" {
		organizationId = os.Getenv(EnvSolarwindsOrganizationId)
	}

	c := &Client{
		email:          username,
		password:       password,
		organizationId: organizationId,
		baseURL:        baseURLToUse.String(),
	}
	c.client = http.DefaultClient
	c.InvitationService = &InvitationService{client: c}
	c.ActiveUserService = &ActiveUserService{client: c}
	c.UserService = &UserService{
		ActiveUserService: c.ActiveUserService,
		InvitationService: c.InvitationService,
	}
	return c, nil
}

func (c *Client) Init() error {
	auth, err := c.login()
	if err != nil {
		return err
	}
	if err := c.obtainSwiSettings(); err != nil {
		return err
	}
	return c.obtainToken(auth)
}

func (c *Client) NewRequest(method string, rsc string, params io.Reader) (*http.Request, error) {
	baseURL, err := url.Parse(c.baseURL + rsc)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, baseURL.String(), params)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  cookieNameSwiSettings,
		Value: c.swiSettings,
	})
	req.Header.Set(headerNameCSRFToken, c.csrfToken)
	return req, err
}

func (c *Client) MakeGraphQLRequest(graphQLRequest *GraphQLRequest) (*GraphQLResponse, error) {
	body, err := ToJsonNoEscape(graphQLRequest)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", graphQLEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	graphQLResp, err := NewGraphQLResponse(resp.Body, graphQLRequest.ResponseType)
	if err != nil {
		return nil, err
	}
	if !graphQLResp.isSuccess() {
		return nil, fmt.Errorf("request failed with message: %v", graphQLResp.message())
	}
	return graphQLResp, err
}

// login provides user credentials and gets a 'swicus' value in return. This value serves
// as a proof that one has been authenticated.
func (c *Client) login() (*loginResult, error) {
	params := map[string]string{
		"response_type": "code",
		"scope":         "openid swicus",
		"client_id":     "adminpanel",
		"redirect_uri":  "https://my.solarwinds.cloud/common/auth/callback",
		"state":         RandString(10),
	}
	paramsToUse := url.Values{}
	for k, v := range params {
		paramsToUse.Add(k, v)
	}
	payload := loginPayload{
		Email:            c.email,
		Password:         c.password,
		LoginQueryParams: paramsToUse.Encode(),
	}
	body, err := ToJsonNoEscape(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+"/v1/login", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("visit callback failed, status %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	result := &loginResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	if swicus, err := retrieveCookie(resp, cookieNameSwicus); err != nil {
		return nil, err
	} else {
		result.Swicus = swicus
	}
	return result, nil
}

// obtainSwiSettings is used to retrieve 'swi-settings' cookie. The value is contained
// in a redirect response. This step does not depend on any previous steps.
func (c *Client) obtainSwiSettings() error {
	resp, err := http.Get(c.baseURL + "/common/login")
	if err != nil {
		return err
	}
	swiSettings, err := retrieveCookie(resp.Request.Response, cookieNameSwiSettings)
	if err != nil {
		return err
	}
	c.swiSettings = swiSettings
	return nil
}

// obtainToken uses the 'swicus' and 'swi-settings' to obtain a CSRF token.
func (c *Client) obtainToken(auth *loginResult) error {
	var url string
	if c.organizationId != "" {
		url = fmt.Sprintf("%s/%s/%s/users", c.baseURL, "settings", c.organizationId)
	} else {
		url = c.baseURL + "/settings"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:  cookieNameSwicus,
		Value: auth.Swicus,
	})
	req.AddCookie(&http.Cookie{
		Name:  cookieNameSwiSettings,
		Value: c.swiSettings,
	})
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("visit callback URL failed, status %d", resp.StatusCode)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	if token, err := extractCSRFToken(doc); err != nil {
		return err
	} else {
		c.csrfToken = token
	}
	return nil
}

func extractCSRFToken(start *html.Node) (string, error) {
	var token string
	var head *html.Node
	if first := start.FirstChild; first.Type == html.DoctypeNode {
		head = first.NextSibling.FirstChild
	} else {
		head = first.FirstChild
	}
outer:
	for c := head.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "meta" && len(c.Attr) == 2 {
			for _, attr := range c.Attr {
				if attr.Key == "name" && attr.Val != "csrf-token" {
					continue outer
				}
			}
			for _, attr := range c.Attr {
				if attr.Key == "content" {
					token = attr.Val
				}
			}
			if token != "" {
				break
			}
		}
	}
	if token == "" {
		return "", errors.New("response of callback URL does not contain CSRF token")
	}
	return token, nil
}
