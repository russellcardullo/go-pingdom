package pingdom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTMSCheck_Valid(t *testing.T) {

	tests := []struct {
		name     string
		tmsCheck TMSCheck
		wantErr  error
	}{
		{
			name: "RequireParams",
			tmsCheck: TMSCheck{
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
			wantErr: nil,
		},
		{
			name: "NoError",
			tmsCheck: TMSCheck{
				Name: "NoError",
				Steps: []TMSCheckStep{
					{
						Args: map[string]string{
							"url": "www.google.com",
						},
						Fn: "go_to",
					},
				},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: nil,
		},
		{
			name: "EmptyName",
			tmsCheck: TMSCheck{
				Name: "",
				Steps: []TMSCheckStep{
					{
						Args: map[string]string{
							"url": "www.google.com",
						},
						Fn: "go_to",
					},
				},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `Name`. Must contain non-empty string."),
		},
		{
			name: "NilSteps",
			tmsCheck: TMSCheck{
				Name:                     "NilSteps",
				Steps:                    nil,
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `Steps`. Must contain non-empty value."),
		},
		{
			name: "EmptySteps",
			tmsCheck: TMSCheck{
				Name:                     "EmptySteps",
				Steps:                    []TMSCheckStep{},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `Steps`. Must contain non-empty value."),
		},
		{
			name: "InvalidInterval",
			tmsCheck: TMSCheck{
				Name: "InvalidInterval",
				Steps: []TMSCheckStep{
					{
						Args: map[string]string{
							"url": "www.google.com",
						},
						Fn: "go_to",
					},
				},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 13,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `Interval`. Please provide one of the following valid values instead: [5 10 20 60 720 1440]."),
		},
		{
			name: "InvalidSeverityLevel",
			tmsCheck: TMSCheck{
				Name: "InvalidSeverityLevel",
				Steps: []TMSCheckStep{
					{
						Args: map[string]string{
							"url": "www.google.com",
						},
						Fn: "go_to",
					},
				},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high1",
				Tags:                     []string{"aaa", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `SeverityLevel`. Please provide one of the following valid values instead: [high,low]."),
		},
		{
			name: "InvalidTags",
			tmsCheck: TMSCheck{
				Name: "InvalidTags",
				Steps: []TMSCheckStep{
					{
						Args: map[string]string{
							"url": "www.google.com",
						},
						Fn: "go_to",
					},
				},
				Active:                   true,
				ContactIDs:               []int{12345},
				CustomMessage:            "custome_msg",
				IntegrationIDs:           []int{12345},
				Interval:                 10,
				Metadata:                 &TMSCheckMetaData{},
				Region:                   "us-east",
				SendNotificationWhenDown: 1,
				SeverityLevel:            "high",
				Tags:                     []string{"aaa$", "bbb"},
				TeamIDs:                  []int{12345},
			},
			wantErr: fmt.Errorf("Invalid value for `Tags`. The tag name may contain the characters 'A-Z', 'a-z', '0-9', '_' and '-'."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tmsCheck.Valid()
			if err == nil {
				assert.NoError(t, tt.wantErr)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}

		})
	}
}

func TestTMSCheck_RenderForJSONAPI(t *testing.T) {
	tests := []struct {
		name     string
		tmsCheck TMSCheck
		wantJson string
	}{
		{
			name: "RequireParams",
			tmsCheck: TMSCheck{
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
			wantJson: `{"name":"RequireParams","steps":[{"args":{"url":"www.google.com"},"fn":"go_to"}],"active":false}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.tmsCheck.RenderForJSONAPI(); got != tt.wantJson {
				t.Errorf("TMSCheck.RenderForJSONAPI() = %v, want %v", got, tt.wantJson)
			}
		})
	}
}
