package pingdom

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestOccurrenceServiceList(t *testing.T) {
	setup()
	defer teardown()

	respStr := `
{
  "occurrences": [
    {
      "id": 6110986,
      "maintenanceid": 224724,
      "from": 1617699622,
      "to": 1617703222,
      "duration": 60,
      "durationunit": "minute"
    },
    {
      "id": 6110987,
      "maintenanceid": 224724,
      "from": 1618304422,
      "to": 1618308022,
      "duration": 60,
      "durationunit": "minute"
    }
  ]
}`
	mux.HandleFunc("/maintenance.occurrences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, respStr)
	})
	want := listOccurrenceResponse{}
	err := json.Unmarshal([]byte(respStr), &want)
	assert.NoError(t, err)

	occurrences, err := client.Occurrences.List(ListOccurrenceQuery{})
	assert.NoError(t, err)
	assert.Equal(t, want.Occurrences, occurrences, "Occurrence.List() should return correct result")
}

func TestOccurrenceServiceRead(t *testing.T) {
	setup()
	defer teardown()

	respStr := `
{
  "occurrence": {
      "id": 6110986,
      "maintenanceid": 224724,
      "from": 1617699622,
      "to": 1617703222,
      "duration": 60,
      "durationunit": "minute"
    }
}
`
	mux.HandleFunc("/maintenance.occurrences/6110986", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, respStr)
	})

	want := readOccurrenceResponse{}
	err := json.Unmarshal([]byte(respStr), &want)
	assert.NoError(t, err)

	occurrence, err := client.Occurrences.Read(6110986)
	assert.NoError(t, err)
	assert.Equal(t, want.Occurrence, *occurrence, "Occurrence.Read() should return correct result")
}

func TestOccurrenceServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	respStr := `
{

    "message": "Occurrence successfully modified!"

}
`
	mux.HandleFunc("/maintenance.occurrences/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = fmt.Fprint(w, respStr)
	})

	want := &PingdomResponse{}
	err := json.Unmarshal([]byte(respStr), want)
	assert.NoError(t, err)

	update := Occurrence{
		From: 1,
		To:   2,
	}
	msg, err := client.Occurrences.Update(12345, update)
	assert.NoError(t, err)
	assert.Equal(t, want, msg, "Occurrence.Update() should return correct result")

}

func TestOccurrenceServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	respStr := `
{

    "message": "Occurrence successfully deleted!"

}
`
	mux.HandleFunc("/maintenance.occurrences/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		_, _ = fmt.Fprint(w, respStr)
	})
	want := &PingdomResponse{}
	err := json.Unmarshal([]byte(respStr), want)
	assert.NoError(t, err)

	msg, err := client.Occurrences.Delete(1234)
	assert.NoError(t, err)
	assert.Equal(t, want, msg, "Occurrence.Delete() should return correct result")
}

func TestOccurrenceServiceMultiDelete(t *testing.T) {
	setup()
	defer teardown()

	respStr := `
{

    "message": "5 occurrences successfully deleted."

}
`
	idsToDelete := []int64{1, 2, 3, 4, 5}

	mux.HandleFunc("/maintenance.occurrences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		actualIds := r.URL.Query()["occurrenceids"]
		assert.Equal(t, len(idsToDelete), len(actualIds))
		_, _ = fmt.Fprint(w, respStr)
	})
	want := &PingdomResponse{}
	err := json.Unmarshal([]byte(respStr), want)
	assert.NoError(t, err)

	msg, err := client.Occurrences.MultiDelete(idsToDelete)
	assert.NoError(t, err)
	assert.Equal(t, want, msg, "Occurrence.MultiDelete() should return correct result")
}
