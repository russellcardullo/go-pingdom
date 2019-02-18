package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// CheckRecipe provides an interface to Pingdom checks
type CheckRecipe struct {
	client *Client
}

func (cs *CheckRecipe) List(params ...map[string]string) (*listRecipesJsonResponse, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/tms.recipes", param)
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
	m := &listRecipesJsonResponse{}

	err = json.Unmarshal([]byte(bodyString), &m)
	return m, err
}

func (cs *CheckRecipe) Summary(id int, params ...map[string]string) (*RecipeSummaryResponse, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/tms.summary.performance/"+strconv.Itoa(id), param)
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
	m := &RecipeSummaryResponse{}

	err = json.Unmarshal([]byte(bodyString), &m)
	fmt.Print(m.Summary.Hours)
	return m, err
}
