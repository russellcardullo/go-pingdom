package pingdom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var detailedCheckJSON = `
{
	"id" : 85975,
	"name" : "My check 7",
	"resolution" : 1,
	"sendnotificationwhendown" : 0,
	"notifyagainevery" : 0,
	"notifywhenbackup" : false,
	"created" : 1240394682,
	"type" : {
		"http" : {
			"url" : "/",
			"port" : 80,
			"requestheaders" : {
				"User-Agent" : "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				"Prama" : "no-cache"
			}
		}
	},
	"hostname" : "s7.mydomain.com",
	"status" : "up",
	"severity_level": "HIGH",
	"lasterrortime" : 1293143467,
	"lasttesttime" : 1294064823,
	"tags": [],
	"responsetime_threshold": 2300
}
`

func TestPingdomError(t *testing.T) {
	pe := PingdomError{StatusCode: 400, StatusDesc: "Bad Request", Message: "Missing param foo"}
	want := "400 Bad Request: Missing param foo"
	assert.Equal(t, want, pe.Error())
}

func TestCheckResponseUnmarshal(t *testing.T) {
	var ck CheckResponse
	err := json.Unmarshal([]byte(detailedCheckJSON), &ck)
	assert.NoError(t, err)
	assert.Equal(t, "http", ck.Type.Name)
	assert.NotNil(t, ck.Type.HTTP)
	assert.Equal(t, 2, len(ck.Type.HTTP.RequestHeaders))
	assert.Equal(t, "HIGH", ck.SeverityLevel)
}

var detailedContactJSON = `
{
	"contacts": [
		{
			"id": 1,
			"name": "John Doe",
			"paused": false,
			"type": "user",
			"owner": true,
			"notification_targets": {
				"email": [
					{
					"severity": "HIGH",
					"address": "johndoe@teamrocket.com"
					}
				],
				"sms": [
					{
					"severity": "HIGH",
					"country_code": "00",
					"number": "111111111",
					"provider": "provider's name"
					}
				]
			},
			"teams": [
				{
					"id": 123456,
					"name": "The Dream Team"
				}
			]
		},
		{
			"id": 2,
			"name": "John \"Hannibal\" Smith",
			"paused": true,
			"type": "user",
			"notification_targets": {
			"email": [
				{
				"severity": "HIGH",
				"address": "hannibal@ateam.org"
				}
			],
			"sms": [
				{
				"severity": "HIGH",
				"country_code": "00",
				"number": "222222222",
				"provider": "provider's name"
				}
			]
			},
			"teams": []
		}
	]
}
`

func TestCheckContactResponseUnmarshal(t *testing.T) {
	var contacts listContactsJSONResponse
	err := json.Unmarshal([]byte(detailedContactJSON), &contacts)
	contact := contacts.Contacts[0]

	expectedNotificationTargets := NotificationTargetsResponse{
		SMS: []SMSNotificationResponse{
			SMSNotificationResponse{
				Severity:    "HIGH",
				CountryCode: "00",
				Number:      "111111111",
				Provider:    "provider's name",
			},
		},
		Email: []EmailNotificationResponse{
			EmailNotificationResponse{
				Severity: "HIGH",
				Address:  "johndoe@teamrocket.com",
			},
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", contact.Name)
	assert.NotNil(t, contact.ID)
	assert.Equal(t, expectedNotificationTargets, contact.NotificationTargets)
}
