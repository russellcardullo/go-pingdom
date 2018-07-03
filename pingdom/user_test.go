package pingdom

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"fmt"
)

func TestUserServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
	    "users": [
        {
            "id": 12,
            "name": "John Doe",
            "paused": "NO",
            "access_level": "contact",
            "use_severity_levels": false,
            "sms": [
                {
                    "id": 352,
                    "severity": "HIGH",
                    "country_code": "1",
                    "number": "6095555555",
                    "provider": "nexmo"
                }
            ]
        },
        {
            "id": 234,
            "name": "Jane Doe",
            "paused": "NO",
            "access_level": "default",
            "use_severity_levels": false,
            "email": [
                {
                    "id": 314,
                    "severity": "HIGH",
                    "address": "test@billtrust.com"
                }
            ]
		}
		]
		}`)
	})
	want := []UsersResponse{
		{
			Id:       12,
			Paused:   "NO",
			Username: "John Doe",
			Sms:      []UserSmsResponse{
				{
					Id: 352,
					Severity: "HIGH",
					CountryCode: "1",
					Number: "6095555555",
					Provider: "nexmo",
				},
			},
			Email:    nil,
		},
		{
			Id:       234,
			Paused:   "NO",
			Username: "Jane Doe",
			Sms:      nil,
			Email:    []UserEmailResponse{
				{
					Id: 314,
					Severity: "HIGH",
					Address: "test@billtrust.com",
				},
			},
		},
	}

	users, err := client.Users.List()
	assert.NoError(t, err)
	assert.Equal(t, want, users, "Users.List() should return correct result")
}

func TestUserServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"user": {
				"id": 23439
			}
		}`)
	})

	want := &UsersResponse{
		Id: 23439,
	}

	u := User{
		Username : "testUser",
	}

	user, err := client.Users.Create(&u)
	assert.NoError(t, err)
	assert.Equal(t, want, user, "Users.Create() should return correct result")
}