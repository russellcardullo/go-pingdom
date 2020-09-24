package pingdom

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTmsCheckServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
    "checks": [
        {
            "id": 84591,
            "name": "TestLoginRedirect",
            "type": "script",
            "active": true,
            "status": "successful",
            "created_at": 1572197307,
            "interval": 10,
            "region": "au",
            "modified_at": 1599061779,
            "tags": []
        }
    ],
    "limit": 1000,
    "offset": 0
}`)
	})
	want := []TmsCheckResponse{
		{
			ID:                       84591,
			Name:                     "TestLoginRedirect",
			Active:                   true,
			Interval:                 10,
			Region:                   "au",
			SendNotificationWhenDown: 0,
			Tags:                     []string{},
		},
	}

	check, err := client.TmsChecks.List(nil)
	assert.NoError(t, err)
	assert.Equal(t, want, check)
}
func TestTmsCheckServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
		"id": 1003,
		"name": "Test redirect"
		}`)
	})

	want := &TmsCheckResponse{
		ID:   1003,
		Name: "Test redirect",
	}

	checks, err := client.TmsChecks.Create(TmsCheck{
		Name: "Test redirect",
		Steps: []TmsStep{
			{
				Function: "go_to",
				Args:     map[string]string{"url": "https://example.com/"},
			},
			{
				Function: "url",
				Args:     map[string]string{"url": "https://example.com/redirected"},
			},
		}})
	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}
func TestTmsCheckServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/checks/85975", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
    "check": {
        "id": 85975,
        "type": "script",
        "active": true,
        "status": "successful",
        "created_at": 1572197307,
        "interval": 10,
        "region": "au",
        "modified_at": 1599061779,
        "name": "TestRedirect",
        "steps": [
            {
                "fn": "go_to",
                "args": {
                    "url": "https://example.com"
                }
            },
            {
                "fn": "url",
                "args": {
                    "url": "https://example.com/redirected"
                }
            }
        ],
        "contact_ids": [],
        "team_ids": [123456],
        "integration_ids": [
            12345
        ],
        "custom_message": "",
        "send_notification_when_down": 1,
        "severity_level": "high",
        "tags": []
    }
}`)
	})

	want := &TmsCheckResponse{
		ID:                       85975,
		Name:                     "TestRedirect",
		SendNotificationWhenDown: 1,
		TeamIds:                  []int{123456},
		Tags:                     []string{},
		Active:                   true,
		ContactIds:               []int{},
		IntegrationIds:           []int{12345},
		SeverityLevel:            "high",
		Region:                   "au",
		Interval:                 10,
		Steps: []TmsStep{
			{
				Function: "go_to",
				Args:     map[string]string{"url": "https://example.com"},
			},
			{
				Function: "url",
				Args:     map[string]string{"url": "https://example.com/redirected"},
			},
		},
	}

	check, err := client.TmsChecks.Read(85975)
	assert.NoError(t, err)
	assert.Equal(t, want, check)
}
func TestTmsCheckServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "active": true,
  "contact_ids": [
    12345678,
    19876654
  ],
  "created_at": 1553070682,
  "modified_at": 1553070968,
  "custom_message": "My custom message",
  "interval": 10,
  "name": "Updated Check",
  "region": "us-west",
  "send_notification_when_down": 1,
  "severity_level": "low",
  "status": "successful",
  "steps": [
    {
      "args": {
        "checkbox": "string",
        "element": "string",
        "form": "string",
        "input": "string",
        "option": "string",
        "password": "string",
        "radio": "string",
        "seconds": "string",
        "select": "string",
        "url": "http://www.google.com",
        "username": "string",
        "value": "string"
      },
      "fn": "go_to"
    }
  ],
  "team_ids": [
    12345678,
    135790
  ],
  "integration_ids": [
    1234,
    1359
  ],
  "metadata": {
    "width": 1950,
    "height": 1080,
    "disableWebSecurity": true,
    "authentications": {
      "httpAuthentications": [
        {
          "credentials": {
            "password": "secret",
            "userName": "admin"
          },
          "host": "https://example.com/auth"
        }
      ]
    }
  },
  "tags": [
    "tag1",
    "tag2"
  ],
  "type": [
    "script"
  ]
}`)
	})

	updateCheck := TmsCheck{
		Name: "Updated Check",
	}

	want := &TmsCheckResponse{
		Name: "Updated Check",
		Steps: []TmsStep{
			{
				Function: "go_to",
				Args: map[string]string{
					"checkbox": "string",
					"form":     "string",
					"password": "string",
					"radio":    "string",
					"seconds":  "string",
					"element":  "string",
					"input":    "string",
					"option":   "string",
					"select":   "string",
					"url":      "http://www.google.com",
					"username": "string",
					"value":    "string",
				},
			},
		},
		Active:         true,
		ContactIds:     []int{12345678, 19876654},
		CustomMessage:  "My custom message",
		IntegrationIds: []int{1234, 1359},
		Interval:       10,
		Metadata: map[string]interface{}{
			"authentications": map[string]interface{}{
				"httpAuthentications": []interface{}{
					map[string]interface{}{
						"credentials": map[string]interface{}{
							"userName": "admin",
							"password": "secret",
						},
						"host": "https://example.com/auth",
					},
				},
			},
			"width":              1950.,
			"height":             1080.,
			"disableWebSecurity": true,
		},
		Region:                   "us-west",
		SendNotificationWhenDown: 1,
		SeverityLevel:            "low",
		Tags:                     []string{"tag1", "tag2"},
		TeamIds:                  []int{12345678, 135790},
	}

	msg, err := client.TmsChecks.Update(12345, updateCheck)
	assert.NoError(t, err)
	assert.Equal(t, want, msg)
}
func TestTmsCheckServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
"message": "Deletion of check 12345 was successful"
}`)
	})

	want := &PingdomResponse{Message: "Deletion of check 12345 was successful"}

	msg, err := client.TmsChecks.Delete(12345)
	assert.NoError(t, err)
	assert.Equal(t, want, msg)
}
func TestTmsCheckStatusChangeReportList(t *testing.T) {
	setup()
	defer teardown()

	request := TmsStatusReportListRequest{
		Order: ASC,
	}

	from, _ := time.Parse(time.RFC3339, "2020-07-10T10:51:55.000Z")
	to, _ := time.Parse(time.RFC3339, "2020-07-14T07:25:15.000Z")
	expectedResponse := TmsStatusChangeResponse{
		Report: TmsStatusChange{
			CheckId: 123,
			Name:    "My awesome check",
			States: []TmsState{
				{
					From:    from,
					To:      to,
					Status:  "down",
					Message: "URL should be 'http://www.example12345.com/' but is 'http://www.example.com/'.",
				},
			},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/tms/check/report/status"), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
  "report": {
    "check_id": 123,
    "name": "My awesome check",
    "states": [
      {
        "status": "down",
        "from": "2020-07-10T10:51:55.000Z",
        "to": "2020-07-14T07:25:15.000Z",
        "error_in_step": 2,
        "message": "URL should be 'http://www.example12345.com/' but is 'http://www.example.com/'."
      }
    ]
  }
}`)
	})

	resp, err := client.TmsChecks.StatusReportList(request)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, *resp)
}
func TestTmsCheckStatusChangeReportSingleTransactionCheck(t *testing.T) {
	id := 1234
	setup()
	defer teardown()

	request := TmsStatusReportListByIdRequest{}

	from, _ := time.Parse(time.RFC3339, "2020-07-10T10:51:55.000Z")
	to, _ := time.Parse(time.RFC3339, "2020-07-14T07:25:15.000Z")
	expectedResponse := TmsStatusChangeResponse{
		Report: TmsStatusChange{
			CheckId: 123,
			Name:    "My awesome check",
			States: []TmsState{
				{
					From:    from,
					To:      to,
					Status:  "down",
					Message: "URL should be 'http://www.example12345.com/' but is 'http://www.example.com/'.",
				},
			},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/tms/check/%d/report/status", id), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
  "report": {
    "check_id": 123,
    "name": "My awesome check",
    "states": [
      {
        "status": "down",
        "from": "2020-07-10T10:51:55.000Z",
        "to": "2020-07-14T07:25:15.000Z",
        "error_in_step": 2,
        "message": "URL should be 'http://www.example12345.com/' but is 'http://www.example.com/'."
      }
    ]
  }
}`)
	})

	resp, err := client.TmsChecks.StatusReportById(id, request)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, *resp)
}
func TestTmsCheckPerformanceReportSingleTransactionCheck(t *testing.T) {
	id := 1234
	setup()
	defer teardown()

	request := TmsPerformanceReportRequest{}

	from, _ := time.Parse(time.RFC3339, "2020-07-10T10:51:55.000Z")
	//to, _ := time.Parse(time.RFC3339, "2020-07-14T07:25:15.000Z")
	expectedResponse := TmsPerformanceReportResponse{
		Report: TmsPerformanceReport{
			CheckId: 123,
			Name:    "My awesome check",
			Intervals: []TmsInterval{
				{
					AverageResponse: 123,
					Downtime:        10,
					From:            from,
					Steps: []TmsStepStatus{
						{
							Step: TmsStep{
								Function: "go_to",
								Args: map[string]string{
									"checkbox": "string",
									"element":  "string",
									"form":     "string",
									"input":    "string",
									"radio":    "string",
									"password": "string",
									"option":   "string",
									"seconds":  "string",
									"url":      "http://www.google.com",
									"username": "string",
									"value":    "string",
									"select":   "string",
								},
							},
							AverageResponse: 123,
						},
					},
					Unmonitored: 50,
					Uptime:      230,
				},
			},
			Resolution: DAY,
		},
	}

	mux.HandleFunc(fmt.Sprintf("/tms/check/%d/report/performance", id), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
  "report": {
    "check_id": 123,
    "name": "My awesome check",
    "resolution": "day",
    "intervals": [
      {
        "average_response": 123,
        "from": "2020-07-10T10:51:55.000Z",
        "downtime": 10,
        "uptime": 230,
        "unmonitored": 50,
        "steps": [
          {
            "step": {
              "args": {
                "checkbox": "string",
                "element": "string",
                "form": "string",
                "input": "string",
                "option": "string",
                "password": "string",
                "radio": "string",
                "seconds": "string",
                "select": "string",
                "url": "http://www.google.com",
                "username": "string",
                "value": "string"
              },
              "fn": "go_to"
            },
            "average_response": 123
          }
        ]
      }
    ]
  }
}`)
	})

	resp, err := client.TmsChecks.PerformanceReport(id, request)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, *resp)
}
func TestTmsCheckServiceCreateReal(t *testing.T) {
	// test client
	client, _ = NewClientWithConfig(ClientConfig{
		APIToken: "Ji7RnHUVIl2p7pPcd5KKafXLgDhc8NzIJOuqwXn7VSSbsp-iT7A296YQB6-iOA-CYi0-XJk",
	})

	defer teardown()

	want := &TmsCheckResponse{
		ID:   1003,
		Name: "Test redirect",
	}

	checks, err := client.TmsChecks.Create(TmsCheck{
		Name: "Test redirect",
		Steps: []TmsStep{
			{
				Function: "go_to",
				Args:     map[string]string{"url": "https://example.com/"},
			},
			{
				Function: "url",
				Args:     map[string]string{"url": "https://example.com/redirected"},
			},
		}})
	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}
