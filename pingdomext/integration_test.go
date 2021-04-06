package pingdomext

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationService_Create(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/data/v3/integration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		assert.Equal(t, "false", r.URL.Query().Get("active"))
		assert.Equal(t, "2", r.URL.Query().Get("provider_id"))
		assert.Equal(t, `{"name":"wlwu-test-5","url":"http://www.example.org"}`, r.URL.Query().Get("data_json"))
		fmt.Fprint(w, `{
			"integration": {
				"id":112107,
				"status":true
			}
		}`)
	})

	want := IntegrationStatus{
		ID:     112107,
		Status: true,
	}

	tests := []struct {
		name        string
		client      *Client
		integration Integration
		want        *IntegrationStatus
		wantErr     bool
	}{

		{
			name:   "valid",
			client: client,
			integration: &WebHookIntegration{
				Active:     false,
				ProviderID: 2,
				UserData: &WebHookData{
					Name: "wlwu-test-5",
					URL:  "http://www.example.org",
				},
			},
			want:    &want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.Create(tt.integration)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegrationService_List(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/data/v3/integration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"integration": [
				{
					"provider_id": 2,
					"name": "webhook",
					"user_data": {
						"name": "Slack #mct-internal",
						"url": "https://hooks.slack.com/services"
					},
					"created_at": 1615819798,
					"id": 112107,
					"description": "Webhook",
					"provider_data": [
						{
							"required": true,
							"description": "URL",
							"name": "url",
							"validation_options": {},
							"validation": "url"
						},
						{
							"required": true,
							"description": "Name",
							"name": "name",
							"validation_options": {},
							"validation": "string"
						},
						{
							"required": false,
							"description": "Use legacy BeepManager format",
							"name": "legacy",
							"validation_options": {},
							"validation": "bool"
						}
					],
					"activated_at": 1615819798,
					"number_of_connected_checks": 2
				},
				{
					"provider_id": 2,
					"name": "webhook",
					"user_data": {
						"name": "wlwu-test-5",
						"url": "http://www.example.org"
					},
					"created_at": 1615969145,
					"id": 112165,
					"description": "Webhook",
					"provider_data": [
						{
							"required": true,
							"description": "URL",
							"name": "url",
							"validation_options": {},
							"validation": "url"
						},
						{
							"required": true,
							"description": "Name",
							"name": "name",
							"validation_options": {},
							"validation": "string"
						},
						{
							"required": false,
							"description": "Use legacy BeepManager format",
							"name": "legacy",
							"validation_options": {},
							"validation": "bool"
						}
					],
					"activated_at": null,
					"number_of_connected_checks": 0
				}
			]
		}`)
	})

	want := []IntegrationGetResponse{
		{
			NumberOfConnectedChecks: 2,
			ID:                      112107,
			Name:                    "webhook",
			Description:             "Webhook",
			ProviderID:              2,
			ActivatedAt:             1615819798,
			CreatedAt:               1615819798,
			UserData: map[string]string{
				"name": "Slack #mct-internal",
				"url":  "https://hooks.slack.com/services",
			},
		},
		{
			NumberOfConnectedChecks: 0,
			ID:                      112165,
			Name:                    "webhook",
			Description:             "Webhook",
			ProviderID:              2,
			ActivatedAt:             0,
			CreatedAt:               1615969145,
			UserData: map[string]string{
				"name": "wlwu-test-5",
				"url":  "http://www.example.org",
			},
		},
	}

	tests := []struct {
		name    string
		client  *Client
		want    []IntegrationGetResponse
		wantErr bool
	}{
		{
			name:    "valid",
			client:  client,
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegrationService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data/v3/integration/112107", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"integration": 
				{
					"provider_id": 2,
					"name": "webhook",
					"user_data": {
						"name": "Slack #mct-internal",
						"url": "https://hooks.slack.com/services"
					},
					"created_at": 1615819798,
					"id": 112107,
					"description": "Webhook",
					"provider_data": [
						{
							"required": true,
							"description": "URL",
							"name": "url",
							"validation_options": {},
							"validation": "url"
						},
						{
							"required": true,
							"description": "Name",
							"name": "name",
							"validation_options": {},
							"validation": "string"
						},
						{
							"required": false,
							"description": "Use legacy BeepManager format",
							"name": "legacy",
							"validation_options": {},
							"validation": "bool"
						}
					],
					"activated_at": 1615819798,
					"number_of_connected_checks": 2
				}
			
		}`)
	})

	mux.HandleFunc("/data/v3/integration/112108", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"integration": null}`)
	})

	want := IntegrationGetResponse{

		NumberOfConnectedChecks: 2,
		ID:                      112107,
		Name:                    "webhook",
		Description:             "Webhook",
		ProviderID:              2,
		ActivatedAt:             1615819798,
		CreatedAt:               1615819798,
		UserData: map[string]string{
			"name": "Slack #mct-internal",
			"url":  "https://hooks.slack.com/services",
		},
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *IntegrationGetResponse
		wantErr bool
	}{
		{
			name:    "valid",
			client:  client,
			args:    args{id: 112107},
			want:    &want,
			wantErr: false,
		},
		{
			name:    "null",
			client:  client,
			args:    args{id: 112108},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid",
			client:  client,
			args:    args{id: 112109},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.Read(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegrationService_Update(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/data/v3/integration/112107", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		assert.Equal(t, "false", r.URL.Query().Get("active"))
		assert.Equal(t, "2", r.URL.Query().Get("provider_id"))
		assert.Equal(t, `{"name":"wlwu-test-5","url":"http://www.example.org"}`, r.URL.Query().Get("data_json"))
		fmt.Fprint(w, `{
			"integration": {
				"status":true
			}
		}`)
	})

	want := IntegrationStatus{
		ID:     0,
		Status: true,
	}

	type args struct {
		id          int
		integration Integration
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *IntegrationStatus
		wantErr bool
	}{

		{
			name:   "valid",
			client: client,
			args: args{
				id: 112107,
				integration: &WebHookIntegration{
					Active:     false,
					ProviderID: 2,
					UserData: &WebHookData{
						Name: "wlwu-test-5",
						URL:  "http://www.example.org",
					},
				},
			},
			want:    &want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.Update(tt.args.id, tt.args.integration)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegrationService_Delete(t *testing.T) {

	setup()
	defer teardown()

	mux.HandleFunc("/data/v3/integration/112169", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"integration": {
				"status":true
			}
		}`)
	})

	mux.HandleFunc("/data/v3/integration/112166", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, "Somthing went wrong!")
	})

	want := IntegrationStatus{
		ID:     0,
		Status: true,
	}

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *IntegrationStatus
		wantErr bool
	}{
		{
			name:   "valid",
			client: client,
			args: args{
				id: 112169,
			},
			want:    &want,
			wantErr: false,
		},
		{
			name:   "valid",
			client: client,
			args: args{
				id: 112166,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.Delete(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func TestIntegrationService_ListProviders(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integrations/provider", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"data": [
					{
						"required": true,
						"description": "E-Mail",
						"name": "email",
						"validation_options": {},
						"validation": "email"
					},
					{
						"required": true,
						"description": "API Token",
						"name": "api_key",
						"validation_options": {
							"min_length": 64,
							"max_length": 64
						},
						"validation": "string"
					},
					{
						"required": true,
						"description": "Name",
						"name": "name",
						"validation_options": {},
						"validation": "string"
					}
				],
				"id": 1,
				"description": "Librato",
				"name": "librato"
			},
			{
				"data": [
					{
						"required": true,
						"description": "URL",
						"name": "url",
						"validation_options": {},
						"validation": "url"
					},
					{
						"required": true,
						"description": "Name",
						"name": "name",
						"validation_options": {},
						"validation": "string"
					},
					{
						"required": false,
						"description": "Use legacy BeepManager format",
						"name": "legacy",
						"validation_options": {},
						"validation": "bool"
					}
				],
				"id": 2,
				"description": "Webhook",
				"name": "webhook"
			}
		]`)
	})

	want := []IntegrationProvider{
		{
			ID:          1,
			Name:        "librato",
			Description: "Librato",
		},
		{

			ID:          2,
			Name:        "webhook",
			Description: "Webhook",
		},
	}

	tests := []struct {
		name    string
		client  *Client
		want    []IntegrationProvider
		wantErr bool
	}{
		{
			name:    "valid",
			client:  client,
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &IntegrationService{
				client: tt.client,
			}
			got, err := cs.ListProviders()
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrationService.ListProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntegrationService.ListProviders() = %v, want %v", got, tt.want)
			}
		})
	}
}
