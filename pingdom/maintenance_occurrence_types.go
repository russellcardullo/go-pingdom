package pingdom

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Occurrence struct {
	Id            int64  `json:"id"`
	MaintenanceId int64  `json:"maintenanceid"`
	From          int64  `json:"from"`
	To            int64  `json:"to"`
	Duration      int    `json:"duration"`
	DurationUnit  string `json:"durationunit"`
}

type ListOccurrenceQuery struct {
	From          int64 `json:"from"`
	To            int64 `json:"to"`
	MaintenanceId int64 `json:"maintenanceid"`
}

type listOccurrenceResponse struct {
	Occurrences []Occurrence `json:"occurrences"`
}

type readOccurrenceResponse struct {
	Occurrence Occurrence `json:"occurrence"`
}

func (q *ListOccurrenceQuery) toParams() map[string]string {
	m := map[string]string{}
	if q.From != 0 {
		m["from"] = strconv.FormatInt(q.From, 10)
	}
	if q.To != 0 {
		m["to"] = strconv.FormatInt(q.To, 10)
	}
	if q.MaintenanceId != 0 {
		m["maintenanceid"] = strconv.FormatInt(q.MaintenanceId, 10)
	}
	return m
}

func (o *Occurrence) Valid() error {
	if o.From == 0 {
		return fmt.Errorf("Invalid value for `From`.  Must contain time")
	}

	if o.To == 0 {
		return fmt.Errorf("Invalid value for `To`.  Must contain time")
	}

	return nil
}

func (o *Occurrence) RenderForJSONAPI() string {
	b := map[string]interface{}{
		"from": o.From,
		"to":   o.To,
	}
	jsonBody, _ := json.Marshal(b)
	return string(jsonBody)
}
