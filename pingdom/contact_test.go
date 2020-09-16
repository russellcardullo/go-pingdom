package pingdom

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
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
		}`)
	})
	want := []ContactResponse{
		{
			ID:     1,
			Paused: false,
			Name:   "John Doe",
			Owner:  true,
			Teams: []ContactTeamResponse{
				{
					ID:   123456,
					Name: "The Dream Team",
				},
			},
			Type: "user",
			NotificationTargets: NotificationTargetsResponse{
				SMS: []SMSNotificationResponse{
					{
						Severity:    "HIGH",
						CountryCode: "00",
						Number:      "111111111",
						Provider:    "provider's name",
					},
				},
				Email: []EmailNotificationResponse{
					{
						Severity: "HIGH",
						Address:  "johndoe@teamrocket.com",
					},
				},
			},
		},
		{
			ID:     2,
			Paused: true,
			Name:   "John \"Hannibal\" Smith",
			Type:   "user",
			Teams:  []ContactTeamResponse{},
			NotificationTargets: NotificationTargetsResponse{
				SMS: []SMSNotificationResponse{
					{
						Severity:    "HIGH",
						CountryCode: "00",
						Number:      "222222222",
						Provider:    "provider's name",
					},
				},
				Email: []EmailNotificationResponse{
					{
						Severity: "HIGH",
						Address:  "hannibal@ateam.org",
					},
				},
			},
		},
	}

	contacts, err := client.Contacts.List()
	assert.NoError(t, err)
	assert.Equal(t, want, contacts, "Contacts.List() should return correct result")
}

func TestContactService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/contacts/123456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"contact": {
				"id": 123456,
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
			}
		  }`)
	})
	want := ContactResponse{
		ID:     123456,
		Paused: false,
		Name:   "John Doe",
		Owner:  true,
		Type:   "user",
		NotificationTargets: NotificationTargetsResponse{
			SMS: []SMSNotificationResponse{
				SMSNotificationResponse{
					Severity:    "HIGH",
					CountryCode: "00",
					Number:      "111111111",
					Provider:    "provider's name",
				},
			},
			Email: []EmailNotificationResponse{
				{
					Address:  "johndoe@teamrocket.com",
					Severity: "HIGH",
				},
			},
		},
		Teams: []ContactTeamResponse{
			{
				ID:   123456,
				Name: "The Dream Team",
			},
		},
	}

	contacts, err := client.Contacts.Read(123456)
	assert.NoError(t, err)
	assert.Equal(t, &want, contacts, "Contacts.Read(123456) should return a contact")
}

func TestContactService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"contact": {
				"id": 23439
			}
		}`)
	})

	want := &ContactResponse{
		ID: 23439,
	}

	u := Contact{
		Name: "testContact",
	}

	contact, err := client.Contacts.Create(&u)
	assert.NoError(t, err)
	assert.Equal(t, want, contact, "Contacts.Create() should return correct result")
}

func TestContactService_Delete(t *testing.T) {
	setup()
	defer teardown()

	contactID := 12941

	mux.HandleFunc("/alerting/contacts/"+strconv.Itoa(contactID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message":"Deletion of contact was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Deletion of contact was successful!",
	}

	response, err := client.Contacts.Delete(contactID)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Contacts.Delete() should return PingdomResponse with message")

}

func TestContactService_Update(t *testing.T) {
	setup()
	defer teardown()

	contactID := 12941
	contact := Contact{
		Name: "updatedName",
	}

	mux.HandleFunc("/alerting/contacts/"+strconv.Itoa(contactID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"message":"Modification of contact was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Modification of contact was successful!",
	}

	response, err := client.Contacts.Update(contactID, &contact)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Contacts.Update() should return PingdomResponse with message")

}
