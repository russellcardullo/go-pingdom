package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaintenanceServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.1/maintenance", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"maintenance": [
				{
					"description": "Maintenance N",
					"id": 85975,
					"from": 1,
					"to": 1524048059,
					"recurrencetype": "none",
					"repeatevery": 0,
					"effectiveto": 1524048059,
					"checks": {
						"uptime": [
							12345,
							23456
						],
						"tms": [
							1234,
							8975
						]
					}
				}
			]
		}`)
	})

	maintenances, err := client.Maintenances.List()
	if assert.Nil(t, err) {
		checks := MaintenanceCheckResponse{
			Uptime: []int{12345, 23456},
			Tms:    []int{1234, 8975},
		}
		want := []MaintenanceResponse{
			MaintenanceResponse{
				ID:             85975,
				Description:    "Maintenance N",
				From:           1,
				To:             1524048059,
				Recurrencetype: "none",
				Repeatevery:    0,
				Effectiveto:    1524048059,
				Checks:         checks,
			},
		}

		assert.Equal(t, want, maintenances, "Maintenances.List() should return correct result")
	}
}

func TestMaintenanceServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.1/maintenance", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"maintenance": {
				"id": 85975
			}
		}`)
	})

	m := MaintenanceWindow{
		Description: "Maintenance N",
		From:        1,
		To:          1524048059,
	}

	maintenances, err := client.Maintenances.Create(&m)
	if assert.Nil(t, err) {
		want := &MaintenanceResponse{
			ID: 85975,
		}

		assert.Equal(t, want, maintenances, "Maintenances.Create() should return correct result")
	}
}

func TestMaintenanceServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.1/maintenance/456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"maintenance": {
					"id": 456,
					"description": "Particular maintenance window",
					"from": 1497520800,
					"to": 1497574800,
					"recurrencetype": "none",
					"repeatevery": 0,
					"effectiveto": 1497574800,
					"checks": {
							"uptime": [506206, 506233, 222],
							"tms": [123, 111]
					}
			}
	}`)
	})

	maintenance, err := client.Maintenances.Read(456)
	if assert.Nil(t, err) {

		checks := MaintenanceCheckResponse{
			Uptime: []int{506206, 506233, 222},
			Tms:    []int{123, 111},
		}

		want := &MaintenanceResponse{
			ID:             456,
			Description:    "Particular maintenance window",
			From:           1497520800,
			To:             1497574800,
			Recurrencetype: "none",
			Repeatevery:    0,
			Effectiveto:    1497574800,
			Checks:         checks,
		}

		assert.Equal(t, want, maintenance, "Maintenances.Read() should return correct result")
	}

}

func TestMaintenanceServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.1/maintenance/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message":"Maintenance window successfully modified!"}`)
	})

	updateMaintenance := MaintenanceWindow{
		Description: "Updated Maintenance N",
		From:        1,
		To:          1524048061,
	}

	msg, err := client.Maintenances.Update(12345, &updateMaintenance)
	if assert.Nil(t, err) {

		want := &PingdomResponse{Message: "Maintenance window successfully modified!"}
		assert.Equal(t, want, msg, "Maintenances.Update() should return correct result")
	}
}

func TestMaintenanceServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.1/maintenance/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message":"Maintenance window successfully deleted!"}`)
	})

	msg, err := client.Maintenances.Delete(12345)
	if assert.Nil(t, err) {

		want := &PingdomResponse{Message: "Maintenance window successfully deleted!"}
		assert.Equal(t, want, msg, "Maintenances.Delete() should return correct result")
	}
}
