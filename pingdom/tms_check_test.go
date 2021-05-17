package pingdom

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTMSCheckService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"checks": [
				{
					"type": "script",
					"id": 104757,
					"name": "test1",
					"active": true,
					"status": "failing",
					"created_at": 1615778672,
					"interval": 10,
					"region": "us-east",
					"modified_at": 1618302994,
					"last_downtime_start": 1618061106,
					"last_downtime_end": 1619580306,
					"tags": []
				},
				{
					"type": "script",
					"id": 106136,
					"name": "test2",
					"active": true,
					"status": "failing",
					"created_at": 1619511011,
					"interval": 10,
					"region": "us-east",
					"modified_at": 1619574751,
					"last_downtime_start": 1619511179,
					"last_downtime_end": 1619578990,
					"tags": []
				},
				{
					"type": "script",
					"id": 106164,
					"name": "test3",
					"active": true,
					"status": "successful",
					"created_at": 1619574885,
					"interval": 10,
					"region": "us-east",
					"modified_at": 1619574885,
					"tags": []
				}
			],
			"limit": 1000,
			"offset": 0
		}`)
	})

	want := []TMSCheckResponse{
		{
			ID:                104757,
			Name:              "test1",
			Type:              "script",
			Active:            true,
			Status:            "failing",
			Interval:          10,
			Region:            "us-east",
			Tags:              []string{},
			LastDowntimeStart: 1618061106,
			LastDowntimeEnd:   1619580306,
			CreatedAt:         1615778672,
			ModifiedAt:        1618302994,
		},
		{
			ID:                106136,
			Name:              "test2",
			Type:              "script",
			Active:            true,
			Status:            "failing",
			Interval:          10,
			Region:            "us-east",
			Tags:              []string{},
			LastDowntimeStart: 1619511179,
			LastDowntimeEnd:   1619578990,
			CreatedAt:         1619511011,
			ModifiedAt:        1619574751,
		},
		{
			ID:                106164,
			Name:              "test3",
			Type:              "script",
			Active:            true,
			Status:            "successful",
			Interval:          10,
			Region:            "us-east",
			Tags:              []string{},
			LastDowntimeStart: 0,
			LastDowntimeEnd:   0,
			CreatedAt:         1619574885,
			ModifiedAt:        1619574885,
		},
	}

	type args struct {
		params []map[string]string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    []TMSCheckResponse
		wantErr bool
	}{
		{
			name:    "Valied",
			client:  client,
			args:    args{},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.List(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/104757", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"check": {
				"id": 104757,
				"type": "script",
				"name": "NoError",
				"steps": [
					{
						"fn": "go_to",
						"args": {
							"url": "www.google.com"
						}
					},
					{
						"fn": "click",
						"args": {
							"element": "Test"
						}
					}
				],
				"contact_ids": [
					12345
				],
				"team_ids": [],
				"integration_ids": [],
				"send_notification_when_down": 1,
				"severity_level": "high",
				"active": true,
				"status": "failing",
				"created_at": 1615778672,
				"interval": 10,
				"region": "us-east",
				"modified_at": 1618302994,
				"last_downtime_start": 1618061106,
				"last_downtime_end": 1619581506,
				"tags": []
			}
		}`)
	})

	want := &TMSCheckDetailResponse{
		TMSCheck: TMSCheck{
			Name: "NoError",
			Steps: []TMSCheckStep{
				{
					Args: map[string]string{
						"url": "www.google.com",
					},
					Fn: "go_to",
				},
				{
					Args: map[string]string{
						"element": "Test",
					},
					Fn: "click",
				},
			},
			Active:                   true,
			ContactIDs:               []int{12345},
			CustomMessage:            "",
			IntegrationIDs:           []int{},
			Interval:                 10,
			Metadata:                 nil,
			Region:                   "us-east",
			SendNotificationWhenDown: 1,
			SeverityLevel:            "high",
			Tags:                     []string{},
			TeamIDs:                  []int{},
		},
		ID:                104757,
		Type:              "script",
		Status:            "failing",
		LastDowntimeStart: 1618061106,
		LastDowntimeEnd:   1619581506,
		CreatedAt:         1615778672,
		ModifiedAt:        1618302994,
	}

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *TMSCheckDetailResponse
		wantErr bool
	}{
		{
			name:   "Valied",
			client: client,
			args: args{
				id: 104757,
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.Read(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_Create(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"check": {
				"id": 106170,
				"name": "wlwu-test-3"
			}
		}`)
	})

	type args struct {
		tmsCheck *TMSCheck
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *TMSCheckDetailResponse
		wantErr bool
	}{
		{
			name:   "Valied",
			client: client,
			args: args{
				tmsCheck: &TMSCheck{
					Name: "RequireParams",
					Steps: []TMSCheckStep{
						{
							Args: map[string]string{
								"url": "www.google.com",
							},
							Fn: "go_to",
						},
					},
				},
			},
			want: &TMSCheckDetailResponse{
				TMSCheck: TMSCheck{
					Name: "wlwu-test-3",
				},
				ID: 106170,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.Create(tt.args.tmsCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_Update(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/104757", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"check": {
				"id": 104757,
				"type": "script",
				"name": "NoError",
				"steps": [
					{
						"fn": "go_to",
						"args": {
							"url": "www.google234.com"
						}
					},
					{
						"fn": "click",
						"args": {
							"element": "Test234"
						}
					}
				],
				"contact_ids": [
					123456
				],
				"team_ids": [123456],
				"integration_ids": [123456],
				"send_notification_when_down": 1,
				"severity_level": "high",
				"active": true,
				"status": "failing",
				"created_at": 1615778672,
				"interval": 10,
				"region": "us-east",
				"modified_at": 1618302994,
				"last_downtime_start": 1618061106,
				"last_downtime_end": 1619581506,
				"tags": ["aaa","bbb"]
			}
		}`)
	})

	changedTMSCheck := TMSCheck{
		Name: "NoError",
		Steps: []TMSCheckStep{
			{
				Args: map[string]string{
					"url": "www.google234.com",
				},
				Fn: "go_to",
			},
			{
				Args: map[string]string{
					"element": "Test234",
				},
				Fn: "click",
			},
		},
		Active:                   true,
		ContactIDs:               []int{123456},
		CustomMessage:            "",
		IntegrationIDs:           []int{123456},
		Interval:                 10,
		Metadata:                 nil,
		Region:                   "us-east",
		SendNotificationWhenDown: 1,
		SeverityLevel:            "high",
		Tags:                     []string{"aaa", "bbb"},
		TeamIDs:                  []int{123456},
	}
	want := &TMSCheckDetailResponse{
		TMSCheck:          changedTMSCheck,
		ID:                104757,
		Type:              "script",
		Status:            "failing",
		LastDowntimeStart: 1618061106,
		LastDowntimeEnd:   1619581506,
		CreatedAt:         1615778672,
		ModifiedAt:        1618302994,
	}

	type args struct {
		id       int
		tmsCheck *TMSCheck
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *TMSCheckDetailResponse
		wantErr bool
	}{
		{
			name:   "Valied",
			client: client,
			args: args{
				id:       104757,
				tmsCheck: &changedTMSCheck,
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.Update(tt.args.id, tt.args.tmsCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_Delete(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/104757", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message": "Deletion of check 104757 was successful."
		}`)
	})

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *PingdomResponse
		wantErr bool
	}{

		{
			name:   "NoError",
			client: client,
			args: args{
				id: 104757,
			},
			want: &PingdomResponse{
				Message: "Deletion of check 104757 was successful.",
			},
			wantErr: false,
		},

		{
			name:   "Error",
			client: client,
			args: args{
				id: 104758,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.Delete(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_GetStatusReport(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/104757/report/status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"report": {
				"check_id": 104757,
				"name": "test",
				"states": [
					{
						"status": "down",
						"from": "2021-04-22T06:34:41Z",
						"to": "2021-04-28T12:45:14Z",
						"message": "Element 'View Create a cluster' does not exist.",
						"error_in_step": 4
					},
					{
						"status": "down",
						"from": "2021-04-28T12:45:14Z",
						"to": "2021-04-28T12:55:06Z",
						"message": "Timed out (>60s)",
						"error_in_step": 2
					},
					{
						"status": "down",
						"from": "2021-04-28T12:55:06Z",
						"to": "2021-04-28T15:05:06Z",
						"message": "Element 'View Create a cluster' does not exist.",
						"error_in_step": 4
					},
					{
						"status": "down",
						"from": "2021-04-28T15:05:06Z",
						"to": "2021-04-28T15:15:06Z",
						"message": "Timed out (>60s)",
						"error_in_step": 1
					},
					{
						"status": "down",
						"from": "2021-04-28T15:15:06Z",
						"to": "2021-04-29T06:25:06Z",
						"message": "Element 'View Create a cluster' does not exist.",
						"error_in_step": 4
					}
				]
			}
		}`)
	})

	want := &TMSCheckStatusReportResponse{
		CheckID: 104757,
		Name:    "test",
		States: []TMSCheckStatus{
			{
				ErrorInStep: 4,
				From:        "2021-04-22T06:34:41Z",
				To:          "2021-04-28T12:45:14Z",
				Message:     "Element 'View Create a cluster' does not exist.",
				Status:      "down",
			},
			{
				ErrorInStep: 2,
				From:        "2021-04-28T12:45:14Z",
				To:          "2021-04-28T12:55:06Z",
				Message:     "Timed out (>60s)",
				Status:      "down",
			},
			{
				ErrorInStep: 4,
				From:        "2021-04-28T12:55:06Z",
				To:          "2021-04-28T15:05:06Z",
				Message:     "Element 'View Create a cluster' does not exist.",
				Status:      "down",
			},
			{
				ErrorInStep: 1,
				From:        "2021-04-28T15:05:06Z",
				To:          "2021-04-28T15:15:06Z",
				Message:     "Timed out (>60s)",
				Status:      "down",
			},
			{
				ErrorInStep: 4,
				From:        "2021-04-28T15:15:06Z",
				To:          "2021-04-29T06:25:06Z",
				Message:     "Element 'View Create a cluster' does not exist.",
				Status:      "down",
			},
		},
	}

	type args struct {
		id     int
		params map[string]string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *TMSCheckStatusReportResponse
		wantErr bool
	}{
		{
			name:   "NoError",
			client: client,
			args: args{
				id: 104757,
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.GetStatusReport(tt.args.id, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.getStatusReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.getStatusReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_ListStatusReports(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/report/status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"report": [
				{
					"check_id": 104757,
					"name": "test",
					"states": [
						{
							"status": "down",
							"from": "2021-04-22T06:34:41Z",
							"to": "2021-04-28T12:45:14Z",
							"message": "Element 'View Create a cluster' does not exist.",
							"error_in_step": 4
						},
						{
							"status": "down",
							"from": "2021-04-28T12:45:14Z",
							"to": "2021-04-28T12:55:06Z",
							"message": "Timed out (>60s)",
							"error_in_step": 2
						},
						{
							"status": "down",
							"from": "2021-04-28T12:55:06Z",
							"to": "2021-04-28T15:05:06Z",
							"message": "Element 'View Create a cluster' does not exist.",
							"error_in_step": 4
						},
						{
							"status": "down",
							"from": "2021-04-28T15:05:06Z",
							"to": "2021-04-28T15:15:06Z",
							"message": "Timed out (>60s)",
							"error_in_step": 1
						},
						{
							"status": "down",
							"from": "2021-04-28T15:15:06Z",
							"to": "2021-04-29T06:25:06Z",
							"message": "Element 'View Create a cluster' does not exist.",
							"error_in_step": 4
						}
					]
				},
				{
					"check_id": 106136,
					"name": "test2",
					"states": [
						{
							"status": "unknown",
							"from": "2021-04-22T06:58:57Z",
							"to": "2021-04-27T08:12:59Z"
						},
						{
							"status": "down",
							"from": "2021-04-27T08:12:59Z",
							"to": "2021-04-29T06:52:59Z",
							"message": "Element 'View Create a cluster' does not exist.",
							"error_in_step": 4
						}
					]
				}
			],
			"limit": 100,
			"offset": 0,
			"omit_empty": false
		}`)
	})

	want := []TMSCheckStatusReportResponse{
		{
			CheckID: 104757,
			Name:    "test",
			States: []TMSCheckStatus{
				{
					ErrorInStep: 4,
					From:        "2021-04-22T06:34:41Z",
					To:          "2021-04-28T12:45:14Z",
					Message:     "Element 'View Create a cluster' does not exist.",
					Status:      "down",
				},
				{
					ErrorInStep: 2,
					From:        "2021-04-28T12:45:14Z",
					To:          "2021-04-28T12:55:06Z",
					Message:     "Timed out (>60s)",
					Status:      "down",
				},
				{
					ErrorInStep: 4,
					From:        "2021-04-28T12:55:06Z",
					To:          "2021-04-28T15:05:06Z",
					Message:     "Element 'View Create a cluster' does not exist.",
					Status:      "down",
				},
				{
					ErrorInStep: 1,
					From:        "2021-04-28T15:05:06Z",
					To:          "2021-04-28T15:15:06Z",
					Message:     "Timed out (>60s)",
					Status:      "down",
				},
				{
					ErrorInStep: 4,
					From:        "2021-04-28T15:15:06Z",
					To:          "2021-04-29T06:25:06Z",
					Message:     "Element 'View Create a cluster' does not exist.",
					Status:      "down",
				},
			},
		},
		{
			CheckID: 106136,
			Name:    "test2",
			States: []TMSCheckStatus{
				{
					From:   "2021-04-22T06:58:57Z",
					To:     "2021-04-27T08:12:59Z",
					Status: "unknown",
				},
				{
					ErrorInStep: 4,
					From:        "2021-04-27T08:12:59Z",
					To:          "2021-04-29T06:52:59Z",
					Message:     "Element 'View Create a cluster' does not exist.",
					Status:      "down",
				},
			},
		},
	}

	type args struct {
		params map[string]string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    []TMSCheckStatusReportResponse
		wantErr bool
	}{
		{
			name:    "NoError",
			client:  client,
			args:    args{},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.ListStatusReports(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.getStatusReports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.getStatusReports() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTMSCheckService_GetPerfomanceReport(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/tms/check/104757/report/performance", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"report": {
				"check_id": 104757,
				"name": "test",
				"resolution": "hour",
				"intervals": [
					{
						"steps": [
							{
								"average_response": 2507177,
								"step": {
									"fn": "go_to",
									"args": {
										"url": "www.google.com"
									}
								}
							},
							{
								"average_response": 2183649,
								"step": {
									"fn": "click",
									"args": {
										"element": "kubernetes"
									}
								}
							}
						],
						"average_response": 16304508,
						"from": "2021-04-28T23:00:00Z"
					},
					{
						"steps": [
							{
								"average_response": 2608565,
								"step": {
									"fn": "go_to",
									"args": {
										"url": "www.google.com"
									}
								}
							},
							{
								"average_response": 4191987,
								"step": {
									"fn": "click",
									"args": {
										"element": "kubernetes"
									}
								}
							}
						],
						"average_response": 19222481,
						"from": "2021-04-29T00:00:00Z"
					}
				]
			}
		}`)
	})

	want := &TMSCheckPerformanceReportResponse{
		CheckID:    104757,
		Name:       "test",
		Resolution: "hour",
		Intervals: []TMSCheckInterval{
			{
				AverageResponse: 16304508,
				From:            "2021-04-28T23:00:00Z",
				Steps: []TMSCheckStepReport{
					{
						AverageResponse: 2507177,
						Step: TMSCheckStep{
							Fn:   "go_to",
							Args: map[string]string{"url": "www.google.com"},
						},
					},
					{
						AverageResponse: 2183649,
						Step: TMSCheckStep{
							Fn:   "click",
							Args: map[string]string{"element": "kubernetes"},
						},
					},
				},
			},
			{
				AverageResponse: 19222481,
				From:            "2021-04-29T00:00:00Z",
				Steps: []TMSCheckStepReport{
					{
						AverageResponse: 2608565,
						Step: TMSCheckStep{
							Fn:   "go_to",
							Args: map[string]string{"url": "www.google.com"},
						},
					},
					{
						AverageResponse: 4191987,
						Step: TMSCheckStep{
							Fn:   "click",
							Args: map[string]string{"element": "kubernetes"},
						},
					},
				},
			},
		},
	}

	type args struct {
		id     int
		params map[string]string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *TMSCheckPerformanceReportResponse
		wantErr bool
	}{
		{
			name:   "NoError",
			client: client,
			args: args{
				id: 104757,
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &TMSCheckService{
				client: tt.client,
			}
			got, err := cs.GetPerfomanceReport(tt.args.id, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TMSCheckService.getPerfomanceReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMSCheckService.getPerfomanceReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
