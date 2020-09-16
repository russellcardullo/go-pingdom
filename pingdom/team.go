package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// TeamService provides an interface to Pingdom teams.
type TeamService struct {
	client *Client
}

// TeamAPI is an interface representing a Pingdom team.
type TeamAPI interface {
	RenderForJSONAPI() string
	Valid() error
}

// List return a list of teams from Pingdom.
func (cs *TeamService) List() ([]TeamResponse, error) {
	req, err := cs.client.NewRequest("GET", "/alerting/teams", nil)
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

	t := &listTeamsJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &t)

	return t.Teams, err
}

// Read return a team object from Pingdom.
func (cs *TeamService) Read(id int) (*TeamResponse, error) {
	req, err := cs.client.NewRequest("GET", "/alerting/teams/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	t := &teamDetailsJSONResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}

	return t.Team, err
}

// Create is used to create a new team.
func (cs *TeamService) Create(team TeamAPI) (*TeamResponse, error) {
	if err := team.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewJSONRequest("POST", "/alerting/teams", team.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	t := &teamDetailsJSONResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t.Team, err
}

// Update is used to update existing team.
func (cs *TeamService) Update(id int, team TeamAPI) (*TeamResponse, error) {
	req, err := cs.client.NewJSONRequest("PUT", "/alerting/teams/"+strconv.Itoa(id), team.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	t := &teamDetailsJSONResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t.Team, err
}

// Delete will delete the Team for the given ID.
func (cs *TeamService) Delete(id int) (*TeamDeleteResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/alerting/teams/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	t := &TeamDeleteResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t, err
}
