package pingdom

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"fmt"
	"strconv"
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

func TestUserServiceContactCreate(t *testing.T) {
	setup()
	defer teardown()

	userId := 12941

	mux.HandleFunc("/users/" + strconv.Itoa(userId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"contact_target": {
				"id": 23439
			}
		}`)
	})

	want := &CreateUserContactResponse{
		Id: 23439,
	}

	c := Contact{
		Email: "test@example.com",
		CountryCode: "1",
		Number: "5559995555",
	}

	contact, err := client.Users.CreateContact(userId, c)
	assert.NoError(t, err)
	assert.Equal(t, want, contact, "Users.CreateContact() should return contact_target.id")
}

func TestUserServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	userId := 12941

	mux.HandleFunc("/users/" + strconv.Itoa(userId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message":"Deletion of user was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Deletion of user was successful!",
	}

	response, err := client.Users.Delete(userId)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.Delete() should return PingdomResponse with message")

}

func TestUserService_DeleteContact(t *testing.T) {
	setup()
	defer teardown()

	userId := 12941
	contactId := 87655

	mux.HandleFunc("/users/" + strconv.Itoa(userId)+"/" + strconv.Itoa(contactId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message":"Deletion of contact target successful"
		}`)
	})

	want := &PingdomResponse{
		Message: "Deletion of contact target successful",
	}

	response, err := client.Users.DeleteContact(userId,contactId)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.DeleteContact() should return PingdomResponse with message")

}