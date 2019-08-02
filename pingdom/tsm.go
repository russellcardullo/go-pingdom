package pingdom

import (
	"encoding/json"
	"io/ioutil"
)

// TsmService provides an interface to Pingdom transactions.
type TsmService struct {
	client *Client
}

// TODO
func (cs *TsmService) List(params ...map[string]string) ([]TsmResponse, error) {
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

	p := &listTsmJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &p)

	return p.Tsm, err
}
