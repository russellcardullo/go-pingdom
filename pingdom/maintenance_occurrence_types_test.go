package pingdom

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListOccurrenceQuery(t *testing.T) {
	q := ListOccurrenceQuery{
		From:          1,
		To:            2,
		MaintenanceId: 3,
	}
	assert.Equal(t, map[string]string{
		"from":          "1",
		"to":            "2",
		"maintenanceid": "3",
	}, q.toParams())
}

func TestOccurrenceValid(t *testing.T) {
	o := Occurrence{
		To: 1,
	}
	assert.Error(t, o.Valid())

	o.From = 1
	o.To = 0
	assert.Error(t, o.Valid())
}

func TestRenderForRESTAPIJSON(t *testing.T) {
	o := Occurrence{
		From: 1,
		To:   2,
	}
	m := map[string]int{}
	err := json.Unmarshal([]byte(o.RenderForJSONAPI()), &m)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{
		"from": 1,
		"to":   2,
	}, m)
}
