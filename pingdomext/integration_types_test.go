package pingdomext

import (
	"reflect"
	"testing"
)

func TestWebHookIntegration_PostParams(t *testing.T) {

	tests := []struct {
		name        string
		integration WebHookIntegration
		wantParams  map[string]string
	}{
		{
			name: "parametrizes webhook integration",
			integration: WebHookIntegration{
				Active:     true,
				ProviderID: 2,
				UserData: &WebHookData{
					Name: "wlwu-test-12",
					URL:  "https://www.example.com",
				},
			},
			wantParams: map[string]string{
				"active":      "true",
				"provider_id": "2",
				"data_json":   `{"name":"wlwu-test-12","url":"https://www.example.com"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.integration.PostParams(); !reflect.DeepEqual(got, tt.wantParams) {
				t.Errorf("WebHookIntegration.PostParams() = %v, want %v", got, tt.wantParams)
			}
		})
	}
}

func TestWebHookIntegration_Valid(t *testing.T) {
	tests := []struct {
		name        string
		integration WebHookIntegration
		wantErr     bool
	}{
		{
			name: "parametrizes webhook integration",
			integration: WebHookIntegration{
				Active:     true,
				ProviderID: 2,
				UserData: &WebHookData{
					Name: "wlwu-test-12",
					URL:  "https://www.example.com",
				},
			},
			wantErr: false,
		},

		{
			name: "parametrizes webhook integration",
			integration: WebHookIntegration{
				Active:     false,
				ProviderID: 3,
				UserData: &WebHookData{
					Name: "wlwu-test-12",
					URL:  "https://www.example.com",
				},
			},
			wantErr: true,
		},

		{
			name: "parametrizes webhook integration",
			integration: WebHookIntegration{
				Active:     false,
				ProviderID: 2,
				UserData: &WebHookData{
					Name: "",
					URL:  "https://www.example.com",
				},
			},
			wantErr: true,
		},

		{
			name: "parametrizes webhook integration",
			integration: WebHookIntegration{
				Active:     false,
				ProviderID: 2,
				UserData: &WebHookData{
					Name: "11111",
					URL:  "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.integration.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("WebHookIntegration.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
