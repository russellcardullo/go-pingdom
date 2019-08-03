package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTsmServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms.recipes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"recipes": [
				{
                                        "name": "example.com",
                                        "status": "SUCCESSFUL",
                                        "kitchen": "eu",
                                        "active": "YES",
                                        "created_at": 1560949555,
                                        "interval": 5
				},
				{
                                        "name": "mydomain.com",
                                        "status": "SUCCESSFUL",
                                        "kitchen": "us-west",
                                        "active": "YES",
                                        "created_at": 1560935224,
                                        "interval": 5
		  	        }
			]
		}`)
	})
	want := []TsmResponse{
		{
			Name:      "example.com",
			Status:    "SUCCESSFUL",
			Kitchen:   "eu",
			Active:    "YES",
			CreatedAt: 1560949555,
			Interval:  5,
		},
		{
			Name:      "mydomain.com",
			Status:    "SUCCESSFUL",
			Kitchen:   "us-west",
			Active:    "YES",
			CreatedAt: 1560935224,
			Interval:  5,
		},
	}

	tsm, err := client.Tsm.List()
	assert.NoError(t, err)
	assert.Equal(t, want, tsm, "Tsm.List() should return correct result")
}
