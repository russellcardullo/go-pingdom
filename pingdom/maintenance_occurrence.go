package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type OccurrenceService struct {
	client *Client
}

func (os *OccurrenceService) List(query ListOccurrenceQuery) ([]Occurrence, error) {
	params := query.toParams()
	req, err := os.client.NewRequest("GET", "/maintenance.occurrences", params)
	if err != nil {
		return nil, err
	}

	resp, err := os.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := &listOccurrenceResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Occurrences, err
}

func (os *OccurrenceService) Read(id int64) (*Occurrence, error) {
	req, err := os.client.NewRequest("GET", "/maintenance.occurrences/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, err
	}

	t := &readOccurrenceResponse{}
	_, err = os.client.Do(req, t)
	if err != nil {
		return nil, err
	}

	return &t.Occurrence, err
}

// Update is used to update an existing Occurrence. Only the 'From',
// and 'To' fields can be updated.
func (os *OccurrenceService) Update(id int64, occurrence Occurrence) (*PingdomResponse, error) {
	if err := occurrence.Valid(); err != nil {
		return nil, err
	}

	req, err := os.client.NewJSONRequest("PUT", "/maintenance.occurrences/"+strconv.FormatInt(id, 10), occurrence.RenderForJSONAPI())
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = os.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// MultiDelete will delete the Occurrence for the given ID.
func (os *OccurrenceService) MultiDelete(ids []int64) (*PingdomResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("empty id list for multiple occurrence delete")
	}
	strIds := make([]string, 0, len(ids))
	for _, id := range ids {
		strIds = append(strIds, strconv.FormatInt(id, 10))
	}
	req, err := os.client.NewRequestMultiParamValue("DELETE", "/maintenance.occurrences", map[string][]string{
		"occurrenceids": strIds,
	})
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = os.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// Delete will delete the Occurrence for the given ID.
func (os *OccurrenceService) Delete(id int64) (*PingdomResponse, error) {
	req, err := os.client.NewRequest("DELETE", "/maintenance.occurrences/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = os.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}
